package api

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Resource struct {
	ID        uint64    `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	AuthorID  uint64    `db:"author_id" json:"author_id"`

	Name        string `db:"name" json:"name"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	ShortDescription string `db:"-" json:"short_description"`
	// Rating      int    `db:"rating"`
	// Downloads   int    `db:"downloads"`
	// Type        int    `db:"type"` // todo: ResourceType
	Visibility string `db:"visibility" json:"visibility"`
	Archived   bool   `db:"archived" json:"archived"`
	DownloadCount int `db:"download_count" json:"download_count"`

	CanManage bool `db:"-" json:"can_manage"`
}

const (
	ResourceVisibilityPublic  string = "public"
	ResourceVisibilityPrivate string = "private"
)

type ResourcePackage struct {
	ID        uint64    `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	ResourceID  uint64 `db:"resource_id" json:"resource_id"` // relation
	AuthorID    uint64 `db:"author_id" json:"author_id"`     // relation
	Version     string `db:"version" json:"version"`
	Description string `db:"description" json:"description"`

	PublishedAt  *time.Time `db:"published_at" json:"published_at"`
	FileUploaded bool       `db:"-" json:"file_uploaded"`
	UploadedAt   *time.Time `db:"uploaded_at" json:"uploaded_at"`
}

func (pkg *ResourcePackage) GetBucketFilename() string {
	return fmt.Sprintf("res%d/pkg%d.zip", pkg.ResourceID, pkg.ID)
}

func (pkg *ResourcePackage) IsDraft() bool {
	return pkg.PublishedAt == nil
}

type ResourceRating struct {
	UserID     uint64 `db:"user_id" json:"user_id"`
	ResourceID uint64 `db:"resource_id" json:"resource_id"`
	Positive   bool   `db:"positive" json:"positive"`
}

// canUserManageResource checks if a given user can manage a given resource
func (a *API) canUserManageResource(ctx *gin.Context, userID uint64, resourceID uint64) (canAccess bool, err error) {
	// Check the resource from context if the resource ID matches
	if data, ok := ctx.Get("resource"); ok {
		resource := data.(*Resource)
		if resource.ID != resourceID {
			// If they are the owner, return true
			if resource.AuthorID == userID {
				fmt.Println("quik owner")
				return true, nil
			}

			// Otherwise just check the resource_collaborators table
			err = a.QB.Select("true").From("resource_collaborators").
				Where("accepted AND resource_id = $1 AND user_id = $2", resourceID, userID).
				ScanContext(ctx, &canAccess)

			if err == sql.ErrNoRows {
				err = nil
			} else if err != nil {
				err = errors.Wrap(err, "sql query failed")
			}

			return
		}
	}

	err = a.DB.GetContext(ctx, &canAccess, `
		select true from resource_collaborators where accepted and resource_id = $1 and user_id = $2
			union distinct
		select true from resources where id = $1 and author_id = $2
	`, resourceID, userID)

	// If the error is a lack of rows, suppress it
	if err == sql.ErrNoRows {
		err = nil
	}

	return
}
