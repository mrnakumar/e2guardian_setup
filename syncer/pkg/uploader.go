package pkg

import (
	"bytes"
	"fmt"
	"github.com/mrnakumar/e2g_utils"
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
	UserName      string
	Password      string
	Url           string
	HeaderKeyPath string
	Interval      uint16
	BaseFolder    string
	FileSuffix    []string
	SizeLimit     int64
	filePath      string
}

type Uploader struct {
	options     UploadOptions
	wg          *sync.WaitGroup
	client      *http.Client
	contentType string
	encryptor   e2g_utils.Encryptor
}

func CreateUploader(options UploadOptions, wg *sync.WaitGroup) (Uploader, error) {
	encryptor, err := e2g_utils.CreateEncryptor(options.HeaderKeyPath)
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

func (u Uploader) Worker() {
	defer u.wg.Done()
	for {
		files, err := e2g_utils.ListFiles(u.options.FileSuffix, u.options.BaseFolder)
		if err != nil {
			log.Printf("failed to list files. Caused by '%s'", err)
		} else {
			for _, file := range files {
				if file.Size > u.options.SizeLimit {
					log.Printf("the file '%s' is larger than allowed size '%d' large. deleting it.", file.Path, u.options.SizeLimit)
					err = os.Remove(file.Path)
					if err != nil {
						log.Printf("filed to delete file. caused by '%v'", err)
					}
				} else {
					u.uploadOne(file.Path)
				}
			}
		}
		time.Sleep(time.Second * time.Duration(u.options.Interval))
	}
}

func (u Uploader) uploadOne(shotPath string) {
	b, contentType, ok := u.makeBody(shotPath)
	if ok {
		u.httpSend(b, contentType)
	}
}

func (u Uploader) makeBody(shotPath string) (*bytes.Buffer, string, bool) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	defer func(w *multipart.Writer) {
		_ = w.Close()
	}(w)

	shotContent, err := ioutil.ReadFile(shotPath)
	if err != nil {
		log.Printf("failed to read shot '%s'. caused by: '%v'", shotPath, err)
		return &bytes.Buffer{}, "", false
	}
	shotName := filepath.Dir(shotPath)
	fw, err := w.CreateFormFile("file", shotName)
	if err != nil {
		log.Printf("failed to create form with file '%s'. caused by: '%s'", shotName, err)
		return &bytes.Buffer{}, "", false
	}

	if _, err = io.Copy(fw, bytes.NewBuffer(shotContent)); err != nil {
		log.Printf("failed to copy contnts of shot '%s' to form. caused by: '%s'", shotPath, err)
		return &bytes.Buffer{}, "", false
	}

	return &b, w.FormDataContentType(), true
}

func (u Uploader) httpSend(data *bytes.Buffer, contentType string) {
	url := u.options.Url
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		log.Printf("failed to create POST request for url '%s'. Caused by: '%v'", url, err)
		return
	}
	req.Header.Set("Content-Type", contentType)
	authHeader, err := u.encryptor.Encrypt([]byte(fmt.Sprintf("%s:%s", u.options.UserName, u.options.Password)))
	if err != nil {
		log.Printf("failed to encrypt auth header to request for url '%s'. Caused by: '%v'", url, err)
		return
	}
	encrypted := string(authHeader)
	encoded := e2g_utils.Base64Encode(encrypted)
	req.Header.Set("Authorization", encoded)
	res, err := u.client.Do(req)
	if err != nil {
		log.Printf("failed to upload. Caused by: '%v'", err)
		return
	}
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		respBody, _ := ioutil.ReadAll(res.Body)
		log.Printf("got not ok from server: '%s'. response body is: '%s' ", res.Status, string(respBody))
	}
}
