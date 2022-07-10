package main

import (
	"e2gserver/pkg"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/mrnakumar/e2g_utils"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	flags := pkg.ParseFlags()
	route := gin.Default()
	decoder, err := pkg.CreateDecoder(flags.IdentityFilePath)
	if err != nil {
		log.Fatalf("failed to create decoder. Reason: '%v'", err)
	}
	go func() {
		// clean old
		for {
			files, err := e2g_utils.ListFiles([]string{".ec"}, flags.BasePath)
			if err != nil {
				log.Printf("failed to list files. Caused by '%s'", err)
			} else {
				for _, file := range files {
					if time.Now().Sub(file.ModTime) > 48*time.Hour {
						err := os.Remove(file.Path)
						if err != nil {
							log.Printf("failed to delete file. Caused by '%v'", err)
						}
					}
				}
			}
			time.Sleep(10 * time.Minute)
		}
	}()
	route.POST("/shoot", pkg.FileHandler(decoder, flags.BasePath, flags.UserName, flags.Password))
	route.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "Hi")
	})
	if flags.DevelopMode {
		log.Fatal(route.Run("localhost:8080"))
	} else {
		log.Fatal(autotls.Run(route, flags.Domain))
	}
}
