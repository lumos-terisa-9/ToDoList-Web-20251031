package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"team_task_hub/backend/internal/cache"
	"team_task_hub/backend/internal/models"
	"team_task_hub/backend/internal/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService 认证服务
type AuthService struct {
	userRepo     *repositories.UserRepository
	emailService *EmailService
	jwtSecret    string
}

// JWTClaims JWT声明
type JWTClaims struct {
	UserID       uint   `json:"user_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	JTI          string `json:"jti"`
	TokenVersion uint   `json:"token_version"`
	jwt.RegisteredClaims
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo *repositories.UserRepository, emailService *EmailService, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		emailService: emailService,
		jwtSecret:    jwtSecret,
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Userid   uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"` // 邮箱验证码
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// 验证密码
func (s *AuthService) isValidPassword(password string) bool {
	if len(password) < 6 || len(password) > 20 {
		return false
	}
	hasLetter, _ := regexp.MatchString(`[a-zA-Z]`, password)
	hasDigit, _ := regexp.MatchString(`[0-9]`, password)
	return hasLetter && hasDigit
}

// Register 用户注册
func (s *AuthService) Register(req *RegisterRequest) (*RegisterResponse, error) {
	// 验证邮箱验证码
	isValid, _, err := s.emailService.VerifyCode(req.Email, req.Code, "register")
	if err != nil {
		return nil, errors.New("验证码验证失败")
	}
	if !isValid {
		return nil, errors.New("邮箱验证码无效或已过期")
	}

	// 验证密码强度
	if !s.isValidPassword(req.Password) {
		return nil, errors.New("密码必须包含字母和数字，长度6-20位")
	}

	//检查id是否已经存在
	if exists, _ := s.userRepo.ExistsById(req.Userid); exists {
		return nil, errors.New("ID已被使用")
	}

	// 检查用户名是否已存在
	if exists, _ := s.userRepo.ExistsByUsername(req.Username); exists {
		return nil, errors.New("用户名已被使用")
	}

	// 检查邮箱是否已存在
	if exists, _ := s.userRepo.ExistsByEmail(req.Email); exists {
		return nil, errors.New("邮箱已被注册")
	}

	// 加密密码
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	//创建用户
	user := &models.User{
		ID:           req.Userid,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("用户创建失败")
	}

	// 标记验证码为已使用
	_ = s.emailService.MarkCodeAsUsed(req.Email, req.Code)

	return &RegisterResponse{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

// LoginRequest 登录请求
type LoginRequest struct {
	Identifier string `json:"identifier"` // 邮箱或用户ID
	Password   string `json:"password"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	UserID      uint   `json:"user_id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 根据标识符查找用户
	user, err := s.findUserByIdentifier(req.Identifier)
	if err != nil {
		return nil, errors.New("账号或密码错误")
	}

	// 验证密码
	if !checkPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("账号或密码错误")
	}

	// 生成JWT令牌
	token, err := s.generateJWT(user)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	return &LoginResponse{
		UserID:      user.ID,
		Username:    user.Username,
		Email:       user.Email,
		AccessToken: token,
	}, nil
}

// 生成随机JTI
func generateJTI() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// generateJWT 生成JWT令牌
func (s *AuthService) generateJWT(user *models.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	// 生成唯一JTI
	jti, err := generateJTI()
	if err != nil {
		return "", fmt.Errorf("生成JWT标识失败: %v", err)
	}

	claims := &JWTClaims{
		UserID:       user.ID,
		Username:     user.Username,
		Email:        user.Email,
		JTI:          jti,
		TokenVersion: user.TokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "team_task_hub",
			Subject:   strconv.FormatUint(uint64(user.ID), 10),
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ValidateJWT 验证JWT令牌
func (s *AuthService) ValidateJWT(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("非预期的签名方法")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("无效的令牌")
	}

	if inBlacklist, _ := cache.IsInJWTBlacklist(claims.JTI); inBlacklist {
		return nil, errors.New("令牌已失效")
	}

	if err := s.validateTokenVersion(claims); err != nil {
		return nil, err
	}

	return claims, nil
}

// validateTokenVersion 验证令牌版本号
func (s *AuthService) validateTokenVersion(claims *JWTClaims) error {
	user, err := s.FindUserByIDWithCache(claims.UserID)
	if err != nil {
		return fmt.Errorf("用户信息验证失败: %v", err)
	}
	if user == nil {
		return errors.New("用户不存在")
	}
	if claims.TokenVersion != user.TokenVersion {
		return fmt.Errorf("令牌版本不匹配（令牌:%d ≠ 当前:%d）",
			claims.TokenVersion, user.TokenVersion)
	}
	return nil
}

// LogoutResponse 登出响应
type LogoutResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// Logout 用户登出（将令牌加入黑名单）
func (s *AuthService) Logout(tokenString string) error {
	// 验证令牌有效性
	claims, err := s.ValidateJWT(tokenString)
	if err != nil {
		return nil
	}

	// 计算令牌剩余有效期
	expirationTime := claims.ExpiresAt.Time
	remainingTime := time.Until(expirationTime)

	if remainingTime > 0 {
		// 将令牌加入黑名单，有效期至原令牌过期
		if err := cache.AddToJWTBlacklist(claims.JTI, remainingTime); err != nil {
			return fmt.Errorf("登出失败: %v", err)
		}
		log.Printf("✅ 用户%d令牌已加入黑名单，剩余有效期: %v", claims.UserID, remainingTime)
	}
	// 如果令牌已过期，无需加入黑名单

	return nil
}

// LogoutUser 通过用户ID登出（登出所有该用户的令牌）
func (s *AuthService) LogoutUser(userID uint) error {
	// 这里可以扩展为记录用户的所有活跃令牌并全部加入黑名单
	// 目前简单实现为记录日志
	log.Printf("✅ 用户%d所有会话已登出", userID)
	return nil
}

// findUserByIdentifier 根据标识符查找用户（缓存优先）
func (s *AuthService) findUserByIdentifier(identifier string) (*models.User, error) {
	identifier = strings.TrimSpace(identifier)

	// 根据格式特征判断类型
	if strings.Contains(identifier, "@") {
		return s.findUserByEmailWithCache(identifier)
	}

	// 尝试解析为数字ID
	if userID, err := strconv.ParseUint(identifier, 10, 32); err == nil {
		return s.FindUserByIDWithCache(uint(userID))
	}

	return nil, errors.New("用户不存在")
}

// findUserByEmailWithCache 通过邮箱查找用户
func (s *AuthService) findUserByEmailWithCache(email string) (*models.User, error) {
	// 通过邮箱缓存获取用户ID
	userID, err := cache.GetUserIDByEmail(email)
	if err == nil {
		// 缓存命中，通过ID获取用户完整信息
		return s.FindUserByIDWithCache(userID)
	}

	// 缓存未命中，查询数据库
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}
	// 将查询结果写入缓存
	// 缓存用户完整信息
	cache.SetUser(user, 30*time.Minute)
	// 缓存邮箱到用户ID的映射
	cache.SetUserByEmail(user.Email, user.ID, 30*time.Minute)

	return user, nil
}

// FindUserByIDWithCache 通过ID查找用户（带缓存）
func (s *AuthService) FindUserByIDWithCache(userID uint) (*models.User, error) {
	// 从缓存获取完整用户信息
	cachedUser, err := cache.GetUser(userID)
	if err == nil && cachedUser != nil {
		return cachedUser, nil
	}

	// 缓存未命中，查询数据库
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}
	// 将查询结果写入缓存
	cache.SetUser(user, 30*time.Minute)
	return user, nil
}

// 更新用户名
func (s *AuthService) UpdateUserName(userID uint, newUserName string) error {
	// 基础验证
	newUserName = strings.TrimSpace(newUserName)
	if len(newUserName) == 0 {
		return fmt.Errorf("用户名不能为空")
	}

	// 检查新用户名是否已被其他用户占用
	exists, err := s.userRepo.ExistsByUsername(newUserName)
	if err != nil {
		return fmt.Errorf("检查用户名失败: %v", err)
	}
	if exists {
		return fmt.Errorf("用户名 '%s' 已被占用", newUserName)
	}

	// 执行更新
	if err := s.userRepo.UpdateUserName(userID, newUserName); err != nil {
		return fmt.Errorf("更新用户名失败: %v", err)
	}
	cache.DeleteUser(userID)
	return nil
}

// 更新头像
func (s *AuthService) UpdateAvatar(userID uint, newAvatar string) error {
	if err := s.userRepo.UpdateAvatar(userID, newAvatar); err != nil {
		return errors.New("更新头像失败")
	}
	cache.DeleteUser(userID)
	return nil
}

// 更新用户邮箱
func (s *AuthService) UpdateEmail(userID uint, oldEmail, newEmail, oldEmailCode, newEmailCode string) error {
	//验证旧邮箱
	if isValidOld, _, err := s.emailService.VerifyCode(oldEmail, oldEmailCode, "change_email_old"); err != nil {
		log.Printf("验证旧邮箱验证码时发生错误: %v", err)
		return fmt.Errorf("验证服务暂时不可用")
	} else if !isValidOld {
		return fmt.Errorf("旧邮箱验证码错误或已过期")
	}
	//验证新邮箱
	if isValidNew, _, err := s.emailService.VerifyCode(newEmail, newEmailCode, "change_email_new"); err != nil {
		log.Printf("验证新邮箱验证码时发生错误: %v", err)
		return fmt.Errorf("验证服务暂时不可用")
	} else if !isValidNew {
		return fmt.Errorf("新邮箱验证码错误或已过期")
	}

	if err := s.userRepo.UpdateEmail(userID, newEmail); err != nil {
		return fmt.Errorf("更新用户邮箱失败")
	}
	cache.DeleteUser(userID)
	return nil
}

// 更新用户密码
func (s *AuthService) UpdatePassword(userID uint, email, code, newPassword string) error {
	//验证密码强度
	if !s.isValidPassword(newPassword) {
		return fmt.Errorf("密码必须包含字母和数字，长度6-20位")
	}

	//验证验证码
	if isValid, _, err := s.emailService.VerifyCode(email, code, "change_password"); err != nil {
		return fmt.Errorf("验证服务暂时不可用")
	} else if !isValid {
		return fmt.Errorf("更改密码验证码无效或者过期")
	}

	//密码加密
	passwordHash, err := hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("密码加密失败")
	}

	//更新密码
	if err := s.userRepo.UpdatePassword(userID, passwordHash); err != nil {
		return fmt.Errorf("更新密码失败")
	}
	cache.DeleteUser(userID)
	return nil
}
