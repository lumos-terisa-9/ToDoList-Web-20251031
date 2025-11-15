package services

import (
	"errors"
	"regexp"
	"strconv"
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
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
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

	// 8. 标记验证码为已使用
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

// findUserByIdentifier 根据标识符查找用户
func (s *AuthService) findUserByIdentifier(identifier string) (*models.User, error) {
	// 尝试解析为数字ID
	if userID, err := strconv.ParseUint(identifier, 10, 32); err == nil {
		// 按ID查找
		return s.userRepo.FindByID(uint(userID))
	}

	// 按邮箱查找
	return s.userRepo.FindByEmail(identifier)
}

// generateJWT 生成JWT令牌
func (s *AuthService) generateJWT(user *models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "team_task_hub",
			Subject:   strconv.FormatUint(uint64(user.ID), 10), // 用户ID作为主题
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

	return claims, nil
}

func (s *AuthService) isValidPassword(password string) bool {
	if len(password) < 6 || len(password) > 20 {
		return false
	}
	hasLetter, _ := regexp.MatchString(`[a-zA-Z]`, password)
	hasDigit, _ := regexp.MatchString(`[0-9]`, password)
	return hasLetter && hasDigit
}
