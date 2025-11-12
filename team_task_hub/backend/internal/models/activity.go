package models

import (
	"time"
)

type Activity struct {
	ID              uint      `gorm:"primaryKey;autoIncrement;type:BIGINT UNSIGNED" json:"id"`
	Title           string    `gorm:"size:200;not null" json:"title"`
	Description     string    `gorm:"type:text" json:"description"`
	OrganizationID  uint      `gorm:"not null;type:BIGINT UNSIGNED" json:"organization_id"`
	CreatorID       uint      `gorm:"not null;type:BIGINT UNSIGNED" json:"creator_id"`
	ActivityDate    time.Time `gorm:"not null" json:"activity_date"`
	Location        string    `gorm:"size:255" json:"location"`
	MaxParticipants int       `json:"max_participants"`
	Status          string    `gorm:"type:ENUM('upcoming','ongoing','completed','cancelled');default:'upcoming'" json:"status"`
	CreatedAt       time.Time `json:"created_at"`

	// 简化外键约束
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization"`
	Creator      User         `gorm:"foreignKey:CreatorID" json:"creator"`
}

func (Activity) TableName() string {
	return "activities"
}

type ActivityParticipation struct {
	ID                  uint      `gorm:"primaryKey;autoIncrement;type:BIGINT UNSIGNED" json:"id"`
	ActivityID          uint      `gorm:"not null;type:BIGINT UNSIGNED" json:"activity_id"`
	UserID              uint      `gorm:"not null;type:BIGINT UNSIGNED" json:"user_id"`
	ParticipationStatus string    `gorm:"type:ENUM('signed_up','attended','absent','leave');default:'signed_up'" json:"participation_status"`
	LeaveReason         string    `gorm:"type:text" json:"leave_reason"`
	Rating              int       `gorm:"check:rating >= 1 AND rating <= 5" json:"rating"`
	ReviewText          string    `gorm:"type:text" json:"review_text"`
	CheckedInAt         time.Time `json:"checked_in_at"`
	CreatedAt           time.Time `json:"created_at"`

	// 简化外键约束
	Activity Activity `gorm:"foreignKey:ActivityID" json:"activity"`
	User     User     `gorm:"foreignKey:UserID" json:"user"`
}

func (ActivityParticipation) TableName() string {
	return "activity_participation"
}
