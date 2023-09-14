// Code generated by goctl. DO NOT EDIT.
package types

type UserInfo struct {
	Id            int64  `json:"id"`
	Username      string `json:"username"`
	FollowerCount int64  `json:"follower_count"`
	FollowCount   int64  `json:"follow_count"`
}

type RegisterReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,password"`
}

type RegisterResp struct {
	AccessToken  string `json:"access_token"`
	RefreshAfter int64  `json:"refresh_after"`
	ExpireTime   int64  `json:"expire_time"`
}

type LoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshAfter int64  `json:"refresh_after"`
	ExpireTime   int64  `json:"expire_time"`
}

type UserInfoReq struct {
	Id int64 `json:"id" validate:"required"`
}

type UserInfoResp struct {
	User UserInfo `json:"user_info"`
}

type UpdateUserReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,password"`
}

type SyncUserToEsByIdReq struct {
	UserId int64 `json:"user_id" validate:"required"`
}
