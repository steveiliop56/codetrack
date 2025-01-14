package utils

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func HandleError(err error, exit bool) {
	if err != nil {
		if exit {
			log.Fatal().Str("error", err.Error()).Msg("A fatal error occurred, exiting")
			os.Exit(1)
		} else {
			log.Error().Str("error", err.Error()).Msg("An error occurred")
		}
	}
}

func APIBadRequest(err error, c *gin.Context) (bool){
	if err != nil {
		c.JSON(400, gin.H{
			"status": 400,
			"message": "Bad request",
		})
		return true
	}
	return false
}

func APIInternalError(err error, c *gin.Context) (bool) {
	if err != nil {
		c.JSON(500, gin.H{
			"status": 500,
			"message": "Internal server error",
		})
		return true
	}
	return false
}