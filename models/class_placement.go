package models

import (
	"time"
)

// ClassPlacement merepresentasikan penempatan student ke kelas tertentu setelah pembayaran valid
type ClassPlacement struct {
	ID             uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	RegistrationID uint          `gorm:"not null;uniqueIndex" json:"registration_id"`
	Registration   *Registration `gorm:"foreignKey:RegistrationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"registration,omitempty"`
	ClassID        uint          `gorm:"not null;index" json:"class_id"`
	Class          *Class        `gorm:"foreignKey:ClassID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"class,omitempty"`
	PlacementDate  time.Time     `gorm:"not null" json:"placement_date"`
	CreatedAt      time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}
