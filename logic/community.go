package logic

import (
	"go_web_bbs/dao/mysql"
	"go_web_bbs/models"
)

func GetCommunityList() ([]*models.Community, error) {
	// 查询数据库 找到所有的community 并返回
	return mysql.GetCommunityList()
}
func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	// 查询数据库 找到所有的community 并返回
	return mysql.GetCommunityDetailById(id)
}
