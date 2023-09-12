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
}

type UserFollowedByUserMessage struct {
	UserId int64 `json:"user_id"`
}

type VideoSectionMessage struct {
}

type VideoTagMessage struct {
}
