package models

import (
	"time"
)

type Todo struct {
	ID          uint   `gorm:"primaryKey;autoIncrement;type:BIGINT UNSIGNED" json:"id"`
	Title       string `gorm:"size:200;not null" json:"title"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	Status      string `gorm:"type:ENUM('pending','in_progress','completed','cancelled');default:'pending';index" json:"status"`
	Urgency     string `gorm:"type:ENUM('low','medium','high');default:'medium'" json:"urgency"`
	Category    string `gorm:"type:ENUM('work','study','fun','personal');default:'personal'" json:"category"`

	StartTime   time.Time `gorm:"not null;index" json:"start_time"`
	EndTime     time.Time `gorm:"not null" json:"end_time"`
	CompletedAt time.Time `gorm:"default:'1970-01-01'" json:"completed_at"`

	// 父子关系
	ParentID    uint `gorm:"not null;default:0;index" json:"parent_id"`
	HasChildren bool `gorm:"default:false;index" json:"has_children"`

	// 重复规则
	RepeatType     string    `gorm:"type:ENUM('none','daily','weekly','monthly');default:'none'" json:"repeat_type"`
	RepeatInterval int       `gorm:"default:1" json:"repeat_interval"`
	RepeatEndDate  time.Time `gorm:"default:'2125-01-01'" json:"repeat_end_date"`

	CreatorUserID  uint `gorm:"not null;default:0" json:"creator_user_id"`
	CreatorOrganID uint `gorm:"not null;default:0" json:"creator_organ_id"`
}

// CalculateNextInstance 计算下一个实例时间
func (t *Todo) CalculateNextInstance(afterTime time.Time) (time.Time, bool) {
	if t.RepeatType == "none" {
		return afterTime, false
	}

	current := t.StartTime
	for !current.After(afterTime) {
		switch t.RepeatType {
		case "daily":
			current = current.AddDate(0, 0, t.RepeatInterval)
		case "weekly":
			current = current.AddDate(0, 0, 7*t.RepeatInterval)
		case "monthly":
			current = current.AddDate(0, t.RepeatInterval, 0)
		default:
			return afterTime, false
		}

		if current.After(t.RepeatEndDate) {
			return afterTime, false
		}
	}

	return current, true
}

func (Todo) TableName() string {
	return "todos"
}

type TodoAssignment struct {
	ID         uint      `gorm:"primaryKey;autoIncrement;type:BIGINT UNSIGNED" json:"id"`
	TodoID     uint      `gorm:"not null;type:BIGINT UNSIGNED" json:"todo_id"`
	AssignedTo uint      `gorm:"not null;type:BIGINT UNSIGNED" json:"assigned_to"`
	AssignedBy uint      `gorm:"not null;type:BIGINT UNSIGNED" json:"assigned_by"`
	AssignedAt time.Time `json:"assigned_at"`
	Status     string    `gorm:"type:ENUM('pending','accepted','rejected');default:'pending'" json:"status"`
	AcceptedAt time.Time `json:"accepted_at"`

	// 简化外键约束
	Todo     Todo `gorm:"foreignKey:TodoID" json:"todo"`
	Assignee User `gorm:"foreignKey:AssignedTo" json:"assignee"`
	Assigner User `gorm:"foreignKey:AssignedBy" json:"assigner"`
}

func (TodoAssignment) TableName() string {
	return "todo_assignments"
}
