package auth

import (
	"net/http"

	"github.com/multitheftauto/community/pkg/models"
	"github.com/pkg/errors"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"golang.org/x/crypto/bcrypt"
)

func (i *Impl) Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")

	u := models.Account{
		Username: username,
		Password: password,
		Email:    email,
		Slug:     slug.Make(username),
	}

	success, err := govalidator.ValidateStruct(&u)
	if !success {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	var count int
	err = i.DB.Get(
		&count,
		"select count(id) from accounts where (username = $1) or (slug = $2) or (email = $3)",
		u.Username,
		u.Slug,
		u.Email,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errors.Wrapf(err, "could not check existence").Error(),
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"status":  "error",
			"message": "Account already exists with that username, slug, or email",
		})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	u.Password = string(hashedPassword)

	_, err = i.DB.NamedExec("insert into accounts (username, password, email, slug) values (:username, :password, :email, :slug)", &u)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errors.Wrap(err, "could not insert").Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
	})
}
