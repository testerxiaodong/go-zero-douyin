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

// DelStateNo 数据库软删除
const (
	DelStateNo  int64 = 0 //未删除
	DelStateYes int64 = 1 //已删除
)

const (
	FollowStateNo  int64 = 0
	FollowStateYes int64 = 1
)

const (
	LikeStateNo  int64 = 0
	LikeStateYes int64 = 1
)

const (
	ElasticSearchVideoIndexName       = "video"
	ElasticSearchVideoSuggestionName  = "video_suggestion"
	ElasticSearchUserIndexName        = "user"
	ElasticSearchVideoDefaultPageSize = 10
)
