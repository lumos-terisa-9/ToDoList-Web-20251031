package services

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"team_task_hub/backend/internal/models"
	"team_task_hub/backend/internal/repositories"
	"time"

	"gopkg.in/gomail.v2"
)

// EmailConfig 邮箱服务配置结构体
type EmailConfig struct {
	Host     string `json:"host"`     // SMTP服务器地址，如 smtp.qq.com
	Port     int    `json:"port"`     // SMTP端口，如 587
	User     string `json:"user"`     // 发件邮箱账号
	Password string `json:"password"` // 邮箱授权码（不是登录密码）
	From     string `json:"from"`     // 发件人显示名称（可选）
}

// EmailService 邮箱服务
type EmailService struct {
	config   *EmailConfig
	codeRepo *repositories.VerificationCodeRepository
}

// NewEmailService 创建邮箱服务实例
func NewEmailService(config *EmailConfig, codeRepo *repositories.VerificationCodeRepository) *EmailService {
	return &EmailService{
		config:   config,
		codeRepo: codeRepo,
	}
}

// SendVerificationCode 发送邮箱验证码
func (s *EmailService) SendVerificationCode(email, businessType string) error {
	limitSeconds := 60

	limited, err := s.codeRepo.IsRateLimited(email, limitSeconds)
	if err != nil {
		return fmt.Errorf("系统错误，请重试: %v", err)
	}
	if limited {
		return fmt.Errorf("请求过于频繁，请%d秒后再试", limitSeconds)
	}

	// 生成验证码
	code, err := s.generateSecureCode()
	if err != nil {
		return fmt.Errorf("生成验证码失败: %v", err)
	}

	// 保存验证码到数据库
	verificationCode := &models.VerificationCode{
		Email:     email,
		Code:      code,
		Business:  businessType,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Used:      false,
	}

	if err := s.codeRepo.Create(verificationCode); err != nil {
		return fmt.Errorf("保存验证码失败: %v", err)
	}

	// 发送邮件
	if err := s.sendEmail(email, code, businessType); err != nil {
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	return nil
}

// VerifyCode 验证邮箱和验证码是否匹配
func (s *EmailService) VerifyCode(email, code, businessType string) (bool, *models.VerificationCode, error) {
	validCode, err := s.codeRepo.FindValidCode(email, code, businessType)
	if err != nil {
		return false, nil, fmt.Errorf("验证码查询失败: %v", err)
	}

	if validCode != nil && validCode.ID != 0 {
		return true, validCode, nil
	}

	return false, nil, nil
}

// MarkCodeAsUsed 标记验证码为已使用
func (s *EmailService) MarkCodeAsUsed(email, code string) error {
	return s.codeRepo.MarkAsUsed(email, code)
}

// generateSecureCode 生成安全的6位数字随机验证码
func (s *EmailService) generateSecureCode() (string, error) {
	const digits = "0123456789"
	code := make([]byte, 6)

	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", fmt.Errorf("生成随机数失败: %v", err)
		}
		code[i] = digits[num.Int64()]
	}

	return string(code), nil
}

// sendEmail 使用SMTP发送邮件
func (s *EmailService) sendEmail(to, code, businessType string) error {
	// 验证配置
	if s.config.Host == "" || s.config.User == "" || s.config.Password == "" {
		return fmt.Errorf("邮箱配置不完整，请检查SMTP配置")
	}

	m := gomail.NewMessage()

	// 设置发件人
	if s.config.From != "" {
		m.SetHeader("From", s.config.From+"<"+s.config.User+">")
	} else {
		m.SetHeader("From", s.config.User)
	}

	// 设置标题
	var subject string
	switch businessType {
	case "register":
		subject = "您的注册验证码 - Team Task Hub"
	case "reset_password":
		subject = "密码重置验证码 - Team Task Hub"
	case "change_email":
		subject = "邮箱修改验证码 - Team Task Hub"
	case "join_org":
		subject = "加入组织验证码 - Team Task Hub"
	default:
		subject = "您的验证码 - Team Task Hub"
	}
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	// 构建HTML邮件内容
	htmlContent := fmt.Sprintf(`
        <!DOCTYPE html>
        <html lang="zh-CN">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>验证码邮件</title>
            <style>
                body {
                    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
                    background-color: #f5f5f5;
                    margin: 0;
                    padding: 20px;
                }
                .container {
                    max-width: 600px;
                    margin: 0 auto;
                    background: white;
                    border-radius: 10px;
                    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
                    overflow: hidden;
                }
                .header {
                    background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
                    color: white;
                    padding: 30px;
                    text-align: center;
                }
                .header h1 {
                    margin: 0;
                    font-size: 24px;
                }
                .content {
                    padding: 30px;
                    color: #333;
                }
                .code {
                    font-size: 32px;
                    font-weight: bold;
                    color: #ff6b6b;
                    text-align: center;
                    margin: 20px 0;
                    letter-spacing: 5px;
                }
                .warning {
                    background: #fff3cd;
                    border-left: 4px solid #ffc107;
                    padding: 15px;
                    margin: 20px 0;
                    border-radius: 4px;
                }
                .footer {
                    background: #f8f9fa;
                    padding: 20px;
                    text-align: center;
                    color: #6c757d;
                    font-size: 12px;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <div class="header">
                    <h1>Team Task Hub</h1>
                    <p>您的注册验证码</p>
                </div>
                <div class="content">
                    <p>您好！</p>
                    <p>您正在注册 Team Task Hub 账户，请使用以下验证码完成注册：</p>
                    
                    <div class="code">%s</div>
                    
                    <div class="warning">
                        <strong>重要提示：</strong>
                        <ul>
                            <li>此验证码 <strong>5分钟</strong> 内有效</li>
                            <li>请勿将此验证码透露给他人</li>
                            <li>如非本人操作，请忽略此邮件</li>
                        </ul>
                    </div>
                    
                    <p>如果这不是您发起的请求，请忽略此邮件。</p>
                </div>
                <div class="footer">
                    <p>此邮件由系统自动发送，请勿回复</p>
                    <p>&copy; 2025 Team Task Hub. 保留所有权利。</p>
                </div>
            </div>
        </body>
        </html>
    `, code)

	// 设置邮件内容为HTML格式
	m.SetBody("text/html", htmlContent)

	textContent := fmt.Sprintf(`
        欢迎注册 Team Task Hub！
        
        您的验证码是：%s
        
        此验证码5分钟内有效，请勿泄露给他人。
        
        如果这不是您发起的请求，请忽略此邮件。
        
        Team Task Hub 团队
    `, code)
	m.AddAlternative("text/plain", textContent)

	// 创建SMTP拨号器并发送邮件
	d := gomail.NewDialer(s.config.Host, s.config.Port, s.config.User, s.config.Password)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("SMTP发送失败: %v", err)
	}

	return nil
}

// TestConnection 测试SMTP连接（用于诊断配置问题）
func (s *EmailService) TestConnection() error {
	if s.config.Host == "" || s.config.User == "" || s.config.Password == "" {
		return fmt.Errorf("邮箱配置不完整")
	}

	d := gomail.NewDialer(s.config.Host, s.config.Port, s.config.User, s.config.Password)

	// 尝试连接但不发送邮件
	sender, err := d.Dial()
	if err != nil {
		return fmt.Errorf("SMTP连接测试失败: %v", err)
	}

	// 关闭连接
	if sender != nil {
		sender.Close()
	}

	return nil
}

// CleanupExpiredCodes 清理过期的验证码
func (s *EmailService) CleanupExpiredCodes() error {
	return s.codeRepo.DeleteExpiredCodes()
}
