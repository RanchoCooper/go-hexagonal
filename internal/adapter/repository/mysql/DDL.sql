DROP TABLE IF EXISTS `example`;

CREATE TABLE `example` (
    `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '名称',
    `alias` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '别称',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Hexagonal示例表';