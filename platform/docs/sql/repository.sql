SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for repository
-- ----------------------------
DROP TABLE IF EXISTS `repository`;
CREATE TABLE `repository` (
  `id` varchar(36) NOT NULL COMMENT '仓库ID',
  `repository_url` varchar(255) NOT NULL COMMENT '仓库URL',
  `last_update_time` datetime DEFAULT NULL COMMENT '最后更新时间',
  `last_sync_time` datetime DEFAULT NULL COMMENT '最后同步时间',
  `token` varchar(255) DEFAULT '' COMMENT '令牌',
  `status` varchar(20) DEFAULT 'active' COMMENT '状态',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `repo_type` tinyint(4) DEFAULT '0' COMMENT '仓库类型',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
