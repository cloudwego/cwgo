SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for idl
-- ----------------------------
DROP TABLE IF EXISTS `idl`;
CREATE TABLE `idl` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `repository_id` bigint(20) NOT NULL COMMENT '仓库ID',
  `main_idl_path` varchar(255) NOT NULL COMMENT '主要IDL路径',
  `content` text COMMENT '内容',
  `service_name` varchar(255) NOT NULL COMMENT '服务名称',
  `last_sync_time` varchar(255) NOT NULL COMMENT '最后同步时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='IDL表';

SET FOREIGN_KEY_CHECKS = 1;
