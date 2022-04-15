package main

import (
	"log"
	"sync"
	"syncer/pkg"
)

const gmailSizeLimit = 25*1024*1024 - 10*1024

func main() {
	flags := pkg.ParseFlags()
	var wg sync.WaitGroup
	wg.Add(1)
	screenShotMaker, err := pkg.CreateScreenShotMaker(&wg, pkg.ScreenShotOptions{
		Interval:         flags.ScreenShotInterval,
		RecipientKeyPath: flags.HeaderKeyPath,
		ShotsPath:        flags.ShotsFolder,
		StorageLimit:     uint64(flags.StorageLimit),
	})
	if err != nil {
		log.Fatalf("faild to create shot maker. caused by: '%v'", err)
	}
	uploader, err := pkg.CreateUploader(pkg.UploadOptions{
		UserName:         flags.UserName,
		Password:         flags.Password,
		Url:              flags.ServerUrl,
		RecipientKeyPath: flags.ShotKeyPath,
		Interval:         flags.SyncInterval,
		BaseFolder:       flags.ShotsFolder,
		FileSuffix:       []string{".zip", pkg.ScreenShotSuffix},
		SizeLimit:        gmailSizeLimit,
	}, &wg)
	if err != nil {
		log.Fatalf("faild to create uploader. caused by: '%v'", err)
	}
	go screenShotMaker.Worker()
	go uploader.Worker()
	wg.Wait()
	log.Println("exiting")
}
