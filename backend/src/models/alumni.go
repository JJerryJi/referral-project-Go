package models

import (
	"gorm.io/gorm"
)

// Alumni model
type Alumni struct {
	gorm.Model         // Includes ID, CreatedAt, UpdatedAt, DeletedAt
	UserID      *uint  `gorm:"unique"`                                         // Foreign key to User table
	User        User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // One-to-One relationship
	CompanyName string `gorm:"size:255;not null"`
}

// GetAllAlumniInfo retrieves all alumni information
func GetAllAlumniInfo(db *gorm.DB) ([]map[string]interface{}, error) {
	var alumniList []Alumni
	if err := db.Preload("User").Find(&alumniList).Error; err != nil {
		return nil, err
	}

	result := []map[string]interface{}{}
	for _, alumni := range alumniList {
		result = append(result, alumni.GetOneAlumniInfo())
	}
	return result, nil
}

// GetOneAlumniInfo retrieves specific alumni information
func (a *Alumni) GetOneAlumniInfo() map[string]interface{} {
	return map[string]interface{}{
		"alumni_id":     a.ID,
		"user":          a.User,
		"company_name":  a.CompanyName,
		"created_time":  a.Model.CreatedAt,
		"modified_time": a.Model.UpdatedAt,
	}
}
