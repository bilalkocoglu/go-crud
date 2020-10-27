package mw

import (
	"encoding/base64"
	_const "github.com/bilalkocoglu/go-crud/pkg/const"
	"github.com/bilalkocoglu/go-crud/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request
		basicToken := req.Header.Get(_const.AuthorizationHeader)
		if basicToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token must be not null.",
			})
			return
		}
		tokenType, token := parseToken(basicToken)
		if tokenType != "Basic" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token type must be Basic.",
			})
			return
		}
		decodeToken, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token decode exception",
			})
			return
		}
		splitToken := strings.Split(string(decodeToken), ":")
		username := splitToken[0]
		password := splitToken[1]

		var user database.User
		err = database.GetUserByUsername(&user, username)
		if err != nil || user.ID == 0 || user.Password != password {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid username or password",
			})
			return
		}
		c.Set(_const.CurrentUser, user)

		c.Next()
	}
}

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request

		bearerToken := req.Header.Get(_const.AuthorizationHeader)

		if bearerToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token must be not null.",
			})
			return
		}

		tokenType, token := parseToken(bearerToken)

		if tokenType != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token type must be Bearer.",
			})
			return
		}

		log.Info().Msg("tokenType: " + tokenType + " token: " + token)

		c.Next()
	}
}

func parseToken(token string) (string, string) {
	parsedToken := strings.Split(token, " ")

	return parsedToken[0], parsedToken[1]
}
