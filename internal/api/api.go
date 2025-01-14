package api

import (
	"codetrack/internal/db"
	"codetrack/internal/services"
	"codetrack/internal/types"
	"codetrack/internal/utils"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	gormsessions "github.com/gin-contrib/sessions/gorm"
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
	Store cookie.Store
	Services *services.Services
}

func (api *Api) Initialize(secret string) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	store := gormsessions.NewStore(api.Database.Database, true, []byte(secret))
	router.Use(zerolog())
	router.Use(sessions.Sessions("codetrack", store))
	api.Router = router
	api.Store = store
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
	api.Router.POST("/api/accounts/register", func (c *gin.Context) {
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

	api.Router.POST("/api/accounts/login", func (c *gin.Context) {
		session := sessions.Default(c)
		user := types.User{}

		bindErr := c.BindJSON(&user)
		if utils.APIBadRequest(bindErr, c) {
			return
		}

		userExists, userExistsErr := api.Services.UserExists(user.Email)
		if utils.APIInternalError(userExistsErr, c) {
			return
		}

		if !userExists {
			c.JSON(400, gin.H{
				"status": 400,
				"message": "User does not exist",
			})
			return
		}

		email := session.Get("email")

		if email == user.Email {
			c.JSON(400, gin.H{
				"status": 400,
				"message": "Already logged in",
			})
			return
		}

		login, loginErr := api.Services.Login(user.Email, user.Password)

		if utils.APIInternalError(loginErr, c) {
			return
		}

		if !login {
			c.JSON(400, gin.H{
				"status": 400,
				"message": "Invalid credentials",
			})
			return
		}

		session.Set("email", user.Email)
		session.Save()

		c.JSON(200, gin.H{
			"status": 200,
			"message": "Logged in",
		})
	})

	api.Router.POST("/api/accounts/logout", api.sessionCheck(), func (c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()

		c.JSON(200, gin.H{
			"status": 200,
			"message": "Logged out",
		})
	})

	api.Router.DELETE("/api/accounts/delete", api.sessionCheck(), func (c *gin.Context) {
		session := sessions.Default(c)
		email := session.Get("email")

		emailString, emailOk := email.(string)

		if !emailOk {
			c.JSON(500, gin.H{
				"status": 500,
				"message": "Internal server error",
			})
			return
		}

		deleteErr := api.Services.DeleteUser(emailString)
		if utils.APIInternalError(deleteErr, c) {
			return
		}

		session.Clear()
		session.Save()

		c.JSON(200, gin.H{
			"status": 200,
			"message": "User deleted",
		})
	})

	api.Router.GET("/api/accounts/me", api.sessionCheck(), func (c *gin.Context) {
		session := sessions.Default(c)
		email := session.Get("email")

		emailString, emailOk := email.(string)

		if !emailOk {
			c.JSON(500, gin.H{
				"status": 500,
				"message": "Internal server error",
			})
			return
		}

		c.JSON(200, gin.H{
			"status": 200,
			"message": "OK",
			"data": gin.H{
				"email": emailString,
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

func (api *Api) sessionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		email := session.Get("email")

		emailString, emailOk := email.(string)

		if !emailOk {
			c.JSON(401, gin.H{
				"status": 401,
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		login, loginErr := api.Services.EmailLogin(emailString)
		if utils.APIInternalError(loginErr, c) {
			c.Abort()
			return
		}

		if !login {
			c.JSON(401, gin.H{
				"status": 401,
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}