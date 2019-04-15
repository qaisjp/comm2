package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/multitheftauto/community/internal/models"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	jwt "gopkg.in/dgrijalva/jwt-go.v3"
)

func (a *API) oauthToken(c *gin.Context) {
	var input struct {
		GrantType  string `json:"grant_type"`
		Username   string `json:"username"`
		Password   string `json:"password"`
		ExpiryTime uint64 `json:"expiry_time"`
	}
	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if input.ExpiryTime < 0 {
		panic("input.ExpiryTime should be uint64")
	}

	switch input.GrantType {
	case "password":
		// Parameters:
		//  - username    - username/email in the system
		//  - password    - the account password
		//  - expiry_time - seconds until token expires

		if input.ExpiryTime == 0 {
			input.ExpiryTime = 86400 // 24 hours
		}

		acc := models.Account{}
		err := a.DB.Get(&acc, "SELECT * FROM accounts WHERE (username = $1) or (email = $1)", input.Username)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, &gin.H{
				"status":  "error",
				"message": "No such address exists",
			})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, &gin.H{
				"status":  "error",
				"message": errors.Wrap(err, "could not select account").Error(),
			})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(input.Password))
		if (err != nil) && (err != bcrypt.ErrMismatchedHashAndPassword) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": errors.Wrapf(err, "could not compare hash and password").Error(),
			})
			return
		}

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid password",
			})
			return
		}

		expiryDate := time.Now().Add(time.Duration(input.ExpiryTime) * time.Second)

		// Create the token
		claims := jwt.MapClaims{
			"id":  acc.ID,
			"exp": expiryDate.Unix(),
			"iat": time.Now().Unix(),
		}
		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
		tokenString, err := token.SignedString([]byte(a.Config.JWTSecret))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": errors.Wrap(err, "could not insert token").Error(),
			})

			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status": "success",
			"data":   tokenString,
		})

		return
	case "implicit":
		fallthrough
	case "client_credentials":
		c.JSON(http.StatusNotImplemented, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("Flow %s is not implemented", input.GrantType),
		})
	}

	c.JSON(http.StatusUnprocessableEntity, gin.H{
		"message": "Validation failed",
		"error":   "Invalid action",
	})
	return
}
