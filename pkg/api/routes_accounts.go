package api

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/multitheftauto/community/pkg/models"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (a *API) createAccount(c *gin.Context) {
	var input struct {
		Username string `json:"username" valid:"stringlength(1|255),required"`
		Password string `json:"password" valid:"stringlength(5|100),required"`
		Email    string `json:"email" valid:"email,stringlength(1|254),required"`
	}

	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	u := models.Account{
		Username: input.Username,
		Password: input.Password,
		Email:    input.Email,
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
	err = a.DB.Get(
		&count,
		"select count(id) from accounts where (username = $1) or (email = $2)",
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	u.Password = string(hashedPassword)

	_, err = a.DB.NamedExec("insert into accounts (username, password, email, activated) values (:username, :password, :email, true)", &u)

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
