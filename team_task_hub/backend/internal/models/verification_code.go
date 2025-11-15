package models

import (
	"time"
)

type VerificationCode struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;type:BIGINT UNSIGNED" json:"id"`
	Email     string    `gorm:"size:100;not null;index" json:"email"`
	Code      string    `gorm:"size:10;not null" json:"code"`
	Business  string    `gorm:"size:20;not null;default:'general'" json:"business"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Used      bool      `gorm:"default:false;not null" json:"used"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (VerificationCode) TableName() string {
	return "verification_codes"
}
