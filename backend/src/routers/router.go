package routers

import (
	"github.com/JerryJi/referral-finder/src/handlers"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
)

func Register(db *gorm.DB, r *gin.Engine, rdb *redis.Client) {
	alumniHandler := handlers.NewAlumniHandler(db)
	tokenHandler := handlers.NewTokenHandler(rdb, db)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	user := r.Group("/user")
	{
		user.POST("/alumni", alumniHandler.CreateNewAlumni)
		user.GET("/alumni", alumniHandler.AdminViewAlumniInfo)
		user.GET("/alumni/:id", alumniHandler.ViewAlumniInfo)
		user.PUT("/alumni/:id", alumniHandler.UpdateAlumniInfo)
	}

	//Token-related routes
	token := r.Group("/token")
	{
		token.GET("", tokenHandler.Retrieve)  // Retrieve user info using token
		token.POST("", tokenHandler.Generate) // Authenticate and generate token
	}
}
