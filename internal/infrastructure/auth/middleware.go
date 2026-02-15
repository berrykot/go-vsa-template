package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const ContextKeyUser = "user"

func GetAuthUser(c *gin.Context) (*UserInfo, bool) {
	u, ok := c.Get(ContextKeyUser)
	if !ok {
		return nil, false
	}

	userInfo, ok := u.(*UserInfo)
	if !ok || userInfo == nil {
		return nil, false
	}

	return userInfo, true
}

// Middleware
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//your auth -> set userInfo to context
		if false {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth"})
		}
		userInfo := UserInfo{}
		c.Set(ContextKeyUser, userInfo)
		c.Next()
	}
}
