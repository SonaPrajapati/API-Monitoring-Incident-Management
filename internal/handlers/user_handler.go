package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Profile(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "You are authorized",
	})
}
