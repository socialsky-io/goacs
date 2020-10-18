package auth

import "github.com/gin-gonic/gin"

func ACSBasicAuth() gin.HandlerFunc {
	realm := "GoACS Authorization"

	return func(c *gin.Context) {

	}
}
