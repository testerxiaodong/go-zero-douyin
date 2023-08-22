package xconst

const (
	VideoFeedCount = 10
)

// 文件资源类型
const (
	TEXT    int64 = 0
	FILE    int64 = 1
	VIDEO   int64 = 2
	PICTURE int64 = 3
	UNKONWN int64 = 9999
)

// RedisKeyPrefix
const (
	RedisVideoCommentPrefix     = "video_id:comment_count:"
	RedisUserLikeVideoPrefix    = "user_id:video_id:"
	RedisVideoLikedByUserPrefix = "video_id:user_id:"
)

const (
	RedisExpireTime = 600
)
