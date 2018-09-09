package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/multitheftauto/community/pkg/models"
	"github.com/pkg/errors"
)

func (a *API) createResource(c *gin.Context) {
	account := c.MustGet("account").(*models.Account)

	var input struct {
		Name        string `json:"name"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	r := models.Resource{
		Name:        input.Name,
		Title:       input.Title,
		Description: input.Description,
		Creator:     account.ID,
	}

	result, err := a.DB.NamedExec("insert into resources (name, title, description, creator) values (:name, :title, :description, :creator)", &r)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errors.Wrap(err, "could not insert").Error(),
		})

		return
	}

	fmt.Printf("%x\n", result)

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
	})
}
