-- 2026-09-13 13:01:16.344926 +0800 CST
-- Base: D:\javaee-project\xboson-open-source\xtool\test\a1.sql
-- Diff: D:\javaee-project\xboson-open-source\xtool\test\a2.sql
-- 1;
CREATE DATABASE  `test` /*!40100 DEFAULT CHARACTER SET utf8 */;
-- 2;
USE `test`;
ALTER TABLE `test`.`test.users` ADD COLUMN email VARCHAR(100) default null;
-- 15;
DELETE FROM `test`.`test.users` where `id` = 3;
-- 15;
UPDATE `test`.`test.users` 
	SET `age` = 21, `email` = 'alice@example.com' 
	WHERE `id` = 1;
-- 15;
UPDATE `test`.`test.users` 
	SET `email` = 'bob@example.com' 
	WHERE `id` = 2;
-- 15;
INSERT INTO test.users (`id`, `age`, `name`, `email`) VALUES 
	(4, 25, 'David', 'david@example.com');
-- 21;
CREATE TABLE `new_table` (  `id` INT NOT NULL,  `title` VARCHAR(100) NOT NULL,  PRIMARY KEY (`id`));
-- 24;
INSERT INTO `new_table` (`id`,`title`) VALUES    (1,'hello'),    (2,'world');
ALTER TABLE `test`.`test.composite_pk` ADD COLUMN `select` INT default 10;
-- 45;
DELETE FROM `test`.`test.composite_pk` where `a` = 1 AND `b` = 2;
-- 45;
UPDATE `test`.`test.composite_pk` 
	SET `value` = 'one changed' 
	WHERE `a` = 1 AND `b` = 1;
-- 45;
INSERT INTO test.composite_pk (`a`, `b`, `value`) VALUES 
	(1, 3, 'new row');
-- 51;
CREATE Table `create` (  `insert` INT NOT NULL,  `alert` INT DEFAULT 0,  PRIMARY KEY (`insert`));
-- 52;
INSERT INTO `create` (`insert`,`alert`) VALUES (1, 2);
-- 53;
INSERT INTO `create` VALUES (3, 4);
DROP TABLE IF EXISTS `test`.`test.old_table`;
