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
	go pkg.ScreenShotMaker(&wg, pkg.ScreenShotOptions{
		Interval:         flags.ScreenShotInterval,
		RecipientKeyPath: flags.HeaderKeyPath,
		ShotsPath:        flags.ShotsFolder,
		StorageLimit:     uint64(flags.StorageLimit),
	})
	uploader := pkg.CreateUploader(pkg.UploadOptions{
		UserName:   flags.UserName,
		Password:   flags.Password,
		Url:        flags.ServerUrl,
		Interval:   flags.SyncInterval,
		BaseFolder: flags.ShotsFolder,
		FileSuffix: []string{".zip", pkg.ScreenShotSuffix},
		SizeLimit:  gmailSizeLimit,
	}, &wg)
	go uploader.Worker()
	wg.Wait()
	log.Println("exiting")
}
