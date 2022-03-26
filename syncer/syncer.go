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
	_, err = pkg.CreateEncryptor(strings.TrimSuffix(recipient, "\n"))

	if err != nil {
		log.Fatalf("failed to create encryptor. Caused by : '%v'", err)
	}
	for {
		shots, err := pkg.TakeScreenShot("screen")
		if err != nil {
			log.Printf("failed to take shot. Caused by: '%v'", err)
		} else {
			for _, shot := range shots {
				log.Printf("Name '%s'", shot.Name)

				//encrypted, err := encryptor.Encrypt(shot.Image)
				encrypted := shot.Image
				if err == nil {
					shotPath := path.Join(shotsPath, shot.Name)
					log.Printf("shot path '%s'. Name '%s'", shotPath, shot.Name)
					err := os.WriteFile(shotPath, encrypted, 0644)
					if err != nil {
						log.Printf("failed to save shot. Caused by: '%v'", err)
					}
				} else {
					log.Printf("failed to encrypt shot '%s'", shot.Name)
				}
			}
		}
		time.Sleep(time.Second * time.Duration(interval))
	}
}
