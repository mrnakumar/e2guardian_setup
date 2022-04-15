package pkg

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type UploadOptions struct {
	UserName         string
	Password         string
	Url              string
	RecipientKeyPath string
	Interval         uint16
	BaseFolder       string
	FileSuffix       []string
	SizeLimit        int64
	filePath         string
}

type batch struct {
	files []string
	size  int64
	limit int64
}

type Uploader struct {
	options     UploadOptions
	wg          *sync.WaitGroup
	client      *http.Client
	contentType string
	encryptor   Encryptor
}

func CreateUploader(options UploadOptions, wg *sync.WaitGroup) (Uploader, error) {
	encryptor, err := CreateEncryptor(options.RecipientKeyPath)
	if err != nil {
		return Uploader{}, err
	}
	return Uploader{
		options:     options,
		wg:          wg,
		client:      &http.Client{},
		contentType: "multipart/form-data",
		encryptor:   encryptor,
	}, nil
}

func (b *batch) Add(file fileInfo) bool {
	if b.size+file.size < b.limit {
		b.files = append(b.files, file.path)
		b.size += file.size
		return true
	}
	return false
}

func (b *batch) Reset() {
	for _, file := range b.files {
		err := os.Remove(file)
		if err != nil {
			log.Printf("failed to delete file from batch. Caused by '%v'", err)
		}
	}
	b.files = make([]string, 0)
	b.size = 0
}

func (u Uploader) Worker() {
	defer u.wg.Done()
	for {
		files, err := ListFiles(u.options.FileSuffix, u.options.BaseFolder)
		if err != nil {
			log.Printf("failed to list files. Caused by '%s'", err)
		} else {
			batch := batch{
				files: make([]string, 0),
				size:  0,
				limit: u.options.SizeLimit,
			}
			for _, file := range files {
				if !batch.Add(file) {
					u.upload(batch.files)
					batch.Reset()
					if !batch.Add(file) {
						// single file exceeds the size limit. log and delete
						log.Printf("file too large. Size is '%d'", file.size)
						err = os.Remove(file.path)
						if err != nil {
							log.Printf("file too delete file. caused by '%v'", err)
						}
					}
				}
			}
		}
		time.Sleep(time.Second * time.Duration(u.options.Interval))
	}
}

func (u Uploader) upload(shots []string) {
	for _, shot := range shots {
		u.uploadOne(shot)
	}
}

func (u Uploader) uploadOne(shotPath string) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	defer func(w *multipart.Writer) {
		_ = w.Close()
	}(w)
	shotName := filepath.Dir(shotPath)
	fw, err := w.CreateFormFile("file", shotName)
	if err != nil {
		fmt.Printf("failed to create form with file '%s'. caused by: '%s'", shotName, err)
		return
	}
	shotContent, err := ioutil.ReadFile(shotPath)
	if err != nil {
		log.Printf("failed to read shot '%s'. caused by: '%v'", shotPath, err)
		return
	}
	encryptedShot, err := u.encryptor.Encrypt(shotContent)
	if err != nil {
		log.Printf("failed to encrypt shot '%s'. caused by: '%v'", shotPath, err)
		return
	}

	if _, err = io.Copy(fw, bytes.NewBuffer(encryptedShot)); err != nil {
		fmt.Printf("failed to copy contnts of shot '%s' to form. caused by: '%s'", shotPath, err)
		return
	}
	u.httpSend(&b)
}

func (u Uploader) httpSend(data *bytes.Buffer) {
	url := u.options.Url
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		fmt.Printf("failed to create POST request for url '%s'. Caused by: '%v'", url, err)
		return
	}
	req.Header.Set("Content-Type", u.contentType)
	res, err := u.client.Do(req)
	if err != nil {
		fmt.Printf("failed to upload. Caused by: '%v'", err)
	}
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		respBody, _ := ioutil.ReadAll(res.Body)
		fmt.Printf("got not ok from server: '%s'. response body is: '%s' ", res.Status, string(respBody))
	}
}
