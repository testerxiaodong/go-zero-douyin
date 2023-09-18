## go-zero-douyin
### 项目介绍
faker-douyin的go-zero版本，主要是想学习一下go-zero微服务框架。

因为go-zero深度绑定了orm框架sqlx以及sqlc并帮助了用户处理了缓存击穿问题，与我在faker-douyin中通过rabbitmq处理缓存一致性问题时不一致

准备先继续采用faker-douyin的处理方式，沿用gorm框架(因为代码生成的方式熟悉一点，而且不用手撸sql)，同时也参考了一些go-zero对于缓存的处理

### 项目预计用到技术
- go-zero
- gorm/gen
- mysql
- redis
- asynq
- elasticsearch
- go-stash
- kafka
- filebeat
- jaeger
- prometheus
- grafana
- rabbitmq

### 项目功能点
还是faker-douyin的老功能，不过这次我准备把视频数据上传到阿里云的oss服务，并且限制上传视频文件的大小

项目日志的记录直接使用go-zero绑定的logx，且集成了elk日志系统

go-zero绑定的消息代理是kafka，不太熟悉，先继续用rabbitmq，做完之后再替换为kafka

使用asynq作为分布式任务队列，实现视频的延迟发布（延迟任务）

项目功能点：
- 用户
    - 注册
    - 登陆（先做简单用户名密码登陆，之后换成第三方登陆）
    - 获取用户信息
    - 更新用户信息
- 视频
    - 获取视频流
    - 发布视频(延迟发布)
    - 查询用户视频列表
- 评论
    - 添加评论
    - 删除评论
- 点赞
    - 点赞视频
    - 取消点赞
    - 获取用户点赞视频列表
    - 获取视频点赞数
- 关注
    - 关注用户
    - 取消关注
    - 获取用户关注列表
    - 获取用户被关注数
- 搜索
    - 搜索用户
    - 搜索视频
