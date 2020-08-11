package mysql

import (
	"database/sql"
	"go_web_bbs/models"

	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(post *models.Post) (err error) {
	sqlStr := `insert into post(post_id, title, content, author_id, community_id)
    values (?, ?, ?, ?, ?)`
	_, err = db.Exec(sqlStr, post.PostID, post.Title,
		post.Content, post.AuthorID, post.CommunityID)
	if err != nil {
		zap.L().Error("insert post failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}

//GetPostByID 根据id获取帖子详情
func GetPostByID(postId int64) (data *models.Post, err error) {
	data = new(models.Post)
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
    from post
    where post_id = ?`
	err = db.Get(data, sqlStr, postId)
	if err == sql.ErrNoRows {
		err = ErrorInvalidID
		return
	}
	if err != nil {
		zap.L().Error("query post failed", zap.String("sql", sqlStr), zap.Error(err))
		err = ErrorQueryFailed
		return
	}
	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (postList []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
    from post
    limit ?, ?`

	postList = make([]*models.Post, 0, 1) // 不要写成make([]*models.Post, 2)
	err = db.Select(&postList, sqlStr, (page-1)*size, size)
	return
}
