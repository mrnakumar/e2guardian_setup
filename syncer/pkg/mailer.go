package pkg

import (
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type MailOptions struct {
	From       string
	To         string
	Password   string
	Host       string
	Port       int
	Subject    string
	Interval   uint16
	BaseFolder string
	FileSuffix []string
	SizeLimit  int64
	filePath   string
}

type fileInfo struct {
	path string
	size int64
}

type batch struct {
	files []string
	size  int64
	limit int64
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
func Mailer(wg *sync.WaitGroup, options MailOptions) {
	defer wg.Done()
	for {
		files, err := listFiles(options.FileSuffix, options.BaseFolder)
		if err != nil {
			log.Printf("failed to list files. Caused by '%s'", err)
		} else {
			batch := batch{
				files: make([]string, 0),
				size:  0,
				limit: options.SizeLimit,
			}
			for _, file := range files {
				if !batch.Add(file) {
					sendMail(options, batch.files)
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
			if batch.size > 0 {
				sendMail(options, batch.files)
				batch.Reset()
			}
		}
		time.Sleep(time.Second * time.Duration(options.Interval))
	}
}

func sendMail(options MailOptions, attachments []string) {
	d := gomail.NewDialer(options.Host, options.Port, options.From, options.Password)
	s, err := d.Dial()
	if err != nil {
		log.Printf("failed to connect to gmail. Calused by '%v'", err)
		return
	}
	m := gomail.NewMessage()
	m.SetHeader("From", options.From)
	m.SetHeader("To", options.To)
	m.SetHeader("Subject", options.Subject)
	for _, attachment := range attachments {
		m.Attach(attachment)
	}
	if err := gomail.Send(s, m); err != nil {
		log.Printf("could not send email. Caused by '%v'", err)
	}
}

func listFiles(suffixes []string, basePath string) ([]fileInfo, error) {
	infos, err := ioutil.ReadDir(basePath)
	if err != nil {
		return nil, err
	}
	var files []fileInfo
	for _, info := range infos {
		if info.Size() > 0 && matchSuffix(suffixes, info.Name()) {
			files = append(files, fileInfo{path: filepath.Join(basePath, info.Name()), size: info.Size()})
		}
	}
	return files, nil
}

func matchSuffix(suffixes []string, fileName string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(fileName, suffix) {
			return true
		}
	}
	return false
}
