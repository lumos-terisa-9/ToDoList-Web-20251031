package services

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// hashPassword 密码加密（业务逻辑）
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// checkPassword 密码验证（业务逻辑）
func checkPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// 创建今天的时间范围
func TodayRange() (start, end time.Time) {
	now := time.Now()
	start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end = start.Add(24 * time.Hour)
	return start, end
}

func DayRange(data time.Time) (start, end time.Time) {
	var dayStart time.Time
	var dayEnd time.Time
	dayStart = time.Date(data.Year(), data.Month(), data.Day(), 0, 0, 0, 0, data.Location())
	dayEnd = start.Add(24 * time.Hour)
	return dayStart, dayEnd
}

// 创建本周的时间范围（周一到周日）
func ThisWeekRange() (start, end time.Time) {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	start = time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, now.Location())
	end = start.AddDate(0, 0, 7)
	return start, end
}

// 创建未来N天的时间范围
func NextDaysRange(days int) (start, end time.Time) {
	now := time.Now()
	start = now
	end = now.AddDate(0, 0, days)
	return start, end
}

// 创建特定月份的时间范围
func MonthRange(year int, month time.Month) (start, end time.Time) {
	start = time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	end = time.Date(year, month+1, 1, 0, 0, 0, 0, time.Local)
	return start, end
}

// parseDateString 解析日期字符串（优化版：使用单一标准格式）
func ParseDateString(dateStr string) (time.Time, error) {
	// 定义系统规定的标准日期格式
	const standardDateFormat = "2006-01-02"

	// 清理输入字符串可能存在的空格
	dateStr = strings.TrimSpace(dateStr)

	// 直接使用规定的标准格式进行解析
	parsedTime, err := time.Parse(standardDateFormat, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("无法解析日期 '%s'，要求格式为 YYYY-MM-DD (如: 2024-01-15)", dateStr)
	}

	return parsedTime, nil
}
