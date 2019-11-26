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

func (a *API) checkResource(c *gin.Context) {
	resourceID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return
	}

	// Check if the resource exists
	var resource models.Resource
	if err := a.DB.Get(&resource, "select * from resources where id = $1", resourceID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "That resource could not be found",
			})
			c.Abort()
			return
		}

		a.Log.WithField("err", err).Errorln("Could not find resource")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": errors.Wrap(err, "Could not find resource"),
		})
		c.Abort()
		return
	}

	// Store the resource
	c.Set("resource", &resource)
}

func (a *API) createResource(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

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
		AuthorID:    user.ID,
	}

	result, err := a.DB.NamedExec("insert into resources (name, title, description, author_id) values (:name, :title, :description, :author_id)", &r)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": errors.Wrap(err, "could not insert").Error(),
		})

		return
	}

	fmt.Printf("%x\n", result)

	c.Status(http.StatusCreated)
}

func (a *API) getResource(c *gin.Context) {
	resource := c.MustGet("resource").(*models.Resource)
	c.JSON(http.StatusOK, resource)
}

func (a *API) voteResource(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	resource := c.MustGet("resource").(*models.Resource)

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
		Account:  user.ID,
		Resource: resource.ID,
		Positive: input.Positive,
	}

	result, err := a.DB.NamedExec(
		`insert into resource_votes
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
