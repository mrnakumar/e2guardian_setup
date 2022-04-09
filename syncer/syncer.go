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
		RecipientKeyPath: flags.KeyPath,
		ShotsPath:        flags.ShotsFolder,
		StorageLimit:     uint64(flags.StorageLimit),
	})
	go pkg.Uploader(&wg, pkg.MailOptions{
		From:       flags.FromEmail,
		To:         flags.ToEmail,
		Password:   flags.Password,
		Host:       "smtp.gmail.com",
		Port:       25,
		Subject:    "SHOTS",
		Interval:   flags.SyncInterval,
		BaseFolder: flags.ShotsFolder,
		FileSuffix: []string{".zip", pkg.ScreenShotSuffix},
		SizeLimit:  gmailSizeLimit,
	})
	wg.Wait()
	log.Println("exiting")
}
