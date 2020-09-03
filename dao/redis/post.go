package redis

import (
	"go_web_bbs/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// 1. 根据用户请求中携带的order参数确定要查询的redis key
func GetPostIDsInOrder(p *models.ParamPostList) (ids []string, err error) {
	// 从redis获取id
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	// 2. 确定查询的索引起始点
	return getIDsFormKey(key, p.Page, p.Size)
}

// 根据ids查询每篇帖子的 投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	//data = make([]int64, 0, len(ids))
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZSetPF + id)
	//	//查找key中分数是1的元素的数量-》统计每篇帖子的赞成票的数量
	//	v1 := client.ZCount(key, "1", "1").Val()
	//	data = append(data, v1)
	//}

	// 使用pipeline一次发送多条命令，减少RTT
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}

	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

func getIDsFormKey(key string, page, size int64) (ids []string, err error) {
	start := (page - 1) * size
	end := start + size - 1
	// 3. ZRevRange查询 按分数从大到小查询指定数量的元素
	return client.ZRevRange(key, start, end).Result()
}

// 1. 根据 Community 查询IDs，用户请求中携带的order参数确定要查询的redis key
func GetCommunityPostIDsInOrder(p *models.ParamCommunityList) (ids []string, err error) {

	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}

	// 使用 zinterstore 把区分的帖子set与帖子分数的zset生成一个新的zset
	// 针对新的zset 按之前的逻辑取数据

	// 社区的key
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))

	//利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if client.Exists(key).Val() < 1 {
		// 不存在，需要计算
		pipeline := client.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{ // ZInterStore计算
			Aggregate: "MAX",
		}, cKey, orderKey)
		pipeline.Expire(key, 60*time.Second) // 设置超时时间
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	// 存在则直接根据key查询ids
	return getIDsFormKey(key, p.Page, p.Size)
}
