package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func WatchSessionCheck(c *gin.Context) {
	cookie, err := c.Cookie("prinflix_session_token")
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(403)
		return
	}

	fmt.Println(cookie)

	c.Next()
}