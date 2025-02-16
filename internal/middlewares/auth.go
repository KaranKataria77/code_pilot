package middlewares

import (
	"code_pilot/internal/utils"
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
)

func IsAuthorized(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	_, err := utils.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		c.Abort()
		return
	}
	log.Println("User validated")
	c.Next()
}
