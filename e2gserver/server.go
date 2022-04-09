package main

import (
	"e2gserver/pkg"
	"encoding/base64"
	"filippo.io/age"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

/*import (
	"e2gserver/pkg"
	"encoding/base64"
	"filippo.io/age"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)*/

func main() {
	flags := pkg.ParseFlags()
	route := gin.Default()
	decoder, err := pkg.CreateDecoder(flags.IdentityFilePath)
	if err != nil {
		log.Fatalf("failed to create decoder. Reason: '%v'", err)
	}
	route.POST("/shoot", pkg.FileHandler(decoder, flags.BasePath, flags.UserName, flags.Password))
	route.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "Hi")
	})
	log.Fatal(autotls.Run(route, flags.Domain))
}

func generateId() {
	id, _ := age.GenerateX25519Identity()
	keyPri := id.String()
	keyPub := id.Recipient().String()
	println(base64.StdEncoding.EncodeToString([]byte(keyPri)))
	println(base64.StdEncoding.EncodeToString([]byte(keyPub)))
}
