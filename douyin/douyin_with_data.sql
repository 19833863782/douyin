/*
Navicat MySQL Data Transfer

Source Server         : test
Source Server Version : 50015
Source Host           : localhost:3306
Source Database       : douyin

Target Server Type    : MYSQL
Target Server Version : 50015
File Encoding         : 65001

Date: 2023-02-23 14:23:53
*/

SET FOREIGN_KEY_CHECKS=0;
-- ----------------------------
-- Table structure for `dy_comment`
-- ----------------------------
DROP TABLE IF EXISTS `dy_comment`;
CREATE TABLE `dy_comment` (
  `id` int(11) NOT NULL auto_increment,
  `sender_id` int(11) NOT NULL,
  `video_id` int(11) NOT NULL,
  `content` text NOT NULL,
  `status` int(11) NOT NULL,
  `create_time` timestamp NOT NULL default CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL default '0000-00-00 00:00:00',
  PRIMARY KEY  (`id`),
  KEY `user_id` (`sender_id`),
  KEY `video_id` (`video_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of dy_comment
-- ----------------------------
INSERT INTO `dy_comment` VALUES ('1', '1', '1', '岳飞战小儿', '1', '2023-02-20 19:48:57', '2022-02-20 13:00:00');
INSERT INTO `dy_comment` VALUES ('2', '2', '1', '实验', '78', '2023-02-20 20:59:33', '0000-00-00 00:00:00');
INSERT INTO `dy_comment` VALUES ('3', '1', '1', '评论', '78', '2023-02-20 21:05:17', '0000-00-00 00:00:00');
INSERT INTO `dy_comment` VALUES ('4', '1', '1', '11111', '78', '2023-02-20 21:06:12', '0000-00-00 00:00:00');

-- ----------------------------
-- Table structure for `dy_favorite`
-- ----------------------------
DROP TABLE IF EXISTS `dy_favorite`;
CREATE TABLE `dy_favorite` (
  `video_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `create_time` timestamp NOT NULL default CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL default '0000-00-00 00:00:00',
  PRIMARY KEY  (`video_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of dy_favorite
-- ----------------------------
INSERT INTO `dy_favorite` VALUES ('1', '1', '2023-02-20 20:55:56', '0000-00-00 00:00:00');
INSERT INTO `dy_favorite` VALUES ('1', '2', '2023-02-20 21:08:59', '0000-00-00 00:00:00');

-- ----------------------------
-- Table structure for `dy_feed`
-- ----------------------------
DROP TABLE IF EXISTS `dy_feed`;
CREATE TABLE `dy_feed` (
  `user_id` int(11) NOT NULL,
  `video_id` int(11) NOT NULL,
  `create_time` timestamp NOT NULL default CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL default '0000-00-00 00:00:00',
  PRIMARY KEY  (`user_id`,`video_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of dy_feed
-- ----------------------------
INSERT INTO `dy_feed` VALUES ('1', '1', '2023-02-20 19:49:05', '2022-02-20 14:00:00');

-- ----------------------------
-- Table structure for `dy_follow`
-- ----------------------------
DROP TABLE IF EXISTS `dy_follow`;
CREATE TABLE `dy_follow` (
  `follow_id` int(11) NOT NULL,
  `follower_id` int(11) NOT NULL,
  `create_time` timestamp NOT NULL default CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL default '0000-00-00 00:00:00',
  PRIMARY KEY  (`follow_id`,`follower_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of dy_follow
-- ----------------------------
INSERT INTO `dy_follow` VALUES ('1', '1', '2023-02-20 20:56:21', '0000-00-00 00:00:00');
INSERT INTO `dy_follow` VALUES ('1', '2', '2023-02-20 20:57:24', '0000-00-00 00:00:00');
INSERT INTO `dy_follow` VALUES ('2', '1', '2023-02-21 20:57:07', '0000-00-00 00:00:00');

-- ----------------------------
-- Table structure for `dy_user`
-- ----------------------------
DROP TABLE IF EXISTS `dy_user`;
CREATE TABLE `dy_user` (
  `id` int(11) NOT NULL auto_increment,
  `name` varchar(20) NOT NULL,
  `follow_count` int(11) NOT NULL,
  `follower_count` int(11) NOT NULL,
  `password` varchar(32) NOT NULL,
  `create_time` timestamp NOT NULL default CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL default '0000-00-00 00:00:00',
  PRIMARY KEY  (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of dy_user
-- ----------------------------
INSERT INTO `dy_user` VALUES ('1', '2933142084@qq.com', '6', '12', 'd51839bf6568ae88db6d9a20c2ffb3e9', '2023-02-20 19:38:27', '2022-02-20 19:50:00');
INSERT INTO `dy_user` VALUES ('2', '123', '2', '1', 'e10adc3949ba59abbe56e057f20f883e', '2023-02-20 20:57:07', '0000-00-00 00:00:00');

-- ----------------------------
-- Table structure for `dy_video`
-- ----------------------------
DROP TABLE IF EXISTS `dy_video`;
CREATE TABLE `dy_video` (
  `id` int(11) NOT NULL auto_increment,
  `title` varchar(50) NOT NULL,
  `play_url` varchar(100) NOT NULL,
  `cover_url` varchar(100) NOT NULL,
  `favorite_count` int(11) NOT NULL,
  `comment_count` int(11) NOT NULL,
  `author_id` int(11) NOT NULL,
  `create_time` timestamp NOT NULL default CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL default '0000-00-00 00:00:00',
  PRIMARY KEY  (`id`),
  KEY `user_id` (`author_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of dy_video
-- ----------------------------
INSERT INTO `dy_video` VALUES ('1', '满江红', 'https://www.baidu.com', 'https://www.baidu.com', '13', '13', '1', '2023-02-20 19:48:58', '2023-02-20 19:44:00');
