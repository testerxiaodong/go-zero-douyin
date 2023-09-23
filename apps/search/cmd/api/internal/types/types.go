// Code generated by goctl. DO NOT EDIT.
package types

type Video struct {
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
	CreateTime   string `json:"create_time"`
	UpdateTime   string `json:"update_time"`
}

type User struct {
	Id            int64  `json:"id"`
	Username      string `json:"username"`
	FollowerCount int64  `json:"follower_count"`
	FollowCount   int64  `json:"follow_count"`
}

type SearchVideoReq struct {
	Keyword   string `json:"keyword" validate:"required"`
	Page      int64  `json:"page" validate:"required"`
	PageSize  int64  `json:"page_size" validate:"required"`
	Sort      int64  `json:"sort,optional"`
	Highlight int64  `json:"highlight,optional"`
}

type SearchVideoResp struct {
	Total  int64    `json:"total"`
	Videos []*Video `json:"videos"`
}

type CompleteVideoReq struct {
	Input string `json:"input" validate:"required"`
}

type CompleteVideoResp struct {
	Suggestions []string `json:"suggestions"`
}

type SearchUserReq struct {
	Keyword   string `json:"keyword" validate:"required"`
	Page      int64  `json:"page" validate:"required"`
	PageSize  int64  `json:"page_size" validate:"required"`
	Sort      int64  `json:"sort,optional"`
	Highlight int64  `json:"highlight,optional"`
}

type SearchUserResp struct {
	Total int64   `json:"total"`
	Users []*User `json:"users"`
}
