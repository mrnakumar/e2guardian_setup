package main

import (
	"e2gserver/pkg"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	flags := pkg.ParseFlags()
	route := gin.Default()
	route.POST("/shot", pkg.FileHandler(flags.BasePath, flags.UserName, flags.Password))
	route.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "Hi")
	})
	log.Fatal(autotls.Run(route, flags.Domain))
}
