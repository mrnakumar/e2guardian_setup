package pkg

import (
	"github.com/gin-gonic/gin"
	"github.com/mrnakumar/e2g_utils"
	"strings"
)

const AuthError = "authError"

type AuthChecker struct {
	Decoder  e2g_utils.Decoder
	UserName string
	Password string
}

func (a AuthChecker) AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		decodedHeader, err := e2g_utils.Base64Decode(authHeader)
		if err != nil {
			c.Set(AuthError, true)
		} else {
			authHeaderDecoded, err := a.Decoder.Decrypt(decodedHeader)
			if err != nil {
				c.Set(AuthError, true)
				return
			}
			headerParts := strings.Split(string(authHeaderDecoded), ":")
			key, pwd := headerParts[0], headerParts[1]
			if a.UserName != key || a.Password != pwd {
				c.Set(AuthError, true)
			}
		}
		c.Next()
	}
}
