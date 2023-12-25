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
-- Table structure for repository
-- ----------------------------
DROP TABLE IF EXISTS `repository`;
CREATE TABLE `repository`  (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'repository id',
    `repository_type` tinyint(4) NOT NULL COMMENT 'repo type',
    `domain` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'repo domain',
    `owner` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'repo owner',
    `repository_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'repo name',
    `branch` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'repo branch',
    `store_type` tinyint(4) NOT NULL COMMENT 'store type',
    `last_update_time` datetime NULL DEFAULT NULL COMMENT 'last update time',
    `last_sync_time` datetime NULL DEFAULT NULL COMMENT 'last sync time',
    `token_id` bigint(20) NULL DEFAULT NULL COMMENT 'repository token id',
    `status` tinyint(1) NULL DEFAULT 1 COMMENT 'status',
    `is_deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'is deleted',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE INDEX `domain`(`domain`, `owner`, `repository_name`) USING BTREE,
    INDEX `repository_type`(`repository_type`) USING BTREE,
    INDEX `store_type`(`store_type`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'repository table' ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
