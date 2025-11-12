package models

import (
	"time"
)

type Todo struct {
	ID           uint      `gorm:"primaryKey;autoIncrement;type:BIGINT UNSIGNED" json:"id"`
	Title        string    `gorm:"size:200;not null" json:"title"`
	Description  string    `gorm:"type:text" json:"description"`
	UrgencyLevel string    `gorm:"type:ENUM('low','medium','high');default:'medium'" json:"urgency_level"`
	TodoType     string    `gorm:"type:ENUM('personal','organizational');not null" json:"todo_type"`
	CreatorID    uint      `gorm:"not null;type:BIGINT UNSIGNED" json:"creator_id"`
	Status       string    `gorm:"type:ENUM('pending','in_progress','completed','cancelled');default:'pending'" json:"status"`
	DueDate      time.Time `json:"due_date"`
	CompletedAt  time.Time `json:"completed_at"`
	CreatedAt    time.Time `json:"created_at"`
	AIGenerated  bool      `gorm:"default:false" json:"ai_generated"`

	// 简化外键约束
	OrganizationID *uint         `gorm:"type:BIGINT UNSIGNED" json:"organization_id"`
	Organization   *Organization `gorm:"foreignKey:OrganizationID" json:"organization"`
	Creator        User          `gorm:"foreignKey:CreatorID" json:"creator"`
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
