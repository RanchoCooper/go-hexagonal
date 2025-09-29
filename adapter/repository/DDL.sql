CREATE DATABASE `go-hexagonal`;

USE `go-hexagonal`;

DROP TABLE IF EXISTS `example`;

CREATE TABLE `example` (
    `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Primary key ID',
    `name` VARCHAR(255) NOT NULL COMMENT 'Name',
    `alias` VARCHAR(255) DEFAULT NULL COMMENT 'Alias',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Deletion time',
    PRIMARY KEY (`id`),
    KEY `idx_name` (`name`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Hexagonal example table';
