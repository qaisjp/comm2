package api

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/multitheftauto/community/internal/models"
	"github.com/multitheftauto/community/internal/resource"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// checkResourcePkg is a middleware that verifies that the package id
// exists for the current resource being accessed.
func (a *API) checkResourcePkg(c *gin.Context) {
	resource := c.MustGet("resource").(*Resource)

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
	var pkg ResourcePackage
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
	resource := c.MustGet("resource").(*Resource)

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
	pkg := c.MustGet("resource_pkg").(*ResourcePackage)
	user := c.MustGet("user").(*models.User)
	if pkg.Draft {
		ok := false
		if user != nil {
			var err error
			ok, err = a.canUserManageResource(c, user.ID, pkg.ResourceID)
			if err != nil {
				a.Log.WithError(err).Errorln("could not download package")
				c.Status(http.StatusInternalServerError)
				return
			}
		}

		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"message": "That package does not exist"})
			return
		}
	}
	c.JSON(http.StatusOK, pkg)
}

func (a *API) downloadResourcePackage(c *gin.Context) {
	resource := c.MustGet("resource_pkg").(*ResourcePackage)

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
	pkg := c.MustGet("resource_pkg").(*ResourcePackage)
	header, err := c.FormFile("file")
	if err == http.ErrNotMultipart {
		c.Status(http.StatusUnsupportedMediaType)
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"message": err.Error()})
		return
	} else if err == http.ErrMissingFile {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "missing file"})
		return
	} else if err != nil {
		a.Log.WithError(err).Errorln("could not upload package")
		c.Status(http.StatusInternalServerError)
		return
	}

	filename := "pkg" + strconv.FormatUint(pkg.ID, 10) + ".zip"

	// Make sure the file is a zip
	if typ := header.Header.Get("Content-Type"); typ != "application/zip" {
		c.Status(http.StatusUnsupportedMediaType)
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"message": "Expected zip, got " + typ})
		return
	}

	var zipBody []byte
	if f, err := header.Open(); err != nil {
		a.Log.WithError(err).Errorln("could not open file header's associated file when uploading package")
		c.Status(http.StatusInternalServerError)
		return
	} else if zipBody, err = ioutil.ReadAll(f); err != nil {
		a.Log.WithError(err).Errorln("could not read zip body when uploading package")
		c.Status(http.StatusInternalServerError)
		return
	} else if err := f.Close(); err != nil {
		a.Log.WithError(err).Errorln("could not close input file when uploading package")
		c.Status(http.StatusInternalServerError)
		return
	}

	f := bytes.NewReader(zipBody)

	if ok, reason, err := resource.CheckResourceZip(f, int64(len(zipBody))); err != nil {
		a.Log.WithError(err).Errorln("failed to check resource zip when uploading package")
		c.Status(http.StatusInternalServerError)
		return
	} else if !ok {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": reason,
		}) // todo(report): explain we picked this status code because https://httpstatuses.com/422
		return
	}

	w, err := a.Bucket.NewWriter(c, filename, nil)
	if err != nil {
		a.Log.WithError(err).Errorln("could not create new bucket writer when uploading package")
		c.Status(http.StatusInternalServerError)
		return
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		a.Log.WithError(err).Errorln("could not seek to start of file when uploading package")
		c.Status(http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(w, f); err != nil {
		a.Log.WithError(err).Errorln("could not copy from request to bucket when uploading package")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := w.Close(); err != nil {
		a.Log.WithError(err).Errorln("could not close writer when uploading package")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, pkg)
}
