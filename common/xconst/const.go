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
	RedisVideoCommentPrefix       = "video_id:comment_count:"
	RedisUserLikeVideoPrefix      = "user_id:video_id:"
	RedisVideoLikedByUserPrefix   = "video_id:user_id:"
	RedisUserFollowUserPrefix     = "follower_id:user_id:"
	RedisUserFollowedByUserPrefix = "user_id:follower_id:"
)

// RedisKey
const (
	RedisVideoTag     = "video_tag"
	RedisVideoSection = "video_section"
)

// RedisLockKeyPrefix
const (
	RedisBuildVideoSectionCacheLockKey         = "build_video_section"
	RedisBuildVideoTagCacheLockKey             = "build_video_Tag"
	RedisBuildVideoCommentCountCacheLockPrefix = "build_video_comment_count_key_"
	RedisBuildUserLikeVideoCacheLockPrefix     = "build_user_like_video_id_list_key_"
	RedisBuildVideoLikedByUserCacheLockPrefix  = "build_video_liked_by_user_key_"
	RedisBuildUserFollowCountCacheLockPrefix   = "build_user_follow_count_key_"
	RedisBuildUserFollowerCountCacheLockPrefix = "build_user_follower_count_key_"
)

const (
	RedisExpireTime = 600
)

const (
	ElasticSearchVideoIndexName       = "video"
	ElasticSearchVideoSuggestionName  = "video_suggestion"
	ElasticSearchUserIndexName        = "user"
	ElasticSearchVideoDefaultPageSize = 10
)
