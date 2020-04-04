package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
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
			"message": "You don't have permission to access that resource.",
		})
		ctx.Abort()
		return
	}
}

func (a *API) checkResource(ctx *gin.Context) {
	var fieldVal interface{} = ctx.Param("resource_id")
	fieldName := "name"
	resourceID, err := strconv.ParseUint(ctx.Param("resource_id"), 10, 64)
	if err == nil {
		fieldVal = resourceID
		fieldName = "id"
	}

	user := ctx.MustGet("user").(*User)

	// Check if the resource exists
	var resource Resource
	if err := a.DB.Get(&resource, "select * from resources where "+pq.QuoteIdentifier(fieldName)+" = $1 and author_id=$2", fieldVal, user.ID); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "That resource could not be found",
			})
			ctx.Abort()
			return
		}

		a.Log.WithField("err", err).Errorln("Could not find resource")
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": errors.Wrap(err, "Could not find resource"),
		})
		ctx.Abort()
		return
	}

	// Populate resource.CanManage field
	currentUser := ctx.MustGet("current_user").(*User)
	if currentUser != nil {
		// Throw an error and abort if the author ID and user does not match
		ok, err := a.canUserManageResource(ctx, user.ID, resource.ID)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			ctx.Abort()
			return
		}
		resource.CanManage = ok
	}

	if resource.Visibility != ResourceVisibilityPublic && !resource.CanManage {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "You don't have permission to access that resource.",
		})
		ctx.Abort()
		return
	}

	// Store the resource
	ctx.Set("resource", &resource)
}

// listResources is an endpoint that allows you to list and search through resources.
//
// todo:
// - support authentication to include hidden stuff
// - exclude hidden stuff for unauthenticated requests
// - support search/filter fields
// - support pagination / cursors
func (a *API) listResources(ctx *gin.Context) {
	resources := []*Resource{}
	user := ctx.MustGet("current_user").(*User)

	suffix := "where visibility = $1"
	args := []interface{}{ResourceVisibilityPublic}
	if user != nil {
		args = append(args, user.ID)
		suffix = `
			, resource_collaborators as c
			where
			(r.visibility = $1) or
			(r.author_id = $2) or
			(
				(c.resource_id = r.id) and
				(c.user_id = $2)
				-- should probably check if accepted too, but it's fine
			)
		`
	}

	err := a.DB.SelectContext(ctx, &resources, "select r.* from resources as r "+suffix, args...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		a.Log.WithError(err).Errorln("could not select resources for listResources")
		return
	}

	ctx.JSON(http.StatusOK, resources)
}

func (a *API) deleteResource(ctx *gin.Context) {
	user := ctx.MustGet("current_user").(*User)
	resource := ctx.MustGet("resource").(*Resource)

	// Only the creator of a resource can delete it
	if user.ID != resource.AuthorID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Only the creator of a resource can delete it",
		})
		return
	}

	_, err := a.QB.Delete("resources").Where("id = $1", resource.ID).ExecContext(ctx)
	if err != nil {
		a.somethingWentWrong(ctx, err).Errorln("could not delete resource")
		return
	}

	ctx.Status(http.StatusNoContent)
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

	// Default title to name
	if r.Title == "" {
		r.Title = r.Name
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

func (a *API) getResource(ctx *gin.Context) {
	resource := ctx.MustGet("resource").(*Resource)

	type ResourceUserInfo struct {
		PublicUserInfo
		IsCreator bool `json:"is_creator"`
	}

	extended := struct {
		*Resource
		Authors []ResourceUserInfo `json:"authors"`
	}{
		Resource: resource,
	}

	// Add the creator to Authros
	{
		var creator User
		if err := a.DB.GetContext(ctx, &creator, "select * from users where id = $1", resource.AuthorID); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong."})
			a.Log.WithError(err).Errorln("getResource failed to get creator details")
			return
		}

		extended.Authors = append(extended.Authors, ResourceUserInfo{creator.PublicInfo(), true})
	}

	// Add the rest of the collaborators
	{
		authors := []User{}
		if err := a.DB.SelectContext(ctx, &authors, "select u.* from resource_collaborators as c, users as u where c.accepted and c.resource_id = $1 and c.user_id = u.id", resource.ID); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong."})
			a.Log.WithError(err).Errorln("getResource failed to get creator details")
			return
		}

		for _, a := range authors {
			extended.Authors = append(extended.Authors, ResourceUserInfo{a.PublicInfo(), false})
		}
	}

	ctx.JSON(http.StatusOK, extended)
}

