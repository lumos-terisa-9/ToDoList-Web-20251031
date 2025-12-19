package models

import (
	"time"
)

type Activity struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement;type:BIGINT UNSIGNED" json:"id"`
	Title              string    `gorm:"size:200;not null" json:"title"`
	Description        string    `gorm:"type:text" json:"description"`
	OrganizationID     uint      `gorm:"not null;type:BIGINT UNSIGNED;index" json:"organization_id"`
	StartTime          time.Time `gorm:"not null" json:"start_time"`
	EndTime            time.Time `gorm:"not null" json:"end_time"`
	ParticipationLimit string    `gorm:"type:ENUM('public','org_only','admin_assign');not null;default:'org_only'" json:"participation_limit"`

	Status             string    `gorm:"type:ENUM('active','completed','cancelled');not null;default:'active'" json:"status"`
	CompletedAt        time.Time `gorm:"index" json:"completed_at,omitempty"`
	CancellationReason string    `gorm:"type:text" json:"cancellation_reason,omitempty"`
	// 简化外键约束
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization"`
}

func (Activity) TableName() string {
	return "activities"
}

type ActivityParticipation struct {
	ID         uint   `gorm:"primaryKey;autoIncrement;type:BIGINT UNSIGNED" json:"id"`
	ActivityID uint   `gorm:"not null;type:BIGINT UNSIGNED;index" json:"activity_id"`
	UserID     uint   `gorm:"not null;type:BIGINT UNSIGNED;index" json:"user_id"`
	Status     string `gorm:"type:ENUM('pending', 'completed', 'cancelled');not null;default:'pending'" json:"status"`
	IsUnread   bool   `gorm:"default:false;not null" json:"is_unread"`

	Rating     int    `gorm:"check:rating >= 1 AND rating <= 10" json:"rating"`
	ReviewText string `gorm:"type:text" json:"review_text"`

	// 简化外键约束
	Activity Activity `gorm:"foreignKey:ActivityID" json:"activity"`
	User     User     `gorm:"foreignKey:UserID" json:"user"`
}

func (ActivityParticipation) TableName() string {
	return "activity_participation"
}
