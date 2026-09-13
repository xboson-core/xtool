-- 2026-09-13 14:53:32.5171157 +0800 CST
-- Base: D:\javaee-project\xboson-open-source\xtool\test\a1.sql
-- Diff: D:\javaee-project\xboson-open-source\xtool\test\a2.sql
-- 2;
USE `test`;
-- 27;
CREATE TABLE `new_table` (
  `id` INT NOT NULL,
  `title` VARCHAR(100) NOT NULL,
  PRIMARY KEY (`id`)
);
INSERT INTO `new_table` (`id`,`title`) VALUES
    (1,'hello'),
    (2,'world');
-- 57;
CREATE Table `create` (
  `insert` INT NOT NULL,
  `alert` INT DEFAULT 0,
  PRIMARY KEY (`insert`)
);
INSERT INTO `create` (`insert`,`alert`) VALUES (1, 2);
INSERT INTO `create` VALUES (3, 4);
ALTER TABLE `test`.`users` ADD COLUMN email VARCHAR(100) default null;
DELETE FROM `test`.`users` where `id` = 3;
UPDATE `test`.`users` 
	SET `age` = 21, `email` = 'alice@example.com' 
	WHERE `id` = 1;
UPDATE `test`.`users` 
	SET `email` = 'bob@example.com' 
	WHERE `id` = 2;
INSERT INTO `test`.`users` (`id`, `age`, `name`, `email`) VALUES 
	(4, 25, 'David', 'david@example.com');
DROP TABLE IF EXISTS `test`.`old_table`;
ALTER TABLE `test`.`composite_pk` ADD COLUMN `select` INT default 10;
DELETE FROM `test`.`composite_pk` where `a` = 1 AND `b` = 2;
UPDATE `test`.`composite_pk` 
	SET `value` = 'one changed' 
	WHERE `a` = 1 AND `b` = 1;
INSERT INTO `test`.`composite_pk` (`a`, `b`, `value`) VALUES 
	(1, 3, 'new row');
