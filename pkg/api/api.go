package api

import (
	"context"
	"net/http"
	"time"

	"github.com/multitheftauto/community/pkg/api/jwt"
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

	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "multitheftauto-api",
		Key:        []byte(conf.JWTSecret),
		Timeout:    time.Hour * 24,
		MaxRefresh: time.Hour * 24,

		Authenticator: a.jwtAuthenticate,
		Authorizator:  a.jwtAuthorize,
		Unauthorized:  a.jwtUnauthorized,

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	}

	router.POST("/v1/auth/login", authMiddleware.LoginHandler)
	router.POST("/v1/accounts", a.createAccount)

	// verifyAuth := authMiddleware.MiddlewareFunc()
	// resources := resources.Impl{API: a}
	// router.GET("/v1/resources", resources.List)
	// router.POST("/v1/resources", verifyAuth, resources.Patch)

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
