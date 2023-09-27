### 1. "用户视频分区个性化推荐"

1. route definition

- Url: /recommend/v1/video
- Method: GET
- Request: `VideoRecommendReq`
- Response: `VideoRecommendResp`

2. request definition



```golang
type VideoRecommendReq struct {
	SectionId int64 `json:"section_id" validate:"required"`
	Count int64 `json:"count" validate:"required"`
}
```


3. response definition



```golang
type VideoRecommendResp struct {
	Videos []*Video `json:"videos"`
}
```

