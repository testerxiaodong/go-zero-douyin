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

CREATE TABLE `user_gorse_sink` (
  `user_id` STRING,
  `labels` STRING,
  `comment` STRING,
  `subscribe` STRING,
  PRIMARY KEY (`user_id`) NOT ENFORCED
  ) WITH (
'connector' = 'jdbc',
'url' = 'jdbc:mysql://192.168.2.248:3306/gorse?serverTimezone=Asia/Shanghai',
'username' = 'root',
'password' = 'my-secret-pw',
'driver' = 'com.mysql.cj.jdbc.Driver',
'table-name' = 'users'
);

INSERT INTO user_gorse_sink (user_id, labels, `comment`, subscribe)
SELECT
    CAST(u.id AS STRING) AS user_id,
    '[]' AS labels,
    '' AS `comment`,
    '[]' AS subscribe
FROM user_source u;

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

CREATE TABLE `video_gorse_sink` (
  `item_id` STRING,
  `is_hidden` tinyint,
  `categories` STRING,
  `time_stamp` timestamp,
  `labels` STRING,
  `comment` STRING,
  PRIMARY KEY (`item_id`) NOT ENFORCED
  ) WITH (
'connector' = 'jdbc',
'url' = 'jdbc:mysql://192.168.2.248:3306/gorse?serverTimezone=Asia/Shanghai',
'username' = 'root',
'password' = 'my-secret-pw',
'driver' = 'com.mysql.cj.jdbc.Driver',
'table-name' = 'items'
);

INSERT INTO video_gorse_sink (item_id, is_hidden, categories, time_stamp, labels, `comment`)
SELECT
    CAST(v.id AS STRING) AS item_id,
    CAST(0 AS TINYINT) AS is_hidden,
    CAST(JSON_ARRAY(CAST(v.section_id AS STRING)) AS STRING) AS categories,
    v.create_time AS time_stamp,
    v.tag_ids AS labels,
    '' AS `comment`
FROM video_source v;

CREATE TABLE `like_source` (
  `id` bigint,
  `video_id` bigint,
  `user_id` bigint,
  `status` tinyint,
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
'table-name' = 'like'
);

CREATE TABLE `like_gorse_sink` (
  `feedback_type` STRING,
  `user_id` STRING,
  `item_id` STRING,
  `time_stamp` timestamp,
  `comment` STRING,
  PRIMARY KEY (`feedback_type`,`user_id`,`item_id`) NOT ENFORCED
  ) WITH (
'connector' = 'jdbc',
'url' = 'jdbc:mysql://192.168.2.248:3306/gorse?serverTimezone=Asia/Shanghai',
'username' = 'root',
'password' = 'my-secret-pw',
'driver' = 'com.mysql.cj.jdbc.Driver',
'table-name' = 'feedback'
);

INSERT INTO like_gorse_sink(feedback_type, user_id, item_id, time_stamp, `comment`)
SELECT
    'like' AS feedback_type,
    CAST(l.user_id AS STRING) user_id,
    CAST(l.video_id AS STRING) item_id,
    l.update_time AS time_stamp,
    '' AS `comment`
FROM like_source l WHERE l.status = 1;