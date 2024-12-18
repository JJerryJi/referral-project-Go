package handlers

import "github.com/gin-gonic/gin"

type AlumniHandlerInterface interface {
	CreateNewAlumni(c *gin.Context)
	ViewAlumniInfo(c *gin.Context)
	UpdateAlumniInfo(c *gin.Context)
	AdminViewAlumniInfo(c *gin.Context)
}
