package main

import (
	"flag"
	"github.com/mrnakumar/downloader/downloader"
	"log"
)

func main() {
	userName := flag.String("userName", "", "User name to authenticate to server")
	password := flag.String("password", "", "Password to connec to server")
	url := flag.String("url", "", "Server url to download screenshots")
	headerEncryptionKeyPath := flag.String("headerEncryptKeyPath", "", "Path to file that has public key to encrypt auth header")
	bodyDecodeKeyPath := flag.String("bodyDecodeKeyPath", "", "Path to file that has private key to decode screenshot")
	downloadFolder := flag.String("downloadPath", "", "Directory path to save screenshots")
	flag.Parse()
	if isEmpty(userName) || isEmpty(password) || isEmpty(url) || isEmpty(headerEncryptionKeyPath) || isEmpty(bodyDecodeKeyPath) || isEmpty(downloadFolder) {
		flag.PrintDefaults()
		return
	}
	scDownloader, err := downloader.MakeDownloader(downloader.DownloaderOptions{
		UserName:             *userName,
		Password:             *password,
		Url:                  *url,
		HeaderEncryptKeyPath: *headerEncryptionKeyPath,
		BodyDecryptKeyPath:   *bodyDecodeKeyPath,
		DownloadBaseFolder:   *downloadFolder,
	})
	if err != nil {
		log.Printf("Failed to create downlader: %v\n", err)
		return
	}
	scDownloader.Download()
	log.Println("Finished downloading all available")
}

func isEmpty(input *string) bool {
	return len(*input) == 0
}
