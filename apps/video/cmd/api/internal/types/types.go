// Code generated by goctl. DO NOT EDIT.
package types

type VideoInfo struct {
	Id           int64  `json:"id"`
	OwnerId      int64  `json:"owner_id"`
	OwnerName    string `json:"owner_name"`
	SectionId    int64  `json:"section_id"`
	TagIds       string `json:"tag_ids"`
	Title        string `json:"title"`
	PlayUrl      string `json:"play_url"`
	CoverUrl     string `json:"cover_url"`
	CommentCount int64  `json:"comment_count"`
	LikeCount    int64  `json:"like_count"`
	CreateTime   int64  `json:"create_time"`
	UpdateTime   int64  `json:"update_time"`
}

type PublishVideoReq struct {
	Title       string `form:"title" validate:"required"`
	SectionId   int64  `form:"section_id" validate:"required"`
	TagIds      string `form:"tag_ids" validate:"required"`
	PublishTime int64  `form:"publish_time,optional"`
}

type VideoFeedReq struct {
	LastTimeStamp int64 `json:"last_time_stamp" validate:"required"`
	SectionId     int64 `json:"section_id" validate:"required"`
}

type VideoFeedResp struct {
	Videos []*VideoInfo `json:"videos"`
}

type UserVideoListReq struct {
	UserId   int64 `json:"user_id" validate:"required"`
	Page     int64 `json:"page,optional"`
	PageSize int64 `json:"page_size,optional"`
}

type UserVideoListResp struct {
	Total  int64        `json:"total"`
	Videos []*VideoInfo `json:"videos"`
}

type VideoDetailReq struct {
	VideoId int64 `json:"video_id"`
}

type VideoDetailResp struct {
	Video *VideoInfo `json:"video"`
}

type DeleteVideoReq struct {
	VideoId int64 `json:"video_id" validate:"required"`
}

type Section struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type AddSectionReq struct {
	Name string `json:"name"`
}

type DelSectionReq struct {
	Id int64 `json:"id"`
}

type GetAllSectionResp struct {
	Sections []*Section `json:"sections"`
}

type Tag struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type AddTagReq struct {
	Name string `json:"name"`
}

type DelTagReq struct {
	Id int64 `json:"id"`
}

type GetAllTagResp struct {
	Tags []*Tag `json:"tags"`
}
