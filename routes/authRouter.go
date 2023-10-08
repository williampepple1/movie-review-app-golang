package routes

import (
	controllers "review-app/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("users/signup", controllers.Signup())
	router.POST("users/login", controllers.Login())
}
