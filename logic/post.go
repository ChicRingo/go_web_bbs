package logic

import (
	"go_web_bbs/dao/mysql"
	"go_web_bbs/models"
	"go_web_bbs/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(post *models.Post) (err error) {
	// 生成帖子ID
	postID := snowflake.GetID()
	post.PostID = postID
	// 创建帖子
	if err := mysql.CreatePost(post); err != nil {
		zap.L().Error("mysql.CreatePost(&post) failed", zap.Error(err))
		return err
	}
	//community, err := mysql.GetCommunityNameByID(fmt.Sprint(post.CommunityID))
	//if err != nil {
	//	zap.L().Error("mysql.GetCommunityNameByID failed", zap.Error(err))
	//	return err
	//}
	//if err := redis.CreatePost(
	//	fmt.Sprint(post.PostID),
	//	fmt.Sprint(post.AuthorId),
	//	post.Title,
	//	TruncateByWords(post.Content, 120),
	//	community.CommunityName); err != nil {
	//	zap.L().Error("redis.CreatePost failed", zap.Error(err))
	//	return err
	//}
	return

}

//func GetPost(postID string) (post *models.ApiPostDetail, err error) {
//	post, err = mysql.GetPostByID(postID)
//	if err != nil {
//		zap.L().Error("mysql.GetPostByID(postID) failed", zap.String("post_id", postID), zap.Error(err))
//		return nil, err
//	}
//	user, err := mysql.GetUserByID(fmt.Sprint(post.AuthorId))
//	if err != nil {
//		zap.L().Error("mysql.GetUserByID() failed", zap.String("author_id", fmt.Sprint(post.AuthorId)), zap.Error(err))
//		return
//	}
//	post.AuthorName = user.UserName
//	community, err := mysql.GetCommunityByID(fmt.Sprint(post.CommunityID))
//	if err != nil {
//		zap.L().Error("mysql.GetCommunityByID() failed", zap.String("community_id", fmt.Sprint(post.CommunityID)), zap.Error(err))
//		return
//	}
//	post.CommunityName = community.CommunityName
//	return post, nil
//}
