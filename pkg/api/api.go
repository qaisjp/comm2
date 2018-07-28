package api

import (
	"time"

	"github.com/multitheftauto/community/pkg/api/auth"
	"github.com/multitheftauto/community/pkg/api/base"
	"github.com/multitheftauto/community/pkg/api/jwt"
	"github.com/multitheftauto/community/pkg/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// NewAPI sets up a new API module.
func NewAPI(
	conf *config.Config,
	log *logrus.Logger,
	db *sqlx.DB,
) *base.API {

	router := gin.Default()

	a := &base.API{
		Config: conf,
		Log:    log,
		DB:     db,
		Gin:    router,
	}

	router.Use(cors.Default())

	auth := auth.Impl{API: a}

	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "multitheftauto-api",
		Key:        []byte(conf.JWTSecret),
		Timeout:    time.Hour * 24,
		MaxRefresh: time.Hour * 24,

		Authenticator: auth.Authenticate,
		Authorizator:  auth.Authorize,
		Unauthorized:  auth.Unauthorized,

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	}

	router.POST("/v1/auth/login", authMiddleware.LoginHandler)
	router.POST("/v1/auth/register", auth.Register)

	// verifyAuth := authMiddleware.MiddlewareFunc()
	// resources := resources.Impl{API: a}
	// router.GET("/v1/resources", resources.List)
	// router.POST("/v1/resources", verifyAuth, resources.Patch)

	return a
}
