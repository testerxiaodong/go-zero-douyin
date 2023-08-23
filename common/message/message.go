package message

type VideoCommentMessage struct {
	VideoId int64 `json:"video_id"`
}

type UserLikeVideoMessage struct {
	UserId int64 `json:"user_id"`
}

type VideoLikedByUserMessage struct {
	VideoId int64 `json:"video_id"`
}

type UserFollowUserMessage struct {
	FollowerId int64 `json:"follower_id"`
	UserId     int64 `json:"user_id"`
}

type UserFollowedByUserMessage struct {
	UserId     int64 `json:"user_id"`
	FollowerId int64 `json:"follower_id"`
}
