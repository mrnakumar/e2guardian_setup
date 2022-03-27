package pkg

import (
	"bytes"
	"fmt"
	"github.com/kbinani/screenshot"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type screenShot struct {
	Name  string
	Image []byte
}

func ScreenShotMaker(wg *sync.WaitGroup, interval uint16, recipientKeyPath string, shotsPath string) {
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
	encryptor, err := CreateEncryptor(strings.TrimSuffix(recipient, "\n"))

	if err != nil {
		log.Fatalf("failed to create encryptor. Caused by : '%v'", err)
	}
	for {
		shots, err := takeScreenShot("screen")
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

func takeScreenShot(NamePrefix string) ([]screenShot, error) {
	n := screenshot.NumActiveDisplays()
	var screenShots []screenShot
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			return nil, err
		}
		var buffer bytes.Buffer
		err = png.Encode(&buffer, img)
		if err != nil {
			return nil, err
		}
		now := time.Now().UnixMilli()
		name := fmt.Sprintf("%s_%d_%d.png", NamePrefix, now, randomNumber())
		screenShots = append(screenShots, screenShot{Name: name, Image: buffer.Bytes()})
	}
	return screenShots, nil
}

func randomNumber() int {
	low := 10000
	high := 99999
	return low + rand.Intn(high-low)
}
