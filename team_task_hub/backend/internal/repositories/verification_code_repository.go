package repositories

import (
	"errors"
	"team_task_hub/backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type VerificationCodeRepository struct {
	db *gorm.DB
}

func NewVerificationCodeRepository(db *gorm.DB) *VerificationCodeRepository {
	return &VerificationCodeRepository{db: db}
}

// Create 保存新的验证码记录
func (r *VerificationCodeRepository) Create(code *models.VerificationCode) error {
	return r.db.Create(code).Error
}

// FindValidCode 查找指定邮箱下最新、未使用且未过期的验证码
func (r *VerificationCodeRepository) FindValidCode(email, code, business string) (*models.VerificationCode, error) {
	var validCode models.VerificationCode
	err := r.db.
		Where("email = ? AND code = ? AND used = ? AND business=? AND expires_at > ?", email, code, false, business, time.Now()).
		Order("created_at DESC").
		First(&validCode).Error
	return &validCode, err
}

// MarkAsUsed 将某个验证码标记为已使用
func (r *VerificationCodeRepository) MarkAsUsed(email, code string) error {
	return r.db.Model(&models.VerificationCode{}).Where("email=? AND code=?", email, code).Update("used", true).Error
}

// DeleteExpiredCodes 清理过期验证码（可作为定时任务或者管理员手动删除）
func (r *VerificationCodeRepository) DeleteExpiredCodes() error {
	return r.db.Where("expires_at <= ?", time.Now()).Delete(&models.VerificationCode{}).Error
}

// IsRateLimited 检查发送频率，防止滥用（例如60秒一次，当前设置的就是六十秒）
func (r *VerificationCodeRepository) IsRateLimited(email string, seconds int) (bool, error) {
	var count int64
	sinceTime := time.Now().Add(-time.Duration(seconds) * time.Second)
	err := r.db.Model(&models.VerificationCode{}).
		Where("email = ? AND created_at > ?", email, sinceTime).
		Count(&count).Error
	return count > 0, err
}

// FindValidJoinOrganizationCode 查找有效的组织加入验证码
func (r *VerificationCodeRepository) FindValidJoinOrganizationCode(organizationID uint, code string) (*models.VerificationCode, error) {
	var verificationCode models.VerificationCode

	// 查询未使用、未过期、匹配组织和代码，且业务类型为join_organization的记录
	err := r.db.Where("organization_id = ? AND code = ? AND business = ? AND used = ? AND expires_at > ?",
		organizationID,
		code,
		"join_organization",
		false,
		time.Now()).
		First(&verificationCode).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &verificationCode, nil
}
