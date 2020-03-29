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

func (a *API) followUser(ctx *gin.Context) {
	user := ctx.MustGet("current_user").(*User)
	targetUser := ctx.MustGet("target_user").(*User)
	method := ctx.Request.Method

	if user.ID == targetUser.ID && method == "PUT" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "You can't follow yourself."})
		return
	}

	var err error
	if method == "GET" {
		var count int
		err = a.DB.GetContext(ctx, &count, "select count(*) from user_followings where source_user_id = $1 and target_user_id = $2", user.ID, targetUser.ID)
		if err == nil && count == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Not following"})
		}
	} else if method == "PUT" {
		_, err = a.DB.ExecContext(ctx, "insert into user_followings(source_user_id, target_user_id) values ($1, $2) on conflict do nothing", user.ID, targetUser.ID)
	} else if method == "DELETE" {
		_, err = a.DB.ExecContext(ctx, "delete from user_followings where source_user_id = $1 and target_user_id = $2", user.ID, targetUser.ID)
	}

	if err != nil {
		a.somethingWentWrong(ctx, err).WithField("source_id", user.ID).WithField("target_id", targetUser.ID).WithField("methid", method).Errorln("could not action following")
		return
	}

	// fmt.Println(ok, user.ID, ctx.Request.Method, targetUser.ID)
	ctx.Status(http.StatusNoContent)
}
