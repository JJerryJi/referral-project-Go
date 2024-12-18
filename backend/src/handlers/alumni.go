package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/JerryJi/referral-finder/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AlumniHandler struct {
	DB *gorm.DB
}

func NewAlumniHandler(db *gorm.DB) *AlumniHandler {
	return &AlumniHandler{DB: db}
}

func (h *AlumniHandler) CreateNewAlumni(c *gin.Context) {
	// Struct to hold incoming request data
	var input AlumniCreate
	// Parse and validate the request body
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	// Start the transaction
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		// Step 1: Create the User
		user := models.User{
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Email:     input.Email,
			Username:  input.Username,
			Password:  input.Password, // Implement password hashing
			Location:  input.Location,
			Role:      input.Role,
		}

		if err := tx.Create(&user).Error; err != nil {
			return err // Rollback the transaction
		}

		// Step 2: Create the Alumni
		alumni := models.Alumni{
			UserID:      &user.ID,
			CompanyName: input.CompanyName,
		}

		if err := tx.Create(&alumni).Error; err != nil {
			return fmt.Errorf("failed to create alumni, rolling back transaction")
		}

		return nil // Commit the transaction
	})

	// Handle the final outcome of the transaction
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	// Return a success response
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Alumni created successfully"})
}

func (h *AlumniHandler) ViewAlumniInfo(c *gin.Context) {
	// Get the alumni ID from the URL parameter, if provided
	alumniID := c.Param("id")

	// If an ID is provided, retrieve the specific alumni
	var alumni models.Alumni
	if err := h.DB.Preload("User").First(&alumni, "id = ?", alumniID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Alumni not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	// Return the alumni information
	c.JSON(http.StatusOK, gin.H{"success": true, "alumni": alumni.GetOneAlumniInfo()})
}

func (h *AlumniHandler) AdminViewAlumniInfo(c *gin.Context) {
	// If no ID is provided, return all alumni (admin-only operation)
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	username := req.Username
	password := req.Password

	var user models.User
	if err := h.DB.First(&user, "Username = ?", username).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Record not found"})
	}
	if user.Role != "admin" || password != user.Password {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Failed Admin Password or not Admin"})
		return
	}

	// retrieve all alumni Information
	alumniList, err := models.GetAllAlumniInfo(h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
	}

	// Return the alumni list
	c.JSON(http.StatusOK, gin.H{"success": true, "alumni": alumniList})
	return

}

func (h *AlumniHandler) UpdateAlumniInfo(c *gin.Context) {
	alumniID := c.Param("id")

	// Fetch the existing alumni record
	var alumni models.Alumni
	if err := h.DB.First(&alumni, "id = ?", alumniID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Alumni not found"})
		return
	}

	// Bind the incoming JSON to the AlumniCreate struct
	var alumniUpdate AlumniCreate
	if err := c.ShouldBindJSON(&alumniUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	// print alumniUpdate
	fmt.Println(alumniUpdate)

	// Apply the updates to the existing alumni object
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		res := map[string]interface{}{
			"FirstName":   alumniUpdate.FirstName,
			"LastName":    alumniUpdate.LastName,
			"Email":       alumniUpdate.Email,
			"Username":    alumniUpdate.Username,
			"Password":    alumniUpdate.Password,
			"Role":        alumniUpdate.Role,
			"CompanyName": alumniUpdate.CompanyName,
		}
		if err := tx.Model(&alumni).Updates(res).Error; err != nil {
			return err // Rollback if an error occurs
		}
		return nil // Commit if no error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	// Return the updated alumni record
	c.JSON(http.StatusOK, gin.H{"success": true, "data": alumni})
}
