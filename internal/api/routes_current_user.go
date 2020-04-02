package api

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

func (a *API) changeCurrentUserUsername(ctx *gin.Context) {
	user := ctx.MustGet("current_user").(*User)
	var fields struct {
		Username string `json:"username"`
	}
	if err := ctx.BindJSON(&fields); err != nil {
		return
	}

	// First check that that username isn't taken...
	var count int
	err := a.DB.GetContext(ctx, &count, "select count(*) from users where username=$1", fields.Username)
	if err != nil {
		a.somethingWentWrong(ctx, err).WithField("uid", user.ID).WithField("username", fields.Username).
			Errorln("could not check if new username is in db")
		return
	}

	// ... by ensuring that the count of that username is 0
	if count != 0 {
		ctx.JSON(http.StatusConflict, gin.H{"message": "That name is taken."})
		return
	}

	_, err = a.QB.Update("users").Set("username", fields.Username).Where(squirrel.Eq{"id": user.ID}).ExecContext(ctx)
	if err != nil {
		a.somethingWentWrong(ctx, err).WithField("uid", user.ID).WithField("username", fields.Username).
			Errorln("could not update new username in db")
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (a *API) changeCurrentUserPassword(ctx *gin.Context) {
	user := ctx.MustGet("current_user").(*User)
	var fields struct {
		Password    string `json:"password"`
		NewPassword string `json:"new_password"`
	}
	if err := ctx.BindJSON(&fields); err != nil {
		return
	}

	// First check that password matches
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(fields.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Incorrect username or password"})
		return
	} else if err != nil {
		a.somethingWentWrong(ctx, err).WithField("uid", user.ID).Errorln("could not verify provided password")
		return
	}

	// Generate new hash
	password, err := bcrypt.GenerateFromPassword([]byte(fields.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		a.somethingWentWrong(ctx, err).WithField("uid", user.ID).Errorln("could not generate new password hash")
		return
	}

	_, err = a.QB.Update("users").Set("password", password).Where(squirrel.Eq{"id": user.ID}).ExecContext(ctx)
	if err != nil {
		a.somethingWentWrong(ctx, err).WithField("uid", user.ID).Errorln("could not update new password in db")
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (a *API) deleteCurrentUser(ctx *gin.Context) {
	user := ctx.MustGet("current_user").(*User)
	// todo: check password or sudo mode

	_, err := a.QB.Delete("users").Where(squirrel.Eq{"id": user.ID}).ExecContext(ctx)
	if err != nil {
		a.somethingWentWrong(ctx, err).WithField("uid", user.ID).Errorln("could not delete account")
		return
	}

	ctx.Status(http.StatusNoContent)
}
