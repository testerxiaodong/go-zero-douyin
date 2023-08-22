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
