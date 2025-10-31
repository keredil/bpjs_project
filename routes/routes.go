package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Endpoint CRUD Program BPJS
	router.GET("/programs", controllers.GetPrograms)
	router.POST("/programs", controllers.CreateProgram)
	router.PUT("/programs/:id", controllers.UpdateProgram)
	router.DELETE("/programs/:id", controllers.DeleteProgram)
}
