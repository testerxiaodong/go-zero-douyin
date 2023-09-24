### 1. "用户注册接口"

1. route definition

- Url: /user/v1/register
- Method: POST
- Request: `RegisterReq`
- Response: `RegisterResp`

2. request definition



```golang
type RegisterReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,password"`
}
```


3. response definition



```golang
type RegisterResp struct {
	AccessToken string `json:"access_token"`
	RefreshAfter int64 `json:"refresh_after"`
	ExpireTime int64 `json:"expire_time"`
}
```

### 2. "用户登陆接口"

1. route definition

- Url: /user/v1/login
- Method: POST
- Request: `LoginReq`
- Response: `LoginResp`

2. request definition



```golang
type LoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
```


3. response definition



```golang
type LoginResp struct {
	AccessToken string `json:"access_token"`
	RefreshAfter int64 `json:"refresh_after"`
	ExpireTime int64 `json:"expire_time"`
}
```

### 3. "获取用户信息"

1. route definition

- Url: /user/v1/detail
- Method: GET
- Request: `UserInfoReq`
- Response: `UserInfoResp`

2. request definition



```golang
type UserInfoReq struct {
	Id int64 `json:"id" validate:"required"`
}
```


3. response definition



```golang
type UserInfoResp struct {
	User UserInfo `json:"user_info"`
}

type UserInfo struct {
	Id int64 `json:"id"`
	Username string `json:"username"`
	FollowerCount int64 `json:"follower_count"`
	FollowCount int64 `json:"follow_count"`
}
```

### 4. "更新用户信息"

1. route definition

- Url: /user/v1/update
- Method: POST
- Request: `UpdateUserReq`
- Response: `-`

2. request definition



```golang
type UpdateUserReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,password"`
}
```


3. response definition


