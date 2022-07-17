package downloader

import (
	"fmt"
	"github.com/mrnakumar/e2g_utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

type Downloader struct {
	headerEncryptor  e2g_utils.Encryptor
	bodyDecryptor    e2g_utils.Decoder
	client           *http.Client
	downloadBasePath string
	userName         string
	password         string
	url              string
}

type DownloaderOptions struct {
	UserName             string
	Password             string
	Url                  string
	HeaderEncryptKeyPath string
	BodyDecryptKeyPath   string
	PollInterval         uint16
	DownloadBaseFolder   string
}

func MakeDownloader(options DownloaderOptions) (Downloader, error) {
	headerEncryptor, err := e2g_utils.CreateEncryptor(options.HeaderEncryptKeyPath)
	if err != nil {
		return Downloader{}, err
	}
	bodyDecryptor, err := e2g_utils.CreateDecoder(options.BodyDecryptKeyPath)
	if err != nil {
		return Downloader{}, err
	}
	return Downloader{
		headerEncryptor:  headerEncryptor,
		bodyDecryptor:    bodyDecryptor,
		client:           &http.Client{},
		downloadBasePath: options.DownloadBaseFolder,
		userName:         options.UserName,
		password:         options.Password,
	}, nil
}
func (d Downloader) Download() {
	authHeader, err := d.authHeader()
	if err != nil {
		log.Print(err)
		return
	}
	for {
		req, err := http.NewRequest("GET", d.url, nil)
		if err != nil {
			log.Println(err)
		} else {
			req.Header.Set("Authorization", authHeader)
			res, err := d.client.Do(req)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Println("Headers ", res.Header)
				if res.StatusCode == http.StatusOK {
					respBody, _ := ioutil.ReadAll(res.Body)
					filePath := path.Join(d.downloadBasePath, "already")
					err = os.WriteFile(filePath, respBody, 0644)
					if err != nil {
						log.Println(err)
					}
					_ = res.Body.Close()
				}
				if res.StatusCode == http.StatusNoContent {
					// No more to download. Print message and exit
					log.Println("Finished downloading and exiting.")
					return
				}
			}
		}
	}
}

func (d Downloader) authHeader() (string, error) {
	authHeader, err := d.headerEncryptor.Encrypt([]byte(fmt.Sprintf("%s:%s", d.userName, d.password)))
	if err != nil {
		return "", fmt.Errorf("failed to encrypt auth header to request for url '%s'. Caused by: '%v'", d.url, err)

	}
	encrypted := string(authHeader)
	encoded := e2g_utils.Base64Encode(encrypted)
	return encoded, nil
}
