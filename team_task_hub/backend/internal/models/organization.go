package models

import "time"

type Organization struct {
	ID           uint   `gorm:"primaryKey;autoIncrement;type:BIGINT UNSIGNED" json:"id"`
	Name         string `gorm:"size:100;not null;index" json:"name"`
	Description  string `gorm:"type:text" json:"description"`
	LogoURL      string `gorm:"size:255" json:"logo_url"`
	CreatorID    uint   `gorm:"not null;type:BIGINT UNSIGNED;index" json:"creator_id"`
	LocationCode string `gorm:"size:9;index" json:"location_code"`

	// 简化外键约束
	Creator User `gorm:"foreignKey:CreatorID" json:"creator"`
}

func (Organization) TableName() string {
	return "organizations"
}

type OrganizationMember struct {
	ID             uint      `gorm:"primaryKey;autoIncrement;type:BIGINT UNSIGNED" json:"id"`
	OrganizationID uint      `gorm:"not null;type:BIGINT UNSIGNED" json:"organization_id"`
	UserID         uint      `gorm:"not null;type:BIGINT UNSIGNED;index" json:"user_id"`
	Role           string    `gorm:"type:ENUM('creator','admin','member');not null;default:'member'" json:"role"`
	Score          int       `gorm:"default:0" json:"score"`
	Status         string    `gorm:"type:ENUM('active','inactive','quit');default:'active'" json:"status"`
	JoinedAt       time.Time `gorm:"autoCreateTime" json:"joined_at"`
	// 简化外键约束
	User         User         `gorm:"foreignKey:UserID" json:"user"`
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization"`
}

func (OrganizationMember) TableName() string {
	return "organization_members"
}

type OrganizationApplication struct {
	ID uint `gorm:"primaryKey;autoIncrement;type:BIGINT UNSIGNED" json:"id"`
	// 申请类型与状态
	ApplicationType string `gorm:"type:ENUM('create_org','join_org','change_org');not null;index" json:"application_type"`
	Status          string `gorm:"type:ENUM('pending','approved','rejected');default:'pending';index" json:"status"`

	// 申请人信息
	ApplicantUserID   uint   `gorm:"not null;type:BIGINT UNSIGNED;index" json:"applicant_user_id"`
	ApplicantUsername string `gorm:"size:100;not null" json:"applicant_username"`

	OrganizationID   uint   `gorm:"type:BIGINT UNSIGNED;index" json:"organization_id"`
	OrganizationName string `gorm:"size:100;index" json:"organization_name"`

	// 申请内容
	ApplicantIntroduction string `gorm:"type:text" json:"applicant_introduction,omitempty"`
	ApplicationReason     string `gorm:"type:text;not null" json:"application_reason"`
	AttachmentURL         string `gorm:"size:255" json:"attachment_url,omitempty"`

	// 审批信息
	AdminRemark string `gorm:"type:text" json:"admin_remark,omitempty"`
}

func (OrganizationApplication) TableName() string {
	return "organization_applications"
}

type OrgInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
