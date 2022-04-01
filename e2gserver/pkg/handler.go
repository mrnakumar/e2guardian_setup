package pkg

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func FileHandler(decoder Decoder, basePath string, eKey string, ePassword string) func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		authHeaderDecoded, err := decoder.Decrypt(authHeader)
		if err != nil {
			c.String(http.StatusUnauthorized, "Invalid auth")
			return
		}
		headerParts := strings.Split(string(authHeaderDecoded), ":")
		key, pwd := headerParts[0], headerParts[1]
		if eKey != key || ePassword != pwd {
			c.String(http.StatusUnauthorized, "")
			return
		}
		file, _ := c.FormFile("file")
		if file == nil {
			log.Printf("missing file")
			c.String(http.StatusBadRequest, "missing file")
			return
		}
		now := time.Now().UnixMilli()
		dstFileName := fmt.Sprintf("%s_%d_%d", file.Filename, now, getRandomNumber())
		dstPath := filepath.Join(basePath, dstFileName)
		err = c.SaveUploadedFile(file, dstPath)
		if err != nil {
			log.Printf("failed to save file '%s'. Caused by '%v'", file.Filename, err)
			c.String(http.StatusInternalServerError, "")
			return
		}
		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	}
}

func getRandomNumber() int {
	max := 9999999999
	return rand.Intn(max)
}
