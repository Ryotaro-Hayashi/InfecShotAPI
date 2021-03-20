-- -----------------------------------------------------
-- Schema proto_api
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `proto_api` DEFAULT CHARACTER SET utf8mb4 ;
USE `proto_api` ;

SET CHARSET utf8mb4;

-- -----------------------------------------------------
-- Table `proto_api`.`user`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `proto_api`.`user` (
    `id` VARCHAR(128) NOT NULL COMMENT 'ユーザID',
    `auth_token` VARCHAR(128) NOT NULL COMMENT '認証トークン',
    `name` VARCHAR(64) NOT NULL COMMENT 'ユーザ名',
    `high_score` INT UNSIGNED NOT NULL COMMENT 'ハイスコア',
    PRIMARY KEY (`id`),
    INDEX `idx_auth_token` (`auth_token` ASC))
ENGINE = InnoDB
COMMENT = 'ユーザ';
