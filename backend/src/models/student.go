package models

import (
	"gorm.io/gorm"
	"time"
)

// Student model
type Student struct {
	gorm.Model
	UserID         *uint     `gorm:"unique"`                                         // Nullable to allow SET NULL
	User           User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // One-to-One relationship
	School         string    `gorm:"size:32;not null"`
	YearInSchool   uint      `gorm:"not null"`
	Major          string    `gorm:"size:32;not null"`
	Degree         string    `gorm:"size:3;not null"`
	GraduationYear time.Time `gorm:"not null"`
}

// DEGREES: Predefined choices for Degree field
var DEGREES = map[string]string{
	"BS": "Bachelor of Science",
	"MS": "Master of Science",
}

// GetAllStudentInfo retrieves all student information
func GetAllStudentInfo(db *gorm.DB) ([]map[string]interface{}, error) {
	var studentList []Student
	if err := db.Preload("User").Find(&studentList).Error; err != nil {
		return nil, err
	}

	// Format the output as a slice of maps
	result := []map[string]interface{}{}
	for _, student := range studentList {
		result = append(result, student.GetOneStudentInfo())
	}
	return result, nil
}

// GetOneStudentInfo retrieves specific student information
func (s *Student) GetOneStudentInfo() map[string]interface{} {
	degreeName := DEGREES[s.Degree]
	return map[string]interface{}{
		"student_id":      s.ID,
		"user":            map[string]interface{}{"first_name": s.User.FirstName, "last_name": s.User.LastName, "email": s.User.Email},
		"school":          s.School,
		"year_in_school":  s.YearInSchool,
		"major":           s.Major,
		"degree":          degreeName,
		"graduation_year": s.GraduationYear,
		"created_time":    s.Model.CreatedAt,
		"modified_time":   s.Model.UpdatedAt,
	}
}
