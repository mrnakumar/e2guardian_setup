package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"syncer/pkg"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go screenShotMaker(&wg, 10, "/keys/pub",
		"/syncer/shots")
	wg.Wait()
	log.Println("exiting")
}

func screenShotMaker(wg *sync.WaitGroup, interval uint64, recipientKeyPath string, shotsPath string) {
	defer wg.Done()
	if _, err := os.Stat(shotsPath); os.IsNotExist(err) {
		err := os.Mkdir(shotsPath, 0644)
		log.Fatalf("failed to create shots path '%s'. Caused by : '%v'", shotsPath, err)
	}
	recipientKey, err := ioutil.ReadFile(recipientKeyPath)
	if err != nil {
		log.Fatalf("failed to read file '%s'. Caused by : '%v'", recipientKeyPath, err)
	}

	recipient := string(recipientKey)
	encryptor, err := pkg.CreateEncryptor(strings.TrimSuffix(recipient, "\n"))

	if err != nil {
		log.Fatalf("failed to create encryptor. Caused by : '%v'", err)
	}
	for {
		shots, err := pkg.TakeScreenShot("screen")
		if err != nil {
			log.Printf("failed to take shot. Caused by: '%v'", err)
		} else {
			for _, shot := range shots {
				encrypted, err := encryptor.Encrypt(shot.Image)
				if err == nil {
					shotPathInProgress := path.Join(shotsPath, shot.Name+".progress")
					err := os.WriteFile(shotPathInProgress, encrypted, 0644)
					if err != nil {
						log.Printf("failed to save shot. Caused by: '%v'", err)
					}
					err = os.Rename(shotPathInProgress, path.Join(shotsPath, shot.Name))
					if err != nil {
						log.Printf("failed to rename shot '%s'. Caused by '%v", shotPathInProgress, err)
					}
				} else {
					log.Printf("failed to encrypt shot '%s'", shot.Name)
				}
			}
		}
		time.Sleep(time.Second * time.Duration(interval))
	}
}
