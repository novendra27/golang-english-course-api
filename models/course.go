package models

import (
	"time"
)

// Course merepresentasikan data katalog kursus bahasa Inggris
type Course struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `gorm:"type:decimal(12,2);not null;default:0" json:"price"`
	Duration    string    `gorm:"type:varchar(50);not null" json:"duration"` // Contoh: "3 Bulan", "48 Jam"
	Status      string    `gorm:"type:varchar(20);not null;default:'active'" json:"status"` // 'active', 'inactive'
	Classes     []Class   `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"classes,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
