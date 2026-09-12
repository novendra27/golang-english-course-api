// Package models defines the database domain models and GORM schema mapping.
package models

import (
	"time"
)

// Class represents a specific class section offered under a Course.
type Class struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID  uint      `gorm:"not null;index" json:"course_id"`
	Course    *Course   `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"course,omitempty"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"` // e.g. "Basic English A"
	Capacity  int       `gorm:"not null;default:20" json:"capacity"`
	Schedule  string    `gorm:"type:varchar(100);not null" json:"schedule"`            // e.g. "Mon & Wed, 19:00 - 21:00"
	Status    string    `gorm:"type:varchar(20);not null;default:'open'" json:"status"` // 'open', 'full', 'closed'
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

