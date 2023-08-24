SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for templateItem
-- ----------------------------
DROP TABLE IF EXISTS `templateItem`;
CREATE TABLE `templateItem` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `template_id` bigint(20) NOT NULL COMMENT '模板ID',
  `name` varchar(255) NOT NULL COMMENT '名称',
  `content` text NOT NULL COMMENT '内容',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='模板项表';

SET FOREIGN_KEY_CHECKS = 1;
