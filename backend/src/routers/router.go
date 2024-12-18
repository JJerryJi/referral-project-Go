package routers

import (
	"github.com/JerryJi/referral-finder/src/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func Register(db *gorm.DB, r *gin.Engine) {
	alumniHandler := handlers.NewAlumniHandler(db)

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

}
