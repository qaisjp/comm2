package auth

import (
	"net/http"

	"github.com/multitheftauto/community/pkg/models"

	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (i *Impl) Authenticate(username string, password string, c *gin.Context) (userID uint64, success bool) {
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
