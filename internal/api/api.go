package api

import (
	"codetrack/internal/db"
	"codetrack/internal/services"
	"codetrack/internal/types"
	"codetrack/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/samber/do/v2"
)

func NewApi(i do.Injector) (*Api, error) {
	return &Api{
		Database: do.MustInvoke[*db.Database](i),
		Services: do.MustInvoke[*services.Services](i),
	}, nil
}

type Api struct {
	Database *db.Database
	Router *gin.Engine
	Services *services.Services
}

func (api *Api) Initialize() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(zerolog())
	api.Router = router
}

func (api *Api) SetupRoutes() {
	// Healthcheck
	api.Router.GET("/api/healthcheck", func (c *gin.Context) {
		c.JSON(200, gin.H{
			"status": 200,
			"message": "OK",
		})
	})

	// Accounts
	api.Router.POST("/api/accounts/new", func (c *gin.Context) {
		user := types.User{}

		bindErr := c.BindJSON(&user)
		if utils.APIBadRequest(bindErr, c) {
			return
		}

		userExists, userExistsErr := api.Services.UserExists(user.Email)
		if utils.APIInternalError(userExistsErr, c) {
			return
		}

		if userExists {
			c.JSON(400, gin.H{
				"status": 400,
				"message": "User already exists",
			})
			return
		}

		registerErr := api.Services.Register(user.Email, user.Password)
		if utils.APIInternalError(registerErr, c) {
			return
		}

		c.JSON(200, gin.H{
			"status": 200,
			"message": "User registered",
			"data": gin.H{
				"email": user.Email,
				"password": user.Password,
			},
		})
	})
}

func (api *Api) Start() {
	log.Info().Int("port", 8080).Msg("Listening for requests")
	api.Router.Run(":8080")
}

func zerolog() gin.HandlerFunc {
	return func(c *gin.Context) {
		tStart := time.Now()

		c.Next()

		code := c.Writer.Status()
		address := c.Request.RemoteAddr
		method := c.Request.Method
		path := c.Request.URL.Path
		
		latency := time.Since(tStart).String()

		switch {
			case code >= 200 && code < 300:
				log.Info().Str("method", method).Str("path", path).Str("address", address).Int("status", code).Str("latency", latency).Msg("Request")
			case code >= 300 && code < 400:
				log.Warn().Str("method", method).Str("path", path).Str("address", address).Int("status", code).Str("latency", latency).Msg("Request")
			case code >= 400:
				log.Error().Str("method", method).Str("path", path).Str("address", address).Int("status", code).Str("latency", latency).Msg("Request")
		}
	}
}