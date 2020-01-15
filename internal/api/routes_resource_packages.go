package api

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/multitheftauto/community/internal/models"
	"github.com/pkg/errors"
)

// checkResourcePkg is a middleware that verifies that the package id
// exists for the current resource being accessed.
func (a *API) checkResourcePkg(c *gin.Context) {
	resource := c.MustGet("resource").(*models.Resource)

	// Parse pkg_id param
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

// createResourcePackage is an endpoint that creates a resource package draft
func (a *API) createResourcePackage(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	resource := c.MustGet("resource").(*models.Resource)

	var input struct {
		Description string `json:"description"`
	}

	if err := c.BindJSON(&input); err != nil {
		return
	}

	var id uint64
	err := a.QB.Insert("resource_packages").
		Columns("resource_id", "author_id", "description", "draft", "filename", "version").
		Values(resource.ID, user.ID, input.Description, true, "", "").Suffix("RETURNING id").
		ScanContext(c, &id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		a.Log.WithError(err).Errorln("database error creating resource package")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"package_id": id})
}

func (a *API) getResourcePackage(c *gin.Context) {
	resource := c.MustGet("resource_pkg").(*models.ResourcePackage)

	// c.Header("Content-Disposition",
	// c.Header("Cache-Control", "no-store")

	r, err := a.Bucket.NewReader(c, resource.Filename, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	c.DataFromReader(http.StatusOK, r.Size(), "application/zip", r, map[string]string{
		"Cache-Control":       "no-store",
		"Content-Disposition": fmt.Sprintf("attachment; filename=\"%s\"", resource.Filename),
	})
	// bytesWritten, copyErr := io.Copy(c.Writer, r)
	// if copyErr != nil {
	// 	fmt.Printf("Error copying file to the http response %s\n", copyErr.Error())
	// 	return
	// }
	// fmt.Printf("%d bytes writte\n", bytesWritten)

	// c.JSON(http.StatusOK, resource)
}

// uploadResourcePackage is an endpoint that uploads a file to an existing resource package
func (a *API) uploadResourcePackage(c *gin.Context) {
	pkg := c.MustGet("resource_pkg").(*models.ResourcePackage)
	header, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	fmt.Println(header.Filename)
	f, err := header.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	w, err := a.Bucket.NewWriter(c, header.Filename, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	bs, err := io.Copy(w, f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	fmt.Printf("%d bytes written\n", bs)

	if err := w.Close(); err != nil {
		fmt.Println("Could not close writer")
	}
	if err := f.Close(); err != nil {
		fmt.Println("Could not close file")
	}
	c.JSON(http.StatusOK, pkg)
}
