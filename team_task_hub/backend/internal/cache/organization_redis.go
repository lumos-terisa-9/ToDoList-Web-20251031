package cache

import (
	"fmt"
	"time"
)

// UserOrganizationOverview 用户组织总览信息
type UserOrganizationOverview struct {
	OrganizationName string    `json:"org_name"`
	CreatorID        uint      `json:"creator_id"`
	LogoURL          string    `json:"logo_url"`
	JoinedAt         time.Time `json:"joined_at"`
}

// UserOrganizationOverviewsKey 用户组织总览缓存键
func UserOrganizationOverviewsKey(userID uint) string {
	return fmt.Sprintf("user:%d:orgs_overview", userID)
}

// SetUserOrganizationOverviews 缓存用户的组织总览列表
func SetUserOrganizationOverviews(userID uint, overviews []UserOrganizationOverview, expiration time.Duration) error {
	return setJson(UserOrganizationOverviewsKey(userID), overviews, expiration)
}

// GetUserOrganizationOverviews 获取用户的组织总览列表
func GetUserOrganizationOverviews(userID uint) ([]UserOrganizationOverview, error) {
	var overviews []UserOrganizationOverview
	err := getJson(UserOrganizationOverviewsKey(userID), &overviews)
	if err != nil {
		return nil, err
	}
	return overviews, nil
}

// DeleteUserOrganizationOverviews 删除用户的组织总览缓存
func DeleteUserOrganizationOverviews(userID uint) error {
	return deleteKey(UserOrganizationOverviewsKey(userID))
}
