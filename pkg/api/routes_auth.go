package api

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/multitheftauto/community/pkg/models"

	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (i *API) jwtAuthenticate(username string, password string, c *gin.Context) (userID uint64, success bool) {
	var u models.Account

	err := i.DB.Get(&u, "SELECT * FROM accounts WHERE username = $1", username)

	if err == sql.ErrNoRows {
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"data":    nil,
			"message": errors.Wrapf(err, "authentication query failed").Error(),
		})

		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if (err != nil) && (err != bcrypt.ErrMismatchedHashAndPassword) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"data":    nil,
			"message": errors.Wrapf(err, "could not compare hash and password").Error(),
		})

		return
	}

	return u.ID, err != bcrypt.ErrMismatchedHashAndPassword
}

func (i *API) jwtAuthorize(userId uint64, c *gin.Context) bool {
	var u models.Account

	err := i.DB.Get(&u, "SELECT * FROM accounts WHERE id = $1", userId)

	if err == sql.ErrNoRows {
		return false
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errors.Wrapf(err, "could not authorize user").Error(),
		})
	}

	c.Set("user", u)

	return u.Activated //(u.Level > 1) && (!u.Banned) && (u.Activated)
}

func (i *API) Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")

	u := models.Account{
		Username: username,
		Password: password,
		Email:    email,
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
		"select count(id) from accounts where (username = $1) or (email = $3)",
		u.Username,
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
			"message": "Account already exists with that username or email",
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

	_, err = i.DB.NamedExec("insert into accounts (username, password, email) values (:username, :password, :email)", &u)

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

func (i *API) jwtUnauthorized(c *gin.Context, code int, message string) {
	if c.Writer.Written() {
		return
	}

	c.JSON(code, gin.H{
		"status":  "error",
		"message": message,
	})
}
