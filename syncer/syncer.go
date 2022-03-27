package main

import (
	"log"
	"sync"
	"syncer/pkg"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go pkg.ScreenShotMaker(&wg, 10, "/keys/pub",
		"/syncer/shots")
	wg.Wait()
	log.Println("exiting")
}
