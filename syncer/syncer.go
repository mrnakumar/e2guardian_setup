package main

import (
	"log"
	"sync"
	"syncer/pkg"
)

func main() {
	flags := pkg.ParseFlags()
	var wg sync.WaitGroup
	wg.Add(1)
	go pkg.ScreenShotMaker(&wg, flags.ScreenShotInterval, flags.KeyPath, flags.ShotsFolder)
	go pkg.Mailer(&wg, pkg.MailOptions{
		From:       flags.FromEmail,
		To:         flags.ToEmail,
		Password:   flags.Password,
		Host:       "smtp.gmail.com",
		Port:       25,
		Subject:    "SHOTS",
		Interval:   10,
		BaseFolder: flags.ShotsFolder,
		FileSuffix: []string{".zip", ".png"},
		SizeLimit:  30000,
	})
	wg.Wait()
	log.Println("exiting")
}
// 25*1024*1024 - 10*1024
