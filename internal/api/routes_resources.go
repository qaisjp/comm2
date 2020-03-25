package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var reResourceName = regexp.MustCompile("[a-zA-Z0-9]+")

func isResourceNameValid(str string) bool {
	return reResourceName.MatchString(str)
}

// mustOwnResource is a middleware that ensures that the
// authenticated user owns the resource being accessed.
func (a *API) mustOwnResource(ctx *gin.Context) {
	// Get our user and resource
	user := ctx.MustGet("current_user").(*User)
	if user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Authentication required."})
		ctx.Abort()
		return
	}

	resource := ctx.MustGet("resource").(*Resource)

	// Throw an error and abort if the author ID and user does not match
	if ok, err := a.canUserManageResource(ctx, user.ID, resource.ID); err != nil {
		ctx.Status(http.StatusInternalServerError)
		ctx.Abort()
		return
	} else if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "You don't own that resource",
		})
		ctx.Abort()
		return
	}
}

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
	var resource Resource
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

// listResources is an endpoint that allows you to list and search through resources.
//
// todo:
// - support authentication to include hidden stuff
// - exclude hidden stuff for unauthenticated requests
// - support search/filter fields
// - support pagination / cursors
func (a *API) listResources(c *gin.Context) {
	resources := []*Resource{}
	err := a.DB.SelectContext(c, &resources, "select * from resources;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		a.Log.WithError(err).Errorln("could not select resources for listResources")
		return
	}

	c.JSON(http.StatusOK, resources)
}

func (a *API) deleteResource(c *gin.Context) {
	user := c.MustGet("current_user").(*User)
	resource := c.MustGet("resource").(*Resource)

	// Only the creator of a resource can delete it
	if user.ID != resource.AuthorID {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Only the creator of a resource can delete it",
		})
	}

	_, err := a.QB.Delete("resources").Where("id = $1", resource.ID).ExecContext(c)
	if err != nil {
		a.Log.WithError(err).Errorln("could not delete resource")
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func (a *API) createResource(c *gin.Context) {
	user := c.MustGet("current_user").(*User)

	var input struct {
		Name        string `json:"name"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := c.BindJSON(&input); err != nil {
		return
	}

	// Expect at least a name to be set
	if input.Name == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "You must provide a name",
		})
		return
	} else if !isResourceNameValid(input.Name) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Name contains invalid characters",
		})
		return
	}

	r := Resource{
		Name:        input.Name,
		Title:       input.Title,
		Description: input.Description,
		AuthorID:    user.ID,
	}

	var count int
	if err := a.DB.GetContext(c, &count, "select count(*) from resources where name=$1", r.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong.",
		})
		a.Log.WithField("err", err).Errorln("could not check for resource existence")
		return
	}

	if count != 0 {
		c.JSON(http.StatusConflict, gin.H{
			"message": "That name is taken.",
		})
		return
	}

	var id int64
	err := a.DB.QueryRowxContext(c,
		"insert into resources (name, title, description, author_id) values ($1, $2, $3, $4) returning id",
		r.Name, r.Title, r.Description, r.AuthorID).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": errors.Wrap(err, "could not insert").Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (a *API) getResource(c *gin.Context) {
	resource := c.MustGet("resource").(*Resource)
	c.JSON(http.StatusOK, resource)
}

func (a *API) voteResource(c *gin.Context) {
	user := c.MustGet("current_user").(*User)
	resource := c.MustGet("resource").(*Resource)

	var input struct {
		Positive bool `json:"positive"`
	}

	if err := c.BindJSON(&input); err != nil {
		a.Log.WithError(err).Errorln("could not BindJSON")
		return
	}

	r := ResourceRating{
		UserID:     user.ID,
		ResourceID: resource.ID,
		Positive:   input.Positive,
	}

	result, err := a.DB.NamedExec(
		`insert into resource_votes
		(resource_id, user_id, positive)
		values (:resource_id, :user_id, :positive)
		on conflict (resource_id, user_id)
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

	c.Status(http.StatusNoContent)
}

func (a *API) listResourcePackages(ctx *gin.Context) {
	resource := ctx.MustGet("resource").(*Resource)
	user := ctx.MustGet("current_user").(*User)

	if user != nil {
		manages, err := a.canUserManageResource(ctx, user.ID, resource.ID)
		if err != nil {
			a.Log.WithError(err).Errorln("failed to check if user manages resource")
			ctx.Status(http.StatusInternalServerError)
			return
		}

		// Unset user if they are not a manager
		if !manages {
			user = nil
		}
	}

	q := a.QB.Select("*").From("resource_packages").Where("resource_id = $1", resource.ID)
	if user == nil {
		q = q.Where("not draft")
	}

	query, values, err := q.ToSql()
	if err != nil {
		a.Log.WithError(err).Errorln("could not convert to sql")
		ctx.Status(http.StatusInternalServerError)
		return
	}

	rows := []*ResourcePackage{}
	if err := a.DB.SelectContext(ctx, &rows, query, values...); err != nil {
		a.Log.WithError(err).Errorln("could not retrieve resource packages")
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, rows)
}
