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

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for idl
-- ----------------------------
DROP TABLE IF EXISTS `idl`;
CREATE TABLE `idl`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `idl_repository_id` bigint(20) NOT NULL COMMENT 'idl_repository id',
  `service_repository_id` bigint(20) NOT NULL COMMENT 'serviec repository id',
  `parent_idl_id` bigint(20) NULL DEFAULT NULL COMMENT 'null if main idl else import idl',
  `idl_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'idl path',
  `commit_hash` char(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'idl file commit hash',
  `service_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'service name',
  `template_id` bigint(20) NOT NULL DEFAULT 0 COMMENT 'template id',
  `status` tinyint(2) NOT NULL DEFAULT 1 COMMENT 'status',
  `last_sync_time` datetime NOT NULL COMMENT 'last update time',
  `is_deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'is deleted',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `parent_idl_id`(`parent_idl_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'IDL table' ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;

SET FOREIGN_KEY_CHECKS = 1;
