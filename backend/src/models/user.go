package models

import (
	"errors"
	"gorm.io/gorm"
)

// User model
type User struct {
	gorm.Model        // Includes ID, CreatedAt, UpdatedAt, DeletedAt
	Username   string `gorm:"size:500;unique"`
	Password   string `gorm:"size:255"`
	FirstName  string `gorm:"size:50;not null"`
	LastName   string `gorm:"size:30;not null"`
	Email      string `gorm:"size:100;uniqueIndex;not null"`
	Role       string `gorm:"size:10;not null"`
	Location   string `gorm:"size:255;not null"`
}

// ROLES: Predefined choices for the Role field
var ROLES = []string{"student", "alumni", "admin"}

// BeforeCreate GORM hook to validate role
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if !isValidRole(u.Role) {
		return errors.New("validation error: the input value for role is not in the list of ['student', 'alumni', 'admin']")
	}
	return nil
}

// Helper function to validate roles
func isValidRole(role string) bool {
	for _, validRole := range ROLES {
		if role == validRole {
			return true
		}
	}
	return false
}
