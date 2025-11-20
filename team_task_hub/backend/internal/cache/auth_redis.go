package cache

import (
	"fmt"
	"strconv"
	"team_task_hub/backend/internal/models"
	"time"
)

// 用户信息缓存
func UserKey(userID uint) string {
	return fmt.Sprintf("user:%d", userID)
}

func SetUser(user *models.User, expiration time.Duration) error {
	return setJson(UserKey(user.ID), user, expiration)
}

func GetUser(userID uint) (*models.User, error) {
	var user models.User
	err := getJson(UserKey(userID), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func DeleteUser(userID uint) error {
	return deleteKey(UserKey(userID))
}

// 生成邮箱缓存键
func UserEmailKey(email string) string {
	return fmt.Sprintf("user:email:%s", email)
}

// SetUserByEmail 缓存邮箱到用户ID的映射
func SetUserByEmail(email string, userID uint, expiration time.Duration) error {
	return setString(UserEmailKey(email), fmt.Sprintf("%d", userID), expiration)
}

// GetUserIDByEmail 通过邮箱获取用户ID
func GetUserIDByEmail(email string) (uint, error) {
	idStr, err := Client.Get(ctx, UserEmailKey(email)).Result()
	if err != nil {
		return 0, err
	}

	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID in cache: %v", err)
	}

	return uint(userID), nil
}

// DeleteUserByEmail 删除邮箱缓存映射
func DeleteUserByEmail(email string) error {
	return deleteKey(UserEmailKey(email))
}

// JWT黑名单缓存
func JWTBlacklistKey(jti string) string {
	return fmt.Sprintf("jwt:blacklist:%s", jti)
}

func AddToJWTBlacklist(jti string, expiration time.Duration) error {
	return setString(JWTBlacklistKey(jti), "1", expiration)
}

func IsInJWTBlacklist(jti string) (bool, error) {
	return exists(JWTBlacklistKey(jti))
}

// RemoveFromJWTBlacklist 从黑名单移除JWT
func RemoveFromJWTBlacklist(jti string) error {
	key := JWTBlacklistKey(jti)
	return Client.Del(ctx, key).Err()
}

// GetJWTBlacklistTTL 获取黑名单剩余时间
func GetJWTBlacklistTTL(jti string) (time.Duration, error) {
	key := JWTBlacklistKey(jti)
	return Client.TTL(ctx, key).Result()
}
