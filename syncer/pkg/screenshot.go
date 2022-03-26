package pkg

import (
	"bytes"
	"fmt"
	"github.com/kbinani/screenshot"
	"image/png"
	"math/rand"
	"time"
)

type ScreenShot struct {
	Name  string
	Image []byte
}

func TakeScreenShot(NamePrefix string) ([]ScreenShot, error) {
	n := screenshot.NumActiveDisplays()
	var screenShots []ScreenShot
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
		screenShots = append(screenShots, ScreenShot{Name: name, Image: buffer.Bytes()})
	}
	return screenShots, nil
}

func randomNumber() int {
	low := 10000
	high := 99999
	return low + rand.Intn(high-low)
}
