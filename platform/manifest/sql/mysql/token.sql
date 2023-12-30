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
-- Table structure for token
-- ----------------------------
DROP TABLE IF EXISTS `token`;
CREATE TABLE `token`  (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
    `owner` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'token owner',
    `owner_id` bigint(20) NULL DEFAULT NULL COMMENT 'token owner id',
    `repository_type` tinyint(4) NOT NULL COMMENT 'repository type (1: gitlab, 2: github)',
    `repository_domain` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'repository api domain',
    `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'token status (1: expired, 2: valid)',
    `token_type` tinyint(4) NOT NULL COMMENT 'token type (1:  personal, 2: organization)',
    `token` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'repository token',
    `expiration_time` datetime NULL DEFAULT NULL COMMENT 'token expiration time',
    `is_deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'is deleted',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
    PRIMARY KEY (`id`) USING BTREE,
    INDEX `owner`(`owner`) USING BTREE,
    INDEX `repository_type`(`repository_type`) USING BTREE,
    INDEX `repository_domain`(`repository_domain`) USING BTREE,
    INDEX `owner_id`(`owner_id`) USING BTREE,
    INDEX `token_type`(`token_type`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 36 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
