SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

/*
 *
 *  * Copyright 2022 CloudWeGo Authors
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *  * You may obtain a copy of the License at
 *  *
 *  *     http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  * Unless required by applicable law or agreed to in writing, software
 *  * distributed under the License is distributed on an "AS IS" BASIS,
 *  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  * See the License for the specific language governing permissions and
 *  * limitations under the License.
 *
 */

-- ----------------------------
-- Table structure for repository
-- ----------------------------
DROP TABLE IF EXISTS `repository`;
CREATE TABLE `repository` (
  `id` bigint(20) NOT NULL COMMENT 'repository id',
  `repository_type` tinyint(4) NOT NULL COMMENT 'repository type',
  `store_type` tinyint(4) NOT NULL COMMENT 'store type',
  `repository_url` varchar(1024) NOT NULL COMMENT 'repository URL',
  `last_update_time` datetime DEFAULT NULL COMMENT 'last update time',
  `last_sync_time` datetime DEFAULT NULL COMMENT 'last sync time',
  `token` varchar(1024) DEFAULT '' COMMENT 'repository token',
  `status` varchar(20) DEFAULT 'active' COMMENT 'status',
  `is_deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'is deleted',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='repository table';

SET FOREIGN_KEY_CHECKS = 1;
