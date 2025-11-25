package models

type Organization struct {
	ID           uint   `gorm:"primaryKey;autoIncrement;type:BIGINT UNSIGNED" json:"id"`
	Name         string `gorm:"size:100;not null" json:"name"`
	Description  string `gorm:"type:text" json:"description"`
	LogoURL      string `gorm:"size:255" json:"logo_url"`
	CreatorID    uint   `gorm:"not null;type:BIGINT UNSIGNED;index" json:"creator_id"`
	LocationCode string `gorm:"size:6;uniqueIndex" json:"location_code"`
	Status       string `gorm:"type:ENUM('pending','approved','rejected');default:'pending'" json:"status"`

	// 简化外键约束
	Creator User `gorm:"foreignKey:CreatorID" json:"creator"`
}

func (Organization) TableName() string {
	return "organizations"
}

type OrganizationMember struct {
	ID             uint   `gorm:"primaryKey;autoIncrement;type:BIGINT UNSIGNED" json:"id"`
	OrganizationID uint   `gorm:"not null;type:BIGINT UNSIGNED" json:"organization_id"`
	UserID         uint   `gorm:"not null;type:BIGINT UNSIGNED;index" json:"user_id"`
	Role           string `gorm:"type:ENUM('creator','admin','member');not null" json:"role"`
	Score          int    `gorm:"default:0" json:"score"`
	Status         string `gorm:"type:ENUM('active','inactive','quit');default:'active'" json:"status"`

	// 简化外键约束
	User         User         `gorm:"foreignKey:UserID" json:"user"`
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization"`
}

func (OrganizationMember) TableName() string {
	return "organization_members"
}
