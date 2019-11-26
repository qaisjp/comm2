package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"github.com/multitheftauto/community/internal/models"
	"github.com/pkg/errors"
)

func (a *API) mustOwnResource(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	resource := c.MustGet("resource").(*models.Resource)

	if resource.AuthorID != user.ID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "You don't own that resource",
		})
		c.Abort()
		return
	}
}

func (a *API) checkResourcePkg(c *gin.Context) {
	resource := c.MustGet("resource").(*models.Resource)

	pkgID, err := strconv.ParseUint(c.Param("pkg_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return
	}

	// Check if the resource package exists
	var pkg models.ResourcePackage
	if err := a.DB.Get(&pkg, "select * from resource_packages where id = $1 and resource_id = $2", pkgID, resource.ID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "That resource package could not be found",
			})
			c.Abort()
			return
		}

		a.Log.WithField("err", err).Errorln("Could not find resource package")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": errors.Wrap(err, "Could not find resource package"),
		})
		c.Abort()
		return
	}

	// Store the resource package
	c.Set("resource_pkg", &pkg)
}

func (a *API) createResourcePackage(c *gin.Context) {
	fmt.Println("Upload resource package")
	user := c.MustGet("user").(*models.User)
	resource := c.MustGet("resource").(*models.Resource)

	var input struct {
		Description string `json:"description"`
	}

	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	s, args, err := qb.Insert("resource_packages").Columns("resource_id", "author_id", "description", "draft", "filename", "version").Values(resource.ID, user.ID, input.Description, true, "", "").ToSql()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Yikes"})
		return
	}
	fmt.Printf("sql: %#v\nargs:%#v\nerr: %#v", s, args, err)

	result, err := a.DB.Exec(s, args...)

	// result, err := a.DB.NamedExec("insert into resource_packages (resource_id, author_id, description, author_id) values (:name, :title, :description, :author_id)", &r)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": errors.Wrap(err, "could not insert").Error(),
		})
		return
	}

	fmt.Printf("Result: %#v\n", result)

	c.Status(http.StatusCreated)
}

func (a *API) getResourcePackage(c *gin.Context) {
	resource := c.MustGet("resource_pkg").(*models.ResourcePackage)
	c.JSON(http.StatusOK, resource)
}
