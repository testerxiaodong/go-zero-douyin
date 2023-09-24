### 1. "分区列表"

1. route definition

- Url: /video/v1/section/list
- Method: GET
- Request: `-`
- Response: `GetAllSectionResp`

2. request definition



3. response definition



```golang
type GetAllSectionResp struct {
	Sections []*Section `json:"sections"`
}
```

### 2. "标签列表"

1. route definition

- Url: /video/v1/tag/list
- Method: GET
- Request: `-`
- Response: `GetAllTagResp`

2. request definition



3. response definition



```golang
type GetAllTagResp struct {
	Tags []*Tag `json:"tags"`
}
```

### 3. "视频详情"

1. route definition

- Url: /video/v1/detail
- Method: GET
- Request: `VideoDetailReq`
- Response: `VideoDetailResp`

2. request definition



```golang
type VideoDetailReq struct {
	VideoId int64 `json:"video_id"`
}
```


3. response definition



```golang
type VideoDetailResp struct {
	Video *VideoInfo `json:"video"`
}
```

### 4. "获取视频流"

1. route definition

- Url: /video/v1/feed
- Method: GET
- Request: `VideoFeedReq`
- Response: `VideoFeedResp`

2. request definition



```golang
type VideoFeedReq struct {
	LastTimeStamp int64 `json:"last_time_stamp" validate:"required"`
	SectionId int64 `json:"section_id" validate:"required"`
}
```


3. response definition



```golang
type VideoFeedResp struct {
	Videos []*VideoInfo `json:"videos"`
}
```

### 5. "获取用户发布视频列表"

1. route definition

- Url: /video/v1/list
- Method: GET
- Request: `UserVideoListReq`
- Response: `UserVideoListResp`

2. request definition



```golang
type UserVideoListReq struct {
	UserId int64 `json:"user_id" validate:"required"`
	Page int64 `json:"page,optional"`
	PageSize int64 `json:"page_size,optional"`
}
```


3. response definition



```golang
type UserVideoListResp struct {
	Total int64 `json:"total"`
	Videos []*VideoInfo `json:"videos"`
}
```

### 6. "新增分区"

1. route definition

- Url: /video/v1/section/add
- Method: POST
- Request: `AddSectionReq`
- Response: `-`

2. request definition



```golang
type AddSectionReq struct {
	Name string `json:"name"`
}
```


3. response definition


### 7. "删除分区"

1. route definition

- Url: /video/v1/section/del
- Method: POST
- Request: `DelSectionReq`
- Response: `-`

2. request definition



```golang
type DelSectionReq struct {
	Id int64 `json:"id"`
}
```


3. response definition


### 8. "新增标签"

1. route definition

- Url: /video/v1/tag/add
- Method: POST
- Request: `AddTagReq`
- Response: `-`

2. request definition



```golang
type AddTagReq struct {
	Name string `json:"name"`
}
```


3. response definition


### 9. "删除标签"

1. route definition

- Url: /video/v1/tag/del
- Method: POST
- Request: `DelTagReq`
- Response: `-`

2. request definition



```golang
type DelTagReq struct {
	Id int64 `json:"id"`
}
```


3. response definition


### 10. "发布视频"

1. route definition

- Url: /video/v1/publish
- Method: POST
- Request: `PublishVideoReq`
- Response: `-`

2. request definition



```golang
type PublishVideoReq struct {
	Title string `form:"title" validate:"required"`
	SectionId int64 `form:"section_id" validate:"required"`
	TagIds string `form:"tag_ids" validate:"required"`
	PublishTime int64 `form:"publish_time,optional"`
}
```


3. response definition


### 11. "删除视频"

1. route definition

- Url: /video/v1/delete
- Method: POST
- Request: `DeleteVideoReq`
- Response: `-`

2. request definition



```golang
type DeleteVideoReq struct {
	VideoId int64 `json:"video_id" validate:"required"`
}
```


3. response definition


