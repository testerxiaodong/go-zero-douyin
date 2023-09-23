/*
 Navicat Premium Data Transfer

 Source Server         : 本地服务器
 Source Server Type    : MySQL
 Source Server Version : 80100
 Source Host           : 127.0.0.1:3306
 Source Schema         : go_zero_douyin

 Target Server Type    : MySQL
 Target Server Version : 80100
 File Encoding         : 65001

 Date: 24/09/2023 00:41:34
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for comment
-- ----------------------------
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '评论id',
  `video_id` bigint unsigned NOT NULL COMMENT '评论所属视频id',
  `user_id` bigint unsigned NOT NULL COMMENT '发表评论的用户id',
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '评论内容',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `delete_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '删除时间',
  `is_delete` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否被删除',
  `version` bigint unsigned NOT NULL DEFAULT '0' COMMENT '版本号',
  PRIMARY KEY (`id`),
  KEY `idx_videoid_userid` (`video_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='视频评论记录';

-- ----------------------------
-- Table structure for comment_count
-- ----------------------------
DROP TABLE IF EXISTS `comment_count`;
CREATE TABLE `comment_count` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '评论数id',
  `video_id` bigint unsigned NOT NULL COMMENT '视频id',
  `comment_count` int unsigned NOT NULL DEFAULT '0' COMMENT '评论数',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `delete_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '删除时间',
  `is_delete` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否被删除',
  `version` bigint unsigned NOT NULL DEFAULT '0' COMMENT '版本号',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_videoid` (`video_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='视频评论数记录';

-- ----------------------------
-- Table structure for follow
-- ----------------------------
DROP TABLE IF EXISTS `follow`;
CREATE TABLE `follow` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '关注id',
  `user_id` bigint unsigned NOT NULL COMMENT '被关注者id',
  `follower_id` bigint unsigned NOT NULL COMMENT '关注者id',
  `status` tinyint unsigned NOT NULL COMMENT '关注状态：0取消关注，1已关注',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `delete_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '删除时间',
  `is_delete` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否被删除',
  `version` bigint unsigned NOT NULL DEFAULT '0' COMMENT '版本号',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_userid_followerid` (`user_id`,`follower_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户关注/被关注记录表';

-- ----------------------------
-- Table structure for follow_count
-- ----------------------------
DROP TABLE IF EXISTS `follow_count`;
CREATE TABLE `follow_count` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '关注数id',
  `user_id` bigint NOT NULL COMMENT '用户id',
  `follow_count` int unsigned NOT NULL DEFAULT '0' COMMENT '关注数',
  `follower_count` int unsigned NOT NULL DEFAULT '0' COMMENT '粉丝数',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `delete_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '删除时间',
  `is_delete` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否被删除',
  `version` bigint unsigned NOT NULL DEFAULT '0' COMMENT '版本号',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_userid` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户关注/粉丝数记录';

-- ----------------------------
-- Table structure for like
-- ----------------------------
DROP TABLE IF EXISTS `like`;
CREATE TABLE `like` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '点赞记录id',
  `video_id` bigint unsigned NOT NULL COMMENT '视频id',
  `user_id` bigint unsigned NOT NULL COMMENT '用户id',
  `status` tinyint unsigned NOT NULL COMMENT '点赞状态：0取消点赞，1已点赞',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `delete_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '删除时间',
  `is_delete` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否被删除',
  `version` bigint unsigned NOT NULL DEFAULT '0' COMMENT '版本号',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_videoid_userid` (`video_id`,`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='视频点赞记录';

-- ----------------------------
-- Table structure for like_count
-- ----------------------------
DROP TABLE IF EXISTS `like_count`;
CREATE TABLE `like_count` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '点赞数id',
  `video_id` bigint unsigned NOT NULL COMMENT '视频id',
  `like_count` int unsigned NOT NULL COMMENT '点赞数',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `delete_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '删除时间',
  `is_delete` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否被删除',
  `version` bigint unsigned NOT NULL DEFAULT '0' COMMENT '版本号',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_videoid` (`video_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='视频点赞数记录';

-- ----------------------------
-- Table structure for section
-- ----------------------------
DROP TABLE IF EXISTS `section`;
CREATE TABLE `section` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '分区id',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '分区名',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `delete_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '删除时间',
  `is_delete` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否被删除',
  `version` bigint unsigned NOT NULL DEFAULT '0' COMMENT '版本号',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`,`is_delete`,`delete_time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='视频分区记录';

-- ----------------------------
-- Table structure for tag
-- ----------------------------
DROP TABLE IF EXISTS `tag`;
CREATE TABLE `tag` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '标签id',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '标签名',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `delete_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '删除时间',
  `is_delete` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否被删除',
  `version` bigint unsigned NOT NULL DEFAULT '0' COMMENT '版本号',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`,`is_delete`,`delete_time`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='视频标签记录';

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户id',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户名',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `delete_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '删除时间',
  `is_delete` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否被删除',
  `version` bigint unsigned NOT NULL DEFAULT '0' COMMENT '版本号',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`,`is_delete`,`delete_time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户记录';

-- ----------------------------
-- Table structure for video
-- ----------------------------
DROP TABLE IF EXISTS `video`;
CREATE TABLE `video` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '视频id',
  `title` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '视频标题',
  `section_id` bigint NOT NULL COMMENT '视频分区id',
  `tag_ids` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '视频标签列表',
  `owner_id` int unsigned NOT NULL COMMENT '视频发布者id',
  `owner_name` varchar(255) NOT NULL COMMENT '视频发布者用户名',
  `play_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '视频下载地址',
  `cover_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '封面下载地址',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `delete_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '删除时间',
  `is_delete` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否被删除',
  `version` bigint unsigned NOT NULL DEFAULT '0' COMMENT '版本号',
  PRIMARY KEY (`id`),
  KEY `idx_ownerid` (`owner_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='视频记录';

SET FOREIGN_KEY_CHECKS = 1;
