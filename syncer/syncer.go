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
	wg.Wait()
	log.Println("exiting")
}
