/*
 *
 * Copyright 2023 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
CREATE DATABASE IF NOT EXISTS `cwgo` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

USE `cwgo`;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for template_item
-- ----------------------------
DROP TABLE IF EXISTS `template_item`;
CREATE TABLE `template_item` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'template item id',
  `template_id` bigint(20) NOT NULL COMMENT 'template id',
  `name` varchar(255) NOT NULL COMMENT 'template item name',
  `content` text COMMENT 'template content',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'is deleted',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='template item table';

SET FOREIGN_KEY_CHECKS = 1;
