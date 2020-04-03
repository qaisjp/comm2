package api

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/multitheftauto/community/internal/resource"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// checkResourcePkg is a middleware that verifies that the package id
// exists for the current resource being accessed.
func (a *API) checkResourcePkg(ctx *gin.Context) {
	resource := ctx.MustGet("resource").(*Resource)

	// Parse pkg_id param
	pkgID, err := strconv.ParseUint(ctx.Param("pkg_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		ctx.Abort()
		return
	}

	// Check if the resource package exists
	var pkg ResourcePackage
	if err := a.DB.Get(&pkg, "select * from resource_packages where id = $1 and resource_id = $2", pkgID, resource.ID); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "That resource package could not be found",
			})
			ctx.Abort()
			return
		}

		a.Log.WithField("err", err).Errorln("Could not find resource package")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": errors.Wrap(err, "Could not find resource package"),
		})
		return
	}

	// If draft, run mustOwnResource middleware
	if pkg.Draft {
		a.mustOwnResource(ctx)
		if ctx.IsAborted() {
			return
		}
	}

	bucketFilename := pkg.GetBucketFilename()
	pkg.FileUploaded, err = a.Bucket.Exists(ctx, bucketFilename)
	if err != nil {
		a.Log.WithError(err).WithField("filename", bucketFilename).WithField("pkg", pkg.ID).Errorln("could not check bucket for file existence")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong.",
		})
		return
	}

	// Store the resource package
	ctx.Set("resource_pkg", &pkg)
}

// createResourcePackage is an endpoint that creates a resource package draft
func (a *API) createResourcePackage(c *gin.Context) {
	user := c.MustGet("current_user").(*User)
	resource := c.MustGet("resource").(*Resource)

	var input struct {
		Version     string `json:"version"`
		Description string `json:"description"`
		Draft       bool   `json:"bool"`
	}

	if err := c.BindJSON(&input); err != nil {
		return
	}

	var id uint64
	err := a.QB.Insert("resource_packages").
		Columns("resource_id", "author_id", "description", "draft", "version").
		Values(resource.ID, user.ID, input.Description, input.Draft, input.Version).Suffix("RETURNING id").
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
	user := c.MustGet("current_user").(*User)
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

func (a *API) downloadResourcePackage(ctx *gin.Context) {
	resource := ctx.MustGet("resource").(*Resource)
	pkg := ctx.MustGet("resource_pkg").(*ResourcePackage)

	// c.Header("Content-Disposition",
	// c.Header("Cache-Control", "no-store")

	if !pkg.FileUploaded {
		ctx.Status(http.StatusNotFound)
		return
	}

	bucketFilename := pkg.GetBucketFilename()
	r, err := a.Bucket.NewReader(ctx, bucketFilename, nil)
	if err != nil {
		a.Log.WithError(err).WithField("filename", bucketFilename).WithField("pkg", pkg.ID).Errorln("could not create reader from bucket")
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.DataFromReader(http.StatusOK, r.Size(), "application/zip", r, map[string]string{
		"Cache-Control":       "no-store",
		"Content-Disposition": fmt.Sprintf("attachment; filename=\"%s\"", resource.Name+".zip"),
	})
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

	w, err := a.Bucket.NewWriter(c, pkg.GetBucketFilename(), nil)
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
	c.Status(http.StatusOK)
}
