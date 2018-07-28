package auth

import "github.com/gin-gonic/gin"

func (i *Impl) Unauthorized(c *gin.Context, code int, message string) {
	if c.Writer.Written() {
		return
	}

	c.JSON(code, gin.H{
		"status":  "error",
		"message": message,
	})
}
