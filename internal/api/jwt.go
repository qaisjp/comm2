package api

import (
	"errors"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/multitheftauto/community/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (a *API) jwtUnauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"message": message,
	})
}

func (a *API) jwtAuthorizator(data interface{}, c *gin.Context) bool {
	return true
}

func (a *API) jwtAuthenticator(c *gin.Context) (_ interface{}, err error) {
	var input struct {
		Username string `valid:"stringlength(1|255),required"`
		Password string `valid:"required"`
	}

	if err = c.BindJSON(&input); err != nil {
		return "", err
	}

	if input.Username == "" || input.Password == "" {
		return "", errors.New("missing username or password")
	}

	success, err := govalidator.ValidateStruct(&input)
	if !success {
		return "", err
	}

	var user models.User

	err = a.DB.Get(&user, "select id, password, is_activated from users where (email = $1) or (username = $1) limit 1", input.Username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return "", err
	}

	if !user.Activated {
		return "", errors.New("user not activated")
	}

	return &user, nil
}

func (a *API) jwtIdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)

	return int(claims["id"].(float64))
}

func (a *API) jwtPayloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*models.User); ok {
		return jwt.MapClaims{
			"id": v.ID,
		}
	}
	return jwt.MapClaims{}
}
