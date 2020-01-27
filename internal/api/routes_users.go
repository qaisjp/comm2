package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (a *API) checkUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return
	}

	// Check if the user exists
	var user User
	if err := a.DB.Get(&user, "select * from users where id = $1", userID); err != nil {
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
