package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (a *API) parseUserID(param string, key string, mustParse bool) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var fieldVal interface{} = ctx.Param(param)
		fieldName := "username"
		userID, err := strconv.ParseUint(ctx.Param(param), 10, 64)
		if err == nil {
			fieldVal = userID
			fieldName = "id"
		}

		// Initial key type
		var user *User
		ctx.Set(key, user)

		// Check if the user exists
		user = &User{}
		if err := a.DB.Get(user, "select * from users where "+pq.QuoteIdentifier(fieldName)+" = $1", fieldVal); err != nil {
			if err == sql.ErrNoRows {
				if mustParse {
					ctx.JSON(http.StatusNotFound, gin.H{
						"message": "That user could not be found",
					})
					ctx.Abort()
				}
				return
			}

			a.Log.WithField("err", err).Errorln("Failed to find user")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to find user",
			})
			ctx.Abort()
			return
		}

		// Store the resource
		ctx.Set(key, user)
	}
}

func (a *API) checkUser(c *gin.Context) {
	var fieldVal interface{} = c.Param("user_id")
	fieldName := "username"
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err == nil {
		fieldVal = userID
		fieldName = "id"
	}

	// Check if the user exists
	var user User
	if err := a.DB.Get(&user, "select * from users where "+pq.QuoteIdentifier(fieldName)+" = $1", fieldVal); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "That user could not be found",
			})
			c.Abort()
			return
		}

		a.Log.WithField("err", err).Errorln("Failed to find user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to find user",
		})
		c.Abort()
		return
	}

	// Store the resource
	c.Set("user", &user)
}

func (a *API) getUser(c *gin.Context) {
	user := c.MustGet("user").(*User)

	// todo: support extra private info
	// currentUser := c.MustGet("current_user").(*User)
	// if user.ID == currentUser.ID {
	// 	c.JSON(http.StatusOK, user.PrivateInfo())
	// 	return
	// }

	c.JSON(http.StatusOK, user.PublicInfo())
}

func (a *API) getUserProfile(ctx *gin.Context) {
	user := ctx.MustGet("user").(*User)
	elevated := false
	if currentUser := ctx.MustGet("current_user").(*User); currentUser != nil {
		elevated = user.ID == currentUser.ID
	}

	var result interface{}

	pu := user.PublicInfo()

	resources := []Resource{}
	{
		resourceQuery := "select r.* from resources r, resource_collaborators c where (r.author_id = $1) or (r.id = c.resource_id and c.user_id = $1) and c.accepted"
		if !elevated {
			resourceQuery += " and status = 'public'"
		}
		if err := a.DB.SelectContext(ctx, &resources, resourceQuery, user.ID); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong."})
			ctx.Error(err)
			a.Log.WithError(err).Errorln("could not get resources for profile info")
			return
		}
	}

	following, err := user.GetFollowing(ctx, a.DB)
	if err != nil {
		a.somethingWentWrong(ctx, err).WithField("user_id", user.ID).Errorln("could not get this user's following list")
		return
	}
	followers, err := user.GetFollowers(ctx, a.DB)
	if err != nil {
		a.somethingWentWrong(ctx, err).WithField("user_id", user.ID).Errorln("could not get this user's followers")
		return
	}

	type BaseProfileInfo struct {
		PublicUserInfo
		Resources []Resource       `json:"resources"`
		Following []PublicUserInfo `json:"following"`
		Followers []PublicUserInfo `json:"followers"`
	}
	base := BaseProfileInfo{
		PublicUserInfo: pu,
		Resources:      resources,
		Following:      UserSlice(following).PublicInfo(),
		Followers:      UserSlice(followers).PublicInfo(),
	}

	if elevated {
		followersMap := make(map[uint64]struct{})
		for _, f := range base.Followers {
			followersMap[f.ID] = struct{}{}
		}
		for i, f := range base.Following {
			_, ok := followersMap[f.ID]
			base.Following[i].FollowsYou = &ok
		}

		result = struct {
			BaseProfileInfo
		}{base}
	} else {
		result = struct {
			BaseProfileInfo
		}{base}
	}

	ctx.JSON(http.StatusOK, result)
}

func (a *API) createUser(c *gin.Context) {
	var input struct {
		Username string `json:"username" valid:"stringlength(1|255),required"`
		Password string `json:"password" valid:"stringlength(5|100),required"`
		Email    string `json:"email" valid:"email,stringlength(1|254),required"`
	}

	if err := c.BindJSON(&input); err != nil {
		return
	}

	u := User{
		Username: input.Username,
		Password: input.Password,
		Email:    input.Email,
	}

	success, err := govalidator.ValidateStruct(&u)
	if !success {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	var count int
	err = a.DB.Get(
		&count,
		"select count(id) from users where (username = $1) or (email = $2)",
		u.Username,
		u.Email,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": errors.Wrapf(err, "could not check existence").Error(),
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"message": "User already exists with that username or email",
		})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	u.Password = string(hashedPassword)

	_, err = a.DB.NamedExec("insert into users (username, password, email, is_activated) values (:username, :password, :email, true)", &u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": errors.Wrap(err, "could not insert").Error(),
		})

		return
	}

	c.Status(http.StatusCreated)
}

func (a *API) getUserFollowers(ctx *gin.Context) {
	user := ctx.MustGet("user").(*User)
	rows, err := user.GetFollowers(ctx, a.DB)
	if err != nil {
		a.somethingWentWrong(ctx, err).WithField("user_id", user.ID).Errorln("could not get this user's followers")
		return
	}
	ctx.JSON(http.StatusOK, UserSlice(rows).PublicInfo())
}

func (a *API) getUserFollowing(ctx *gin.Context) {
	user := ctx.MustGet("user").(*User)
	rows, err := user.GetFollowers(ctx, a.DB)
	if err != nil {
		a.somethingWentWrong(ctx, err).WithField("user_id", user.ID).Errorln("could not get this user's following list")
		return
	}
	ctx.JSON(http.StatusOK, UserSlice(rows).PublicInfo())
}
