package models

import (
	"time"
)

// Registration merepresentasikan pendaftaran student ke suatu Course
type Registration struct {
	ID               uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	StudentID        uint            `gorm:"not null;index" json:"student_id"`
	Student          *Student        `gorm:"foreignKey:StudentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"student,omitempty"`
	CourseID         uint            `gorm:"not null;index" json:"course_id"`
	Course           *Course         `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"course,omitempty"`
	RegistrationDate time.Time       `gorm:"not null" json:"registration_date"`
	Status           string          `gorm:"type:varchar(20);not null;default:'pending'" json:"status"` // 'pending', 'registered', 'cancelled', 'completed'
	Payment          *Payment        `gorm:"foreignKey:RegistrationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"payment,omitempty"`
	ClassPlacement   *ClassPlacement `gorm:"foreignKey:RegistrationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"class_placement,omitempty"`
	CreatedAt        time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}
