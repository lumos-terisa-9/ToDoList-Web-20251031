package cache

import (
	"fmt"
	"time"
)

// 组织成员身份缓存键定义
func OrganizationCreatorKey(orgID uint) string {
	return fmt.Sprintf("org:%d:creator", orgID)
}

func OrganizationAdminKey(orgID, userID uint) string {
	return fmt.Sprintf("org:%d:admin:%d", orgID, userID)
}

func OrganizationMemberKey(orgID, userID uint) string {
	return fmt.Sprintf("org:%d:member:%d", orgID, userID)
}

// 组织创建者缓存操作
func SetOrganizationCreator(orgID, creatorID uint, expiration time.Duration) error {
	return setJson(OrganizationCreatorKey(orgID), creatorID, expiration)
}

func GetOrganizationCreator(orgID uint) (uint, error) {
	var creatorID uint
	err := getJson(OrganizationCreatorKey(orgID), &creatorID)
	return creatorID, err
}

func DeleteOrganizationCreator(orgID uint) error {
	return deleteKey(OrganizationCreatorKey(orgID))
}

// 组织管理员状态缓存操作
func SetOrganizationAdmin(orgID, userID uint, isAdmin bool, expiration time.Duration) error {
	return setJson(OrganizationAdminKey(orgID, userID), isAdmin, expiration)
}

func GetOrganizationAdmin(orgID, userID uint) (bool, error) {
	var isAdmin bool
	err := getJson(OrganizationAdminKey(orgID, userID), &isAdmin)
	return isAdmin, err
}

func DeleteOrganizationAdmin(orgID, userID uint) error {
	return deleteKey(OrganizationAdminKey(orgID, userID))
}

// 组织成员状态缓存操作
func SetOrganizationMember(orgID, userID uint, isMember bool, expiration time.Duration) error {
	return setJson(OrganizationMemberKey(orgID, userID), isMember, expiration)
}

func GetOrganizationMember(orgID, userID uint) (bool, error) {
	var isMember bool
	err := getJson(OrganizationMemberKey(orgID, userID), &isMember)
	return isMember, err
}

func DeleteOrganizationMember(orgID, userID uint) error {
	return deleteKey(OrganizationMemberKey(orgID, userID))
}

// 批量失效用户组织缓存
func InvalidateUserOrganizationCache(orgID, userID uint) error {
	keys := []string{
		OrganizationAdminKey(orgID, userID),
		OrganizationMemberKey(orgID, userID),
	}

	for _, key := range keys {
		if err := deleteKey(key); err != nil {
			return err
		}
	}
	return nil
}
