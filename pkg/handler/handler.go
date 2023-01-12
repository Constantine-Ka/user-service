package handler

import (
	"github.com/Constantine-Ka/user-service/pkg/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Authorization", "Content-Type"}
	corsConfig.ExposeHeaders = []string{"Origin", "Authorization", "Content-Type"}
	router.Use(cors.New(corsConfig))
	auth := router.Group("/auth")
	{
		auth.POST("/sing-up", h.singUp)
		auth.POST("/sing-in", h.singIn)
		auth.GET("/confirm", h.confirm)
		auth.GET("/resetpassword", h.resetLink)
		auth.POST("/resetpassword", h.resetPassword)
		auth.POST("/update-jwt", h.updateToken)

	}
	users := router.Group("/user", h.userIndentity)
	{
		users.GET("/me", h.meUser)
		users.GET("/all", h.getUserAll)
		users.PUT("/links", h.updateLink)
		users.PUT("/update", h.updateUser)
	}
	router.StaticFile("/logs", "./logs/all.json")
	return router
}
