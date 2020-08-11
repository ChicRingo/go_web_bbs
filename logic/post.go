package logic

import (
	"fmt"
	"go_web_bbs/dao/mysql"
	"go_web_bbs/models"
	"go_web_bbs/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(post *models.Post) (err error) {
	// 生成帖子ID
	post.PostID = snowflake.GenID()
	// 创建帖子
	if err := mysql.CreatePost(post); err != nil {
		zap.L().Error("mysql.CreatePost(post) failed", zap.Error(err))
		return err
	}
	return
}

// GetPostById 根据帖子id获取帖子详情数据
func GetPostById(postID int64) (data *models.ApiPostDetail, err error) {
	// 查询接口想用的数据
	post, err := mysql.GetPostByID(postID)
	if err != nil {
		zap.L().Error(
			"mysql.GetPostByID(postID) failed",
			zap.Int64("post_id", postID),
			zap.Error(err),
		)
		return nil, err
	}

	// 根据作者id查询作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error(
			"mysql.GetUserByID() failed",
			zap.String("author_id", fmt.Sprint(post.AuthorID)),
			zap.Error(err),
		)
		return
	}

	// 根据社区id查询社区详细信息
	community, err := mysql.GetCommunityDetailById(post.CommunityID)
	if err != nil {
		zap.L().Error(
			"mysql.GetCommunityByID() failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err),
		)
		return
	}

	// 返回组合的数据
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

// GetPostList 根据页码和每页个数获取帖子分页列表
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	postList, err := mysql.GetPostList(page, size)
	if err != nil {
		return
	}
	data = make([]*models.ApiPostDetail, 0, len(postList))

	// 遍历postList数据
	for _, post := range postList {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error(
				"mysql.GetUserByID() failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err),
			)
			continue
		}

		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailById(post.CommunityID)
		if err != nil {
			zap.L().Error(
				"mysql.GetCommunityByID() failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err),
			)
			continue
		}

		// 组合数据
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}

	// 返回数据
	return
}
