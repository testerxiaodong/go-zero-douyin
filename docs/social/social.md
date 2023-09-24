### 1. "获取视频评论列表"

1. route definition

- Url: /social/v1/comment/list
- Method: GET
- Request: `GetVideoCommentListReq`
- Response: `GetVideoCommentListResp`

2. request definition



```golang
type GetVideoCommentListReq struct {
	VideoId int64 `json:"video_id" validate:"required"`
	Page int64 `json:"page,optional"`
	PageSize int64 `json:"page_size,optional"`
}
```


3. response definition



```golang
type GetVideoCommentListResp struct {
	Comments []*Comment `json:"comments"`
}
```

### 2. "添加评论"

1. route definition

- Url: /social/v1/comment/add
- Method: POST
- Request: `AddCommentReq`
- Response: `-`

2. request definition



```golang
type AddCommentReq struct {
	VideoId int64 `json:"video_id" validate:"required"`
	Content string `json:"content" validate:"required"`
}
```


3. response definition


### 3. "删除评论"

1. route definition

- Url: /social/v1/comment/del
- Method: POST
- Request: `DelCommentReq`
- Response: `-`

2. request definition



```golang
type DelCommentReq struct {
	CommentId int64 `json:"comment_id" validate:"required"`
}
```


3. response definition


### 4. "获取用户点赞视频id列表"

1. route definition

- Url: /social/v1/like/list
- Method: GET
- Request: `GetUserLikeVideoIdListReq`
- Response: `GetUserLikeVideoIdListResp`

2. request definition



```golang
type GetUserLikeVideoIdListReq struct {
	UserId int64 `json:"user_id" validate:"required"`
	Page int64 `json:"page,optional"`
	PageSize int64 `json:"page_size,optional"`
}
```


3. response definition



```golang
type GetUserLikeVideoIdListResp struct {
	VideoIdList []int64 `json:"id_list"`
}
```

### 5. "视频点赞"

1. route definition

- Url: /social/v1/like/add
- Method: POST
- Request: `VideoLikeReq`
- Response: `-`

2. request definition



```golang
type VideoLikeReq struct {
	VideoId int64 `json:"video_id" validate:"required"`
}
```


3. response definition


### 6. "视频取消点赞"

1. route definition

- Url: /social/v1/like/del
- Method: POST
- Request: `VideoUnlikeReq`
- Response: `-`

2. request definition



```golang
type VideoUnlikeReq struct {
	VideoId int64 `json:"video_id" validate:"required"`
}
```


3. response definition


### 7. "获取用户关注数"

1. route definition

- Url: /social/v1/follow/follow_count
- Method: GET
- Request: `GetUserFollowCountReq`
- Response: `GetUserFollowCountResp`

2. request definition



```golang
type GetUserFollowCountReq struct {
	UserId int64 `json:"user_id" validate:"required"`
}
```


3. response definition



```golang
type GetUserFollowCountResp struct {
	FollowCount int64 `json:"follow_count"`
}
```

### 8. "获取用户关注id列表"

1. route definition

- Url: /social/v1/follow/follow_list
- Method: GET
- Request: `GetUserFollowIdListReq`
- Response: `GetUserFollowIdListResp`

2. request definition



```golang
type GetUserFollowIdListReq struct {
	UserId int64 `json:"user_id" validate:"required"`
	Page int64 `json:"page,optional"`
	PageSize int64 `json:"page_size,optional"`
}
```


3. response definition



```golang
type GetUserFollowIdListResp struct {
	UserIdList []int64 `json:"user_id_list"`
}
```

### 9. "获取用户粉丝数"

1. route definition

- Url: /social/v1/follow/follower_count
- Method: GET
- Request: `GetUserFollowerCountReq`
- Response: `GetUserFollowerCountResp`

2. request definition



```golang
type GetUserFollowerCountReq struct {
	UserId int64 `json:"user_id" validate:"required"`
}
```


3. response definition



```golang
type GetUserFollowerCountResp struct {
	FollowerCount int64 `json:"follower_count"`
}
```

### 10. "获取用户粉丝id列表"

1. route definition

- Url: /social/v1/follow/follower_list
- Method: GET
- Request: `GetUserFollowerIdListReq`
- Response: `GetUserFollowerIdListResp`

2. request definition



```golang
type GetUserFollowerIdListReq struct {
	UserId int64 `json:"user_id" validate:"required"`
	Page int64 `json:"page,optional"`
	PageSize int64 `json:"page_size,optional"`
}
```


3. response definition



```golang
type GetUserFollowerIdListResp struct {
	UserIdList []int64 `json:"user_id_list"`
}
```

### 11. "用户关注"

1. route definition

- Url: /social/v1/follow/add
- Method: POST
- Request: `UserFollowReq`
- Response: `-`

2. request definition



```golang
type UserFollowReq struct {
	UserId int64 `json:"user_id" validate:"required"`
}
```


3. response definition


### 12. "用户取消关注"

1. route definition

- Url: /social/v1/follow/del
- Method: POST
- Request: `UserUnfollowReq`
- Response: `-`

2. request definition



```golang
type UserUnfollowReq struct {
	UserId int64 `json:"user_id" validate:"required"`
}
```


3. response definition


