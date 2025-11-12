package models

import (
	"time"
)

type Permission struct {
	ID             uint      `gorm:"primaryKey;autoIncrement;type:BIGINT UNSIGNED" json:"id"`
	OrganizationID uint      `gorm:"not null;type:BIGINT UNSIGNED" json:"organization_id"`
	UserID         uint      `gorm:"not null;type:BIGINT UNSIGNED" json:"user_id"`
	PermissionType string    `gorm:"type:ENUM('publish_todo','delete_todo','manage_members','publish_activity','transfer_ownership');not null" json:"permission_type"`
	GrantedBy      uint      `gorm:"not null;type:BIGINT UNSIGNED" json:"granted_by"`
	GrantedAt      time.Time `json:"granted_at"`

	// 简化外键约束
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization"`
	User         User         `gorm:"foreignKey:UserID" json:"user"`
	Granter      User         `gorm:"foreignKey:GrantedBy" json:"granter"`
}

func (Permission) TableName() string {
	return "permissions"
}
