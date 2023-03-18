package routes

import (
	controller "github.com/sahil/jwt-auth-go/controllers"

	"github.com/gin-gonic/gin"
)

// UserRoutes function
func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controller.NewController().SignUp)
	incomingRoutes.POST("/users/login", controller.NewController().Login)
}
