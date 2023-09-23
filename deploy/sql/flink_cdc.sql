SET execution.checkpointing.interval = 3s;

CREATE TABLE `user_source` (
  `id` bigint,
  `username` STRING,
  `password` STRING,
  `create_time` timestamp,
  `update_time` timestamp,
  `delete_time` timestamp,
  `is_delete` tinyint,
  `version` bigint,
  PRIMARY KEY (`id`) NOT ENFORCED
  ) WITH (
'connector' = 'mysql-cdc',
'hostname' = '192.168.2.248',
'port' = '3306',
'username' = 'root',
'password' = 'my-secret-pw',
'database-name' = 'go_zero_douyin',
'table-name' = 'user'
);

CREATE TABLE `follow_count_source` (
  `id` bigint,
  `user_id` bigint,
  `follow_count` int,
  `follower_count` int,
  `create_time` timestamp,
  `update_time` timestamp,
  `delete_time` timestamp,
  `is_delete` tinyint,
  `version` bigint,
  PRIMARY KEY (`id`) NOT ENFORCED
  ) WITH (
'connector' = 'mysql-cdc',
'hostname' = '192.168.2.248',
'port' = '3306',
'username' = 'root',
'password' = 'my-secret-pw',
'database-name' = 'go_zero_douyin',
'table-name' = 'follow_count'
);

CREATE TABLE `user_es_sink` (
  `id` bigint,
  `username` STRING,
  `follower_count` int,
  `follow_count` int, 
  `suggestion` STRING,
  `create_time` timestamp,
  `update_time` timestamp,
  `version` bigint,
  PRIMARY KEY (`id`)  NOT ENFORCED
 ) WITH (
'connector' = 'elasticsearch-7',
'hosts' = 'http://192.168.2.248:9200',
'index' = 'user'
);

INSERT INTO user_es_sink
SELECT
    u.id, u.username,
    IFNULL(fc.follower_count, 0) AS follower_count,
    IFNULL(fc.follow_count, 0) AS follow_count,
    u.username AS suggestion,
    u.create_time,
    u.update_time,
    u.version
FROM user_source AS u
LEFT JOIN follow_count_source AS fc ON u.id = fc.user_id
WHERE u.is_delete = 0;

CREATE TABLE `video_source` (
  `id` bigint,
  `title` STRING,
  `section_id` bigint,
  `tag_ids` STRING,
  `owner_id` int,
  `owner_name` STRING,
  `play_url` STRING,
  `cover_url` STRING,
  `create_time` timestamp,
  `update_time` timestamp,
  `delete_time` timestamp,
  `is_delete` tinyint,
  `version` bigint,
  PRIMARY KEY (`id`) NOT ENFORCED
   ) WITH (
'connector' = 'mysql-cdc',
'hostname' = '192.168.2.248',
'port' = '3306',
'username' = 'root',
'password' = 'my-secret-pw',
'database-name' = 'go_zero_douyin',
'table-name' = 'video'
);

CREATE TABLE `comment_count_source` (
  `id` bigint,
  `video_id` bigint,
  `comment_count` int,
  `create_time` timestamp,
  `update_time` timestamp,
  `delete_time` timestamp,
  `is_delete` tinyint,
  `version` bigint,
  PRIMARY KEY (`id`) NOT ENFORCED
 ) WITH (
'connector' = 'mysql-cdc',
'hostname' = '192.168.2.248',
'port' = '3306',
'username' = 'root',
'password' = 'my-secret-pw',
'database-name' = 'go_zero_douyin',
'table-name' = 'comment_count'
);

CREATE TABLE `like_count_source` (
  `id` bigint,
  `video_id` bigint,
  `like_count` int,
  `create_time` timestamp,
  `update_time` timestamp,
  `delete_time` timestamp,
  `is_delete` tinyint,
  `version` bigint,
  PRIMARY KEY (`id`) NOT ENFORCED
 ) WITH (
'connector' = 'mysql-cdc',
'hostname' = '192.168.2.248',
'port' = '3306',
'username' = 'root',
'password' = 'my-secret-pw',
'database-name' = 'go_zero_douyin',
'table-name' = 'like_count'
);

CREATE TABLE `video_es_sink` (
  `id` bigint,
  `title` STRING,
  `section_id` bigint ,
  `tag_ids` STRING,
  `owner_id` int,
  `owner_name` STRING,
  `play_url` STRING,
  `cover_url` STRING,
  `comment_count` int,
  `like_count` int,
  `suggestion` STRING,
  `create_time` timestamp,
  `update_time` timestamp,
  `version` bigint,
  PRIMARY KEY (`id`) NOT ENFORCED
) WITH (
'connector' = 'elasticsearch-7',
'hosts' = 'http://192.168.2.248:9200',
'index' = 'video'
);

INSERT INTO video_es_sink
SELECT
    v.id,
    v.title,
    v.section_id,
    v.tag_ids,
    v.owner_id,
    v.owner_name,
    v.play_url,
    v.cover_url,
    IFNULL(cc.comment_count, 0) AS comment_count,
    IFNULL(lc.like_count, 0) AS like_count,
    v.title AS suggestion,
    v.create_time,
    v.update_time,
    v.version FROM video_source v
LEFT JOIN comment_count_source cc ON v.id = cc.video_id
LEFT JOIN like_count_source lc ON v.id = lc.video_id WHERE v.is_delete = 0;
