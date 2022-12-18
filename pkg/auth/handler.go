package auth

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type handler struct {
	DB *mongo.Database
}

func RegisterRoutes(r *gin.Engine, db *mongo.Database) {
	h := &handler{
		DB: db,
	}

	routes := r.Group("/api/v1/auth")
	routes.POST("/register", h.Register)
	routes.POST("/login", h.Login)
	routes.GET("/login/facebook", h.InitFacebookLogin)
	routes.GET("/facebook/callback", h.HandleFacebookLogin)
	routes.GET("/login/google", h.InitGoogleLogin)
	routes.GET("/google/callback", h.HandleGoogleLogin)
	routes.POST("/verify/email", h.VerifyEmail)
	routes.POST("/verify/resendCode", h.ResendCode)
}
