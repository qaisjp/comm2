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
	// Rating      int    `db:"rating"`
	// Downloads   int    `db:"downloads"`
	// Type        int    `db:"type"` // todo: ResourceType
	Status string `db:"status" json:"status"`
}

const (
	ResourceStatusPublic  string = "public"
	ResourceStatusPrivate string = "private"
)

type ResourcePackage struct {
	ID        uint64    `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	ResourceID  uint64 `db:"resource_id" json:"resource_id"` // relation
	AuthorID    uint64 `db:"author_id" json:"author_id"`     // relation
	Version     string `db:"version" json:"version"`
	Description string `db:"description" json:"description"`

	Filename string `db:"filename" json:"filename"`

	Draft bool `db:"draft" json:"draft"`
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
