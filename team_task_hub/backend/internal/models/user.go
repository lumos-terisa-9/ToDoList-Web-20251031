package models

type User struct {
	ID           uint   `gorm:"primaryKey;autoIncrement;type:BIGINT UNSIGNED" json:"id"`
	Username     string `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email        string `gorm:"size:100;uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"size:255;not null" json:"password_hash"`
	TokenVersion uint   `gorm:"default:1;not null" json:"token_version"`
	AvatarURL    string `gorm:"size:255" json:"avatar_url"`
	Role         uint   `gorm:"type:TINYINT UNSIGNED;default:0;not null" json:"role"`
}

func (User) TableName() string {
	return "users"
}
