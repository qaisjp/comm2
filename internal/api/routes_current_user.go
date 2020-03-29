package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *API) getCurrentUser(c *gin.Context) {
	user := c.MustGet("current_user").(*User)
	// todo: support extra private info
	c.JSON(http.StatusOK, user.PrivateInfo())
}
