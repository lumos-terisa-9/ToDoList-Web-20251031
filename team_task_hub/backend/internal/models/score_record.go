package models

import (
	"time"
)

type ScoreRecord struct {
	ID                uint      `gorm:"primaryKey;autoIncrement;type:BIGINT UNSIGNED" json:"id"`
	UserID            uint      `gorm:"not null;type:BIGINT UNSIGNED" json:"user_id"`
	OrganizationID    uint      `gorm:"not null;type:BIGINT UNSIGNED" json:"organization_id"`
	ScoreChange       int       `gorm:"not null" json:"score_change"`
	RelatedActivityID uint      `gorm:"type:BIGINT UNSIGNED" json:"related_activity_id"`
	Review            string    `gorm:"type:text" json:"review"`
	RecordedAt        time.Time `json:"recorded_at"`

	// 简化外键约束
	User            User         `gorm:"foreignKey:UserID" json:"user"`
	Organization    Organization `gorm:"foreignKey:OrganizationID" json:"organization"`
	RelatedActivity Activity     `gorm:"foreignKey:RelatedActivityID" json:"related_activity"`
}

func (ScoreRecord) TableName() string {
	return "score_records"
}
