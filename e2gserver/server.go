package main

import (
	"e2gserver/pkg"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/mrnakumar/e2g_utils"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	flags := pkg.ParseFlags()
	route := gin.Default()
	decoder, err := e2g_utils.CreateDecoder(flags.IdentityFilePath)
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
	authChecker := pkg.AuthChecker{
		Decoder:  decoder,
		UserName: flags.UserName,
		Password: flags.Password,
	}
	route.POST("/shoot", authChecker.AuthCheck(), pkg.FileHandler(flags.BasePath))
	route.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "Hi")
	})
	route.POST("/eat", authChecker.AuthCheck(), func(c *gin.Context) {
		if _, failed := c.Get(pkg.AuthError); failed {
			c.String(http.StatusUnauthorized, "")
			return
		}
		files, err := e2g_utils.ListFiles([]string{".ec"}, flags.BasePath)
		if err != nil {
			log.Printf("failed to list files. Caused by '%s'", err)
			c.String(500, "")
		} else {
			if len(files) > 0 {
				file := files[0]
				c.Header("File-Name", filepath.Base(file.Path))
				c.File(file.Path)
				err = os.Remove(file.Path)
				if err != nil {
					log.Printf("failed to delete file after giving to user. Caused by '%v'", err)
				}
			} else {
				c.String(204, "")
			}
		}
	})
	if flags.DevelopMode {
		log.Fatal(route.Run("localhost:8080"))
	} else {
		//log.Fatal(autotls.Run(route, flags.Domain))
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(flags.Domain),
			Cache:      autocert.DirCache("/var/www/.cache"),
		}

		log.Fatal(autotls.RunWithManager(route, &m))
	}
}
