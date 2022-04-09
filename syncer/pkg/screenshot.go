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

const ScreenShotSuffix = ".png"

type ScreenShotOptions struct {
	Interval         uint16
	RecipientKeyPath string
	ShotsPath        string
	StorageLimit     uint64
}
type screenShot struct {
	Name  string
	Image []byte
}

func ScreenShotMaker(wg *sync.WaitGroup, options ScreenShotOptions) {
	defer wg.Done()
	if _, err := os.Stat(options.ShotsPath); os.IsNotExist(err) {
		err := os.Mkdir(options.ShotsPath, 0644)
		log.Printf("failed to create shots path '%s'. Caused by : '%v'", options.ShotsPath, err)
		return
	}
	recipientKey, err := ioutil.ReadFile(options.RecipientKeyPath)
	if err != nil {
		log.Printf("failed to read file '%s'. Caused by : '%v'", options.RecipientKeyPath, err)
		return
	}

	recipient := string(recipientKey)
	encryptor, err := CreateEncryptor(strings.TrimSuffix(recipient, "\n"))

	if err != nil {
		log.Printf("failed to create encryptor. Caused by : '%v'", err)
		return
	}
	for {
		size, err := Size(options.ShotsPath)
		if err == nil && size < options.StorageLimit {
			shots, err := takeScreenShot("screen")
			if err != nil {
				log.Printf("failed to take shot. Caused by: '%v'", err)
			} else {
				processShots(shots, encryptor, options.ShotsPath)
			}
		} else {
			if err != nil {
				log.Printf("failed to get storage size. Caused by '%v'", err)
			} else {
				log.Printf("storage size limit reached, needs cleanup")
				cleanOld(options.ShotsPath)
			}
		}
		time.Sleep(time.Second * time.Duration(options.Interval))
	}
}

func processShots(shots []screenShot, encryptor Encryptor, shotsPath string) {
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
		name := fmt.Sprintf("%s_%d_%d%s", NamePrefix, now, randomNumber(), ScreenShotSuffix)
		screenShots = append(screenShots, screenShot{Name: name, Image: buffer.Bytes()})
	}
	return screenShots, nil
}

func cleanOld(shotsPath string) {
	files, err := ListFiles([]string{ScreenShotSuffix}, shotsPath)
	if err != nil {
		log.Printf("failed to list files. Caused by '%s'", err)
	} else {
		for _, file := range files {
			if time.Now().Sub(file.modTime) > 24*time.Hour {
				err := os.Remove(file.path)
				if err != nil {
					log.Printf("failed to delete file. Caused by '%v'", err)
				}
			}
		}
	}
}

func randomNumber() int {
	low := 10000
	high := 99999
	return low + rand.Intn(high-low)
}
