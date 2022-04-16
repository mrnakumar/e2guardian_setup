package main

import (
	"log"
	"sync"
	"syncer/pkg"
)

const fileUploadSizeLimit = 25*1024*1024 - 10*1024

func main() {
	flags := pkg.ParseFlags()
	var wg sync.WaitGroup
	wg.Add(1)
	screenShotMaker, err := pkg.CreateScreenShotMaker(&wg, pkg.ScreenShotOptions{
		Interval:      flags.ScreenShotInterval,
		ShotKeyPath: flags.ShotKeyPath,
		ShotsPath:     flags.ShotsFolder,
		StorageLimit:  uint64(flags.StorageLimit),
	})
	if err != nil {
		log.Fatalf("faild to create shot maker. caused by: '%v'", err)
	}
	uploader, err := pkg.CreateUploader(pkg.UploadOptions{
		UserName:      flags.UserName,
		Password:      flags.Password,
		Url:           flags.ServerUrl,
		HeaderKeyPath: flags.HeaderKeyPath,
		Interval:      flags.SyncInterval,
		BaseFolder:    flags.ShotsFolder,
		FileSuffix:    []string{".zip", pkg.ScreenShotSuffix},
		SizeLimit:     fileUploadSizeLimit,
	}, &wg)
	if err != nil {
		log.Fatalf("faild to create uploader. caused by: '%v'", err)
	}
	go screenShotMaker.Worker()
	go uploader.Worker()
	wg.Wait()
	log.Println("exiting")
}
