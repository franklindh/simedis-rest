package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	if err != nil {
		log.Printf("Error: %v, Message: %s\n", err, message)
	}

	c.JSON(statusCode, gin.H{
		"status":  "error",
		"message": message,
		"data":    nil,
	})
}

func SuccessResponse(c *gin.Context, statusCode int, data any, message ...string) {
	msg := "success"
	if len(message) > 0 {
		msg = message[0]
	}

	c.JSON(statusCode, gin.H{
		"status":  "success",
		"message": msg,
		"data":    data,
	})
}
