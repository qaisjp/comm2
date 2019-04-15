package api

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/multitheftauto/community/internal/models"
	jwt "gopkg.in/dgrijalva/jwt-go.v3"
)

func (a *API) authMiddleware(c *gin.Context) {
	// Get the "Authorization" header
	authorization := c.Request.Header.Get("Authorization")
	if authorization == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "Authorization header is malformed",
		})
		c.Abort()
		return
	}

	// Split it into two parts - "Bearer" and token
	parts := strings.SplitN(authorization, " ", 2)
	if parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, &gin.H{
			"status": "error",
			"error":  "Authorization header is malformed",
		})
		c.Abort()
		return
	}

	token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, errors.New("invalid signing algorithm")
		}
		return []byte(a.Config.JWTSecret), nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{
			"status": "error",
			"error":  errors.Wrap(err, "could not parse jwt").Error(),
		})
		c.Abort()
		return
	}

	if !token.Valid {
		c.JSON(http.StatusUnauthorized, &gin.H{
			"status": "error",
			"error":  "Your authentication token is invalid (has it expired?)",
		})
		c.Abort()
		return
	}

	id := token.Claims.(jwt.MapClaims)["id"].(float64)
	userID := uint64(id)

	account := &models.Account{}
	err = a.DB.Get(
		account,
		"select * from accounts where (id = $1)",
		userID,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusGone, gin.H{
			"status":  "error",
			"message": "Your account has been deleted",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errors.Wrapf(err, "could not check existence").Error(),
		})
		return
	}

	// Check for account activeness
	if !account.Activated {
		c.JSON(http.StatusUnauthorized, &gin.H{
			"status": "error",
			"error":  "Your account is not activated",
		})
		c.Abort()
		return
	}

	// Write account and token into environment
	c.Set("account", account)
	c.Set("token", token.Signature)

	c.Header(
		"X-Authenticated-As",
		string(account.ID),
	)
}
