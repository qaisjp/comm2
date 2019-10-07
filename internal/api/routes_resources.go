package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/multitheftauto/community/internal/models"
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
			"message": errors.Wrap(err, "could not insert").Error(),
		})

		return
	}

	fmt.Printf("%x\n", result)

	c.Status(http.StatusCreated)
}

func (a *API) likeResource(c *gin.Context) {
	account := c.MustGet("account").(*models.Account)

	resource, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	// Check if the resource exists
	var count int
	if err := a.DB.Get(&count, "select count(id) from resources where id = $1", resource); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "That resource could not be found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": errors.Wrap(err, "Could not find resource"),
		})
		return
	}

	var input struct {
		Positive bool `json:"positive"`
	}

	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	r := models.ResourceRating{
		Account:  account.ID,
		Resource: resource,
		Positive: input.Positive,
	}

	result, err := a.DB.NamedExec(
		`insert into resource_ratings
		(resource, account, positive)
		values (:resource, :account, :positive)
		on conflict (resource, account)
		do update set positive = :positive
		`,
		&r,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": errors.Wrap(err, "could not update").Error(),
		})

		return
	}

	fmt.Printf("%+v\n", result)

	c.JSON(http.StatusNoContent, gin.H{
		"status": "success",
	})
}
