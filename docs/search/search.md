### 1. "根据关键字搜索视频"

1. route definition

- Url: /search/v1/video
- Method: GET
- Request: `SearchVideoReq`
- Response: `SearchVideoResp`

2. request definition



```golang
type SearchVideoReq struct {
	Keyword string `json:"keyword" validate:"required"`
	Page int64 `json:"page" validate:"required"`
	PageSize int64 `json:"page_size" validate:"required"`
	Sort int64 `json:"sort,optional"`
	Highlight int64 `json:"highlight,optional"`
}
```


3. response definition



```golang
type SearchVideoResp struct {
	Total int64 `json:"total"`
	Videos []*Video `json:"videos"`
}
```

### 2. "视频搜索输入自动补全"

1. route definition

- Url: /search/v1/video/suggestion
- Method: GET
- Request: `CompleteVideoReq`
- Response: `CompleteVideoResp`

2. request definition



```golang
type CompleteVideoReq struct {
	Input string `json:"input" validate:"required"`
}
```


3. response definition



```golang
type CompleteVideoResp struct {
	Suggestions []string `json:"suggestions"`
}
```

### 3. "根据关键字搜索用户"

1. route definition

- Url: /search/v1/user
- Method: GET
- Request: `SearchUserReq`
- Response: `SearchUserResp`

2. request definition



```golang
type SearchUserReq struct {
	Keyword string `json:"keyword" validate:"required"`
	Page int64 `json:"page" validate:"required"`
	PageSize int64 `json:"page_size" validate:"required"`
	Sort int64 `json:"sort,optional"`
	Highlight int64 `json:"highlight,optional"`
}
```


3. response definition



```golang
type SearchUserResp struct {
	Total int64 `json:"total"`
	Users []*User `json:"users"`
}
```

