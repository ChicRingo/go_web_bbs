package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"go_web_bbs/models"

	"go.uber.org/zap"
)

// GetCommunityList
func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := `select community_id, community_name from community`
	err = db.Select(&communityList, sqlStr)
	if err != nil {
		zap.L().Error("mysql.GetCommunityList() failed", zap.Error(err))
	}
	return
}

// GetCommunityDetailByID 根据ID查询社区详情
func GetCommunityDetailByID(id int64) (communityDetail *models.CommunityDetail, err error) {
	communityDetail = new(models.CommunityDetail)
	sqlStr := `select community_id, community_name, introduction, create_time, update_time
    from community
    where community_id = ?`
	err = db.Get(communityDetail, sqlStr, id)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID() failed", zap.Int64("id", id), zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			err = fmt.Errorf("more info: %w [%v]", ErrorInvalidID, id)
		}
	}
	return
}
