package routes

import (
	"code_pilot/internal/controllers"
	"code_pilot/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.GET("/hello-world", controllers.HelloWorld)
	// user route
	api.POST("/user/register", controllers.SignUp)
	api.POST("/user/login", controllers.Login)
	api.GET("/user/:id", middlewares.IsAuthorized, controllers.GetUser)
	// project route
	api.POST("/project/create/:userID", middlewares.IsAuthorized, controllers.CreateProject)
}
