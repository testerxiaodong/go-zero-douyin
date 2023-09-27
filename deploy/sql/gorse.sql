/*
 Navicat Premium Data Transfer

 Source Server         : 本地服务器
 Source Server Type    : MySQL
 Source Server Version : 80100
 Source Host           : 127.0.0.1:3306
 Source Schema         : gorse

 Target Server Type    : MySQL
 Target Server Version : 80100
 File Encoding         : 65001

 Date: 27/09/2023 21:34:50
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for feedback
-- ----------------------------
DROP TABLE IF EXISTS `feedback`;
CREATE TABLE `feedback` (
  `feedback_type` varchar(256) NOT NULL,
  `user_id` varchar(256) NOT NULL,
  `item_id` varchar(256) NOT NULL,
  `time_stamp` timestamp NOT NULL,
  `comment` text NOT NULL,
  PRIMARY KEY (`feedback_type`,`user_id`,`item_id`),
  KEY `user_id` (`user_id`),
  KEY `item_id` (`item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for items
-- ----------------------------
DROP TABLE IF EXISTS `items`;
CREATE TABLE `items` (
  `item_id` varchar(256) NOT NULL,
  `time_stamp` timestamp NOT NULL,
  `labels` json NOT NULL,
  `comment` text NOT NULL,
  `is_hidden` tinyint(1) NOT NULL,
  `categories` json NOT NULL,
  PRIMARY KEY (`item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `user_id` varchar(256) NOT NULL,
  `labels` json NOT NULL,
  `comment` text NOT NULL,
  `subscribe` json NOT NULL,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;
