package models

import (
	"time"
)

// Payment merepresentasikan data pembayaran untuk suatu Registration
type Payment struct {
	ID             uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	RegistrationID uint          `gorm:"not null;uniqueIndex" json:"registration_id"`
	Registration   *Registration `gorm:"foreignKey:RegistrationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"registration,omitempty"`
	Amount         float64       `gorm:"type:decimal(12,2);not null" json:"amount"`
	PaymentMethod  string        `gorm:"type:varchar(50);not null" json:"payment_method"` // 'bank_transfer', 'ewallet', 'credit_card'
	PaymentDate    *time.Time    `json:"payment_date"`
	Status         string        `gorm:"type:varchar(20);not null;default:'pending'" json:"status"` // 'pending', 'paid', 'failed', 'refunded'
	CreatedAt      time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}
