package api

import (
	"context"
	"net/http"

	"github.com/multitheftauto/community/pkg/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// API contains all the dependencies of the API server
type API struct {
	Config *config.Config
	Log    *logrus.Logger
	DB     *sqlx.DB
	Gin    *gin.Engine

	Server *http.Server
}

// NewAPI sets up a new API module.
func NewAPI(
	conf *config.Config,
	log *logrus.Logger,
	db *sqlx.DB,
) *API {

	router := gin.Default()

	a := &API{
		Config: conf,
		Log:    log,
		DB:     db,
		Gin:    router,
	}

	router.Use(cors.Default())

	router.POST("/v1/oauth", a.oauthToken)
	router.POST("/v1/accounts", a.createAccount)

	// resources := resources.Impl{API: a}
	// router.GET("/v1/resources", resources.List)
	router.POST("/v1/resources", a.authMiddleware, a.createResource)
	router.POST("/v1/resources/:id/like", a.authMiddleware, a.likeResource)

	return a
}

// Start binds the API and starts listening.
func (a *API) Start() error {
	a.Server = &http.Server{
		Addr:    a.Config.Address,
		Handler: a.Gin,
	}
	return a.Server.ListenAndServe()
}

// Shutdown shuts down the API
func (a *API) Shutdown(ctx context.Context) error {
	if err := a.Server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
