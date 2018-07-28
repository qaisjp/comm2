package auth

import (
	"net/http"

	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/multitheftauto/community/pkg/models"
	"github.com/pkg/errors"
)

func (i *Impl) Authorize(userId uint64, c *gin.Context) bool {
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

	return (u.Level > 1) && (!u.Banned) && (u.Activated)
}
