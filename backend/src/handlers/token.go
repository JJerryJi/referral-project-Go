package handlers

import (
	"encoding/hex"
	"github.com/JerryJi/referral-finder/src/models"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Token struct {
	RedisClient *redis.Client
	DB          *gorm.DB
}

func NewTokenHandler(rdb *redis.Client, db *gorm.DB) *Token {
	return &Token{rdb, db}
}

// GET: Retrieve user information using a token

func (h *Token) Retrieve(c *gin.Context) {

	ctx := c.Request.Context()
	curToken := c.Query("token")

	if curToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided in Query"})
		return
	}

	idStr, err := h.RedisClient.Get(ctx, curToken).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session Time Out! Please log in again"})
		return
	}

	userID, err := strconv.Atoi(idStr)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
		return
	}

	// Retrieve the user role from the main users table
	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		return
	}

	// Determine which table to search based on the role
	if user.Role == "student" {
		var student models.Student
		if err := h.DB.First(&student, "user_id = ?", userID).Error; err == nil {
			c.JSON(http.StatusOK, gin.H{
				"student_id": student.ID,
				"user_id":    student.UserID,
				"username":   user.Username,
				"email":      user.Email,
			})
			return
		}
	} else if user.Role == "alumni" {
		var alumni models.Alumni
		if err := h.DB.First(&alumni, "user_id = ?", userID).Error; err == nil {
			c.JSON(http.StatusOK, gin.H{
				"alumni_id": alumni.ID,
				"user_id":   alumni.UserID,
				"username":  user.Username,
				"email":     user.Email,
			})
			return
		}
	}

	// If no match is found in the respective table
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})

}

// POST: Authenticate user and generate a token

func (h *Token) Generate(c *gin.Context) {
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	if err := h.DB.Where(&models.User{
		Username: loginData.Username,
		Password: loginData.Password,
	}).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token := generateToken()

	// Store the token in Redis with a 1-hour expiration
	ctx := c.Request.Context()
	err := h.RedisClient.Set(ctx, token, user.ID, time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Helper function to generate a unique token
func generateToken() string {
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		panic("Failed to generate token")
	}
	return hex.EncodeToString(bytes)
}
