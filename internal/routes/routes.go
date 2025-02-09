package routes

import (
	"code_pilot/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.GET("/hello-world", controllers.HelloWorld)
	// user route
	api.POST("/register", controllers.SignUp)
	api.POST("/login", controllers.Login)
}
