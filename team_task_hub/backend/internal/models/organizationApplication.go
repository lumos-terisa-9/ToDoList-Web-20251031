package models

type OrganizationApplication struct {
	ID uint `gorm:"primaryKey;autoIncrement;type:BIGINT UNSIGNED" json:"id"`
	// 申请类型与状态
	ApplicationType string `gorm:"type:ENUM('create_org','join_org');not null;index" json:"application_type"`
	Status          string `gorm:"type:ENUM('pending','approved','rejected');default:'pending';index" json:"status"`

	// 申请人信息
	ApplicantUserID   uint   `gorm:"not null;type:BIGINT UNSIGNED;index" json:"applicant_user_id"`
	ApplicantUsername string `gorm:"size:100;not null" json:"applicant_username"`

	OrganizationName string `gorm:"size:100;not null;index" json:"organization_name"`

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
