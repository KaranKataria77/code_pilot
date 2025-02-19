package routes

import (
	"code_pilot/internal/controllers"
	grpcclient "code_pilot/internal/grpc_client"

	// grpcclient "code_pilot/internal/grpc_client"
	"code_pilot/internal/middlewares"
	"code_pilot/internal/supabase"
	"code_pilot/internal/websocket"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.GET("/hello-world", controllers.HelloWorld)
	// user route
	api.GET("/user", middlewares.IsAuthorized, controllers.GetUser)
	api.POST("/user/register", controllers.SignUp)
	api.POST("/user/login", controllers.Login)
	api.POST("/user/execute-code", grpcclient.CodeExecution)
	api.GET("/ws", websocket.HandleWebsocket)
	api.GET("/fileupload", supabase.GetFolderStructure)
	// project route
	api.POST("/project/create/:userID", middlewares.IsAuthorized, controllers.CreateProject)
}
