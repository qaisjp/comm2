package api

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/multitheftauto/community/internal/resource"
	"gocloud.dev/blob"

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
	if pkg.IsDraft() {
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
	res := c.MustGet("resource").(*Resource)

	fileUploadMode := false
	var meta *resource.XmlMeta
	var f io.Reader

	var input struct {
		Version     string `json:"version"`
		Description string `json:"description"`
		Draft       bool   `json:"bool"`
	}

	header, err := c.FormFile("file")
	if err == nil {
		fileUploadMode = true

		var message string
		var status int
		meta, f, status, message, err = a.readMetaFromheader(c, header)
		if status != http.StatusOK {
			if status == http.StatusInternalServerError {
				a.somethingWentWrong(c, err).Errorln("could not upload resource package")
				return
			}
			c.JSON(status, gin.H{"message": message})
			return
		}

		input.Version = meta.Infos[0].Version
		input.Draft = true

	} else if err != http.ErrNotMultipart {
		a.somethingWentWrong(c, err).Errorln("ctx.FormFile failed unexpectedly")
		return
	} else {
		// If there's no file, we are doing a file-less creation
		if err := c.BindJSON(&input); err != nil {
			return
		}
	}

	var publishedAt interface{} = pq.NullTime{}
	if !input.Draft {
		publishedAt = squirrel.Expr("now()")
	}

	query, args, err := a.QB.Insert("resource_packages").
		Columns("resource_id", "author_id", "description", "published_at", "version", "file_uploaded").
		Values(res.ID, user.ID, input.Description, publishedAt, input.Version, fileUploadMode).Suffix("RETURNING *").
		ToSql()
	if err != nil {
		a.somethingWentWrong(c, err).Errorln("database error creating resource package sql")
		return
	}

	tx, err := a.DB.BeginTxx(c, nil)
	if err != nil {
		a.somethingWentWrong(c, err).Errorln("could not create transaction")
		return
	}
	rollback := func() {
		if err := tx.Rollback(); err != nil {
			a.Log.WithError(err).Errorln("failed to rollback")
		}
	}

	row := tx.QueryRowxContext(c, query, args...)
	if err := row.Err(); err != nil {
		a.somethingWentWrong(c, err).Errorln("could not create resource package row")
		rollback()
		return
	}

	var pkg ResourcePackage
	if err := row.StructScan(&pkg); err != nil {
		a.somethingWentWrong(c, err).Errorln("could not create scan inserted resource package row")
		rollback()
		return
	}

	// Now that we have the ID, we upload the file
	if fileUploadMode {
		if err := pkg.writeToBucket(c, f, a.Bucket); err != nil {
			a.somethingWentWrong(c, err).Errorln("could not write pkg to bucket")
			if err := tx.Rollback(); err != nil {
				a.Log.WithError(err).Errorln("failed to rollback")
			}
			return
		}
	}

	if err := tx.Commit(); err != nil {
		a.somethingWentWrong(c, err).Errorln("could not commit during resource create pkg")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"package_id": pkg.ID})
}

func (a *API) getResourcePackage(c *gin.Context) {
	pkg := c.MustGet("resource_pkg").(*ResourcePackage)
	user := c.MustGet("current_user").(*User)
	if pkg.IsDraft() {
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

func (a *API) readMetaFromheader(ctx *gin.Context, header *multipart.FileHeader) (meta *resource.XmlMeta, f *bytes.Reader, status int, message string, err error) {
	status = http.StatusInternalServerError
	message = "Something went wrong"

	// Make sure the file is a zip
	if typ := header.Header.Get("Content-Type"); typ != "application/zip" {
		status = http.StatusUnsupportedMediaType
		message = "Expected zip, got " + typ
		return
	}

	var zipBody []byte
	var fzip multipart.File
	if fzip, err = header.Open(); err != nil {
		a.Log.WithError(err).Errorln("could not open file header's associated file when uploading package")
		return
	} else if zipBody, err = ioutil.ReadAll(fzip); err != nil {
		a.Log.WithError(err).Errorln("could not read zip body when uploading package")
		return
	} else if err = fzip.Close(); err != nil {
		a.Log.WithError(err).Errorln("could not close input file when uploading package")
		return
	}

	f = bytes.NewReader(zipBody)
	var ok bool
	if meta, ok, message, err = resource.CheckResourceZip(f, int64(len(zipBody))); err != nil {
		a.Log.WithError(err).Errorln("failed to check resource zip when uploading package")
		return
	} else if !ok {
		status = http.StatusUnprocessableEntity
		return
	}

	if _, err = f.Seek(0, io.SeekStart); err != nil {
		a.Log.WithError(err).Errorln("could not seek to start of file when uploading package")
		return
	}

	return meta, f, http.StatusOK, "", nil
}

func (pkg *ResourcePackage) writeToBucket(ctx context.Context, f io.Reader, bucket *blob.Bucket) error {
	w, err := bucket.NewWriter(ctx, pkg.GetBucketFilename(), nil)
	if err != nil {
		return errors.Wrap(err, "could not create new bucket writer when uploading package")
	}

	if _, err = io.Copy(w, f); err != nil {
		return errors.Wrap(err, "could not copy from request to bucket when uploading package")
	}

	if err = w.Close(); err != nil {
		return errors.Wrap(err, "could not close writer when uploading package")
	}
	return nil
}

// uploadResourcePackage is an endpoint that uploads a file to an existing resource package
func (a *API) uploadResourcePackage(ctx *gin.Context) {
	header, err := ctx.FormFile("file")
	if err == http.ErrNotMultipart {
		ctx.Status(http.StatusUnsupportedMediaType)
		ctx.JSON(http.StatusUnsupportedMediaType, gin.H{"message": err.Error()})
		return
	} else if err == http.ErrMissingFile {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": "missing file"})
		return
	} else if err != nil {
		a.Log.WithError(err).Errorln("could not upload package")
		ctx.Status(http.StatusInternalServerError)
		return
	}

	_, f, status, message, err := a.readMetaFromheader(ctx, header)
	if status != http.StatusOK {
		if status == http.StatusInternalServerError {
			a.somethingWentWrong(ctx, err).Errorln("could not upload resource package")
			return
		}
		ctx.JSON(status, gin.H{"message": message})
		return
	}

	pkg := ctx.MustGet("resource_pkg").(*ResourcePackage)
	if err := pkg.writeToBucket(ctx, f, a.Bucket); err != nil {
		a.somethingWentWrong(ctx, err).Errorln("could not write pkg to bucket")
		return
	}

	ctx.Status(http.StatusOK)
}
