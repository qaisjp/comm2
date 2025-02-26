package api

import (
	"context"
	"net/http"
	"time"

	sq "github.com/Masterminds/squirrel"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/multitheftauto/community/internal/config"
	"gocloud.dev/blob"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// API contains all the dependencies of the API server
type API struct {
	Config *config.Config
	Bucket *blob.Bucket
	Log    *logrus.Logger
	Gin    *gin.Engine
	DB     *sqlx.DB
	QB     sq.StatementBuilderType

	Server *http.Server
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

// NewAPI sets up a new API module.
func NewAPI(
	conf *config.Config,
	log *logrus.Logger,
	db *sqlx.DB,
	bucket *blob.Bucket,
) *API {

	// Create default gin router
	router := gin.Default()

	a := &API{
		Config: conf,
		Log:    log,
		Gin:    router,
		DB:     db,
		QB:     sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(db),
		Bucket: bucket,
	}

	// Handle CORS
	corsConf := cors.DefaultConfig()
	corsConf.AddAllowMethods("DELETE", "PATCH")
	corsConf.AddAllowHeaders("Authorization")
	corsConf.AllowAllOrigins = true
	router.Use(cors.New(corsConf))

	// Initialise JWT middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "multitheftauto.com",
		Key:             []byte(conf.JWTSecret),
		Timeout:         time.Hour * 24 * 7, // todo fix refresh in the frontend (this should be time.Hour * 6)
		MaxRefresh:      time.Hour * 24 * 3, // refresh is probably fine as it is though, since we're not using it yet
		IdentityKey:     "current_user",
		PayloadFunc:     a.jwtPayloadFunc,
		IdentityHandler: a.jwtIdentityHandler,
		Authenticator:   a.jwtAuthenticator,
		Authorizator:    a.jwtAuthorizator,
		Unauthorized:    a.jwtUnauthorized,

		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "query: token, cookie: jwt, header: Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",
	})

	if err != nil {
		log.WithField("error", err).Fatal("jwt error")
	}

	// Create JWT middleware
	authMiddlewareFunc := authMiddleware.MiddlewareFunc()
	authRequired := func(ctx *gin.Context) {
		user := ctx.MustGet("current_user").(*User)
		if user == nil {
			ctx.Header("WWW-Authenticate", "JWT realm="+authMiddleware.Realm)
			ctx.Abort()
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "You must be logged in to perform that operation.",
			})
		}
	}
	authMaybeRequired := func(ctx *gin.Context) {
		// Only execute auth if header present
		if _, ok := ctx.Request.Header["Authorization"]; ok {
			authMiddlewareFunc(ctx)
			return
		}

		var user *User
		ctx.Set("current_user", user)
	}

	private := router.Group("/private", authMaybeRequired)
	{
		private.GET("homepage", a.getHomepageResources)
		private.GET("/profile/:user_id", a.checkUser, a.getUserProfile)
		private.DELETE("/account", authRequired, a.deleteCurrentUser)
		private.POST("/account/username", authRequired, a.changeCurrentUserUsername)
		private.POST("/account/password", authRequired, a.changeCurrentUserPassword)
	}

	v1 := router.Group("/v1", authMaybeRequired)
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authMiddleware.LoginHandler)
			auth.POST("/refresh", authMiddleware.RefreshHandler)
			auth.POST("/register", a.createUser)
		}

		v1.GET("/resources", a.listResources)
		v1.POST("/resources", authRequired, a.createResource)
		// todo: make it possible to resolve resource_id on its own. Dummy `_` user id can be used, with a resolveUser middleware
		resources := v1.Group("/resources/:user_id/:resource_id", a.checkUser, a.checkResource)
		{
			resources.GET("", a.getResource)
			resources.PATCH("", authRequired, a.mustOwnResource, a.patchResource)
			resources.DELETE("", authRequired, a.mustOwnResource, a.deleteResource)

			collabMiddle := a.parseUserID("target_user", "target_user", true)
			if collab := resources.Group("collaborators/:target_user", authRequired, a.mustOwnResource, collabMiddle); true {
				collab.PUT("", a.addResourceCollaborator)
				collab.DELETE("", a.deleteResourceCollaborator)
			}

			resources.POST("/transfer", authRequired, a.mustOwnResource, a.transferResource)

			resources.POST("/vote", authRequired, a.voteResource)

			resources.GET("/pkg", a.listResourcePackages)
			resources.POST("/pkg", authRequired, a.mustOwnResource, a.createResourcePackage)
			pkg := resources.Group("/pkg/:pkg_id", a.checkResourcePkg)
			{
				pkg.GET("", a.getResourcePackage)

				pkg.GET("/download", a.downloadResourcePackage)
				pkg.PUT("/upload", authRequired, a.mustOwnResource, a.uploadResourcePackage)
			}
		}

		users := v1.Group("/users/:user_id", a.checkUser)
		{
			users.GET("", a.getUser)
			users.GET("/followers", a.getUserFollowers)
			users.GET("/following", a.getUserFollowing)
		}

		user := v1.Group("/user", authRequired)
		{
			user.GET("", a.getCurrentUser)

			user.GET("/profile", a.getCurrentUserProfile)
			user.PATCH("/profile", a.patchCurrentUserProfile)

			follow := user.Group("/follow/:target_user", a.parseUserID("target_user", "target_user", true))
			{
				follow.GET("", a.followUser)
				follow.PUT("", a.followUser)
				follow.DELETE("", a.followUser)
			}
		}
	}

	return a
}

func (a *API) somethingWentWrong(ctx *gin.Context, err error) *logrus.Entry {
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
	return a.Log.WithError(err)
}
