// cache/organization_cache.go
package cache

import (
	"fmt"
	"strconv"
	"team_task_hub/backend/internal/models"
	"time"
)

// 用户ID映射用户的所有组织名称
func UserOrganizationNamesKey(userID uint) string {
	return fmt.Sprintf("user:org_names:%d", userID)
}

func SetUserOrganizationNames(userID uint, orgNames []string, expiration time.Duration) error {
	return setJson(UserOrganizationNamesKey(userID), orgNames, expiration)
}

func GetUserOrganizationNames(userID uint) ([]string, error) {
	var orgNames []string
	err := getJson(UserOrganizationNamesKey(userID), &orgNames)
	if err != nil {
		return nil, err
	}
	return orgNames, nil
}

func DeleteUserOrganizationNames(userID uint) error {
	return deleteKey(UserOrganizationNamesKey(userID))
}

// 组织ID映射组织信息
func OrganizationKey(orgID uint) string {
	return fmt.Sprintf("organization:%d", orgID)
}

func SetOrganization(org *models.Organization, expiration time.Duration) error {
	return setJson(OrganizationKey(org.ID), org, expiration)
}

func GetOrganization(orgID uint) (*models.Organization, error) {
	var org models.Organization
	err := getJson(OrganizationKey(orgID), &org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func DeleteOrganization(orgID uint) error {
	return deleteKey(OrganizationKey(orgID))
}

// 组织名称映射组织ID
func OrganizationNameIDKey(orgName string) string {
	return fmt.Sprintf("org_name_id:%s", orgName)
}

func SetOrganizationIDByName(orgName string, orgID uint, expiration time.Duration) error {
	return setString(OrganizationNameIDKey(orgName), fmt.Sprintf("%d", orgID), expiration)
}

func GetOrganizationIDByName(orgName string) (uint, error) {
	idStr, err := Client.Get(ctx, OrganizationNameIDKey(orgName)).Result()
	if err != nil {
		return 0, err
	}

	orgID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid organization ID in cache: %v", err)
	}

	return uint(orgID), nil
}

func DeleteOrganizationIDByName(orgName string) error {
	return deleteKey(OrganizationNameIDKey(orgName))
}
