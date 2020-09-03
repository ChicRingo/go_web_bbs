package logic

import (
	"fmt"
	"go_web_bbs/dao/mysql"
	"go_web_bbs/dao/redis"
	"go_web_bbs/models"
	"go_web_bbs/pkg/snowflake"

	"go.uber.org/zap"
)

// 创建帖子
func CreatePost(post *models.Post) (err error) {
	// 生成帖子ID
	post.PostID = snowflake.GenID()
	// 创建帖子
	if err := mysql.CreatePost(post); err != nil {
		zap.L().Error("mysql.CreatePost(post) failed", zap.Error(err))
		return err
	}
	err = redis.CreatePost(post)
	return
}

// 根据帖子id获取帖子详情数据
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

// 根据页码和每页个数获取帖子分页列表
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

// 根据页码和每页个数获取帖子分页列表
func GetPostListByOrder(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {

	// 从redis中查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}

	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}

	zap.L().Debug("redis.GetPostIDsInOrder(p) ids:", zap.Any("ids", ids))

	// 根据从redis中获取的ids去mysql中查询帖子详情
	// 返回的数据还要按照我给定的id的顺序返回
	postList, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("mysql.GetPostListByIDs(ids) postList:", zap.Any("postList", postList))

	// 提前查询好每篇帖子的投票数
	data = make([]*models.ApiPostDetail, 0, len(postList))

	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 遍历postList数据，将帖子的作者和分区信息查询出来填充到帖子中
	for idx, post := range postList {
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
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}

	// 返回数据
	return
}

func GetCommunityPostList(p *models.ParamCommunityList) (data []*models.ApiPostDetail, err error) {

	// 从redis中查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}

	if len(ids) == 0 {
		zap.L().Warn("redis.GetCommunityPostIDsInOrder(p) return 0 data")
		return
	}
	zap.L().Debug("redis.GetCommunityPostIDsInOrder(p) ids:", zap.Any("ids", ids))

	// 根据从redis中获取的ids去mysql中查询帖子详情
	// 返回的数据还要按照给定的id的顺序返回
	postList, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("mysql.GetPostListByIDs(ids) postList:", zap.Any("postList", postList))

	// 提前查询好每篇帖子的投票数
	data = make([]*models.ApiPostDetail, 0, len(postList))

	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 遍历postList数据，将帖子的作者和分区信息查询出来填充到帖子中
	for idx, post := range postList {
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
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}

	// 返回数据
	return
}
