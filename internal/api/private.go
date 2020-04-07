package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (a *API) getHomepageResources(ctx *gin.Context) {
	type ExtendedResource struct {
		Resource
		AuthorUsername string `db:"author_username" json:"author_username"`
	}

	resources := []*ExtendedResource{}

	err := a.DB.SelectContext(ctx, &resources,
		"select r.*, u.username as author_username from resources as r, users as u where visibility = $1 and (r.author_id = u.id) order by updated_at desc limit 6",
		ResourceVisibilityPublic,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		a.Log.WithError(err).Errorln("could not select resources for listResources")
		return
	}

	// clean short descriptions
	for _, r := range resources {
		r.ShortDescription = strings.Split(r.Description, "\n")[0]
	}

	ctx.JSON(http.StatusOK, gin.H{
		"latest": resources,
	})
}
