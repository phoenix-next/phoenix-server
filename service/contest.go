package service

import (
	"fmt"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model"
)

// GetReadableContest 获取用户可读的比赛，并按照指定的sorter排序
func GetReadableContest(userID uint64, sorter int) (contests []model.ContestT) {
	var adminOrg []uint64
	var userOrg []uint64
	for _, tmp := range GetAdminOrganization(userID) {
		adminOrg = append(adminOrg, tmp.OrgID)
	}
	for _, tmp := range GetUserOrganization(userID) {
		userOrg = append(userOrg, tmp.OrgID)
	}
	isAscend := "ascend"
	if sorter < 0 {
		isAscend = "descend"
	}
	column := "id"
	if sorter > 1 {
		column = "name"
	}
	global.DB.Model(&model.Contest{}).
		Where("readable = ?", 2).
		Or("readable = ? AND org_id IN ?", 1, userOrg).
		Or("readable = ? AND org_id IN", 0, adminOrg).
		Order(fmt.Sprintf("%s %s", column, isAscend)).
		Find(&contests)
	return
}