func (a *API) transferResource(ctx *gin.Context) {
	resource := ctx.MustGet("resource").(*Resource)

	// todo: fix inconsistencies where sometimes we use user id and sometimes user name?
	var fields struct {
		NewOwner string `json:"new_owner,omitempty"`
	}
	if err := ctx.BindJSON(&fields); err != nil {
		return
	}

	if fields.NewOwner == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Name owner username cannot be empty"})
		return
	}

	var rows []struct {
		ID       string `db:"id"`
		Username string `db:"username"`
	}
	err := a.DB.SelectContext(ctx, &rows, "select id, username from users where lower(username) = $1", strings.ToLower(fields.NewOwner))
	if err != nil {
		a.somethingWentWrong(ctx, err).WithField("new-owner", fields.NewOwner).Errorln("could not get username")
		return
	}

	if len(rows) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "That user does not exist."})
		return
	} else if len(rows) > 1 {
		panic("not possible to have more than user row here")
	}

	_, err = a.QB.Update("resources").Set("author_id", rows[0].ID).Where(squirrel.Eq{"id": resource.ID}).ExecContext(ctx)
	if err != nil {
		a.somethingWentWrong(ctx, err).WithField("new-owner", fields.NewOwner).Errorln("could not set new owner")
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"new_username": rows[0].Username})
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

func (a *API) patchResource(ctx *gin.Context) {
	resource := ctx.MustGet("resource").(*Resource)
	var fields struct {
		Name        *string `json:"name,omitempty"`
		Title       *string `json:"title,omitempty"`
		Description *string `json:"description,omitempty"`
		Visibility  *string `json:"visibility,omitempty"`
		Archived    *bool   `json:"archived,omitempty"`
	}
	if err := ctx.BindJSON(&fields); err != nil {
		return
	}

	clauses := make(map[string]interface{})
	if fields.Title != nil {
		title := *fields.Title
		clauses["title"] = title
		if title == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Title cannot be empty"})
			return
		}
	}
	if fields.Description != nil {
		clauses["description"] = *fields.Description
	}
	if fields.Visibility != nil {
		vis := *fields.Visibility
		// todo: use ResourceVisibility enum and give an unmarshaller/marshaller
		if vis != ResourceVisibilityPrivate && vis != ResourceVisibilityPublic {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad 'visibility' field provided"})
			return
		}
		clauses["visibility"] = vis
	}
	if fields.Name != nil {
		name := *fields.Name
		clauses["name"] = name

		if name == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Name cannot be empty"})
			return
		}

		var count int
		err := a.DB.GetContext(ctx, &count, "select count(*) from resources where author_id=$1 and name=$2", resource.AuthorID, name)
		if err != nil {
			a.somethingWentWrong(ctx, err).WithField("id", resource.ID).WithField("new-name", name).
				Errorln("could not figure out if resource rename exists")
			return
		}

		if count > 0 {
			ctx.JSON(http.StatusConflict, gin.H{"message": "A resource with that name already exists"})
			return
		}
	}

	// Allow ops if (not_currently_archived || !*fields.Archived)
	shouldAllowChange := !resource.Archived
	if !shouldAllowChange && fields.Archived != nil {
		shouldAllowChange = !*fields.Archived
	}
	if !shouldAllowChange && len(clauses) != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "You can't perform that operation on an archived resource."})
		return
	}

	if fields.Archived != nil {
		clauses["archived"] = *fields.Archived
	}

	if len(clauses) == 0 {
		ctx.Status(http.StatusNotModified)
		return
	}

	result, err := a.QB.Update("resources").Where(squirrel.Eq{"id": resource.ID}).SetMap(clauses).ExecContext(ctx)
	if err != nil {
		a.somethingWentWrong(ctx, err).Errorln("could not update resource")
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		a.somethingWentWrong(ctx, err).Errorln("could not get number of rows affected")
	}

	// todo: rows is always 1. we prolly don't need this code.
	if rows == 0 {
		ctx.Status(http.StatusNotModified)
	} else {
		ctx.Status(http.StatusOK)
	}
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
		q = q.Where("published_at is not null")
	}
	q = q.OrderBy("published_at desc nulls first", "updated_at desc")

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
