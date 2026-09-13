ALTER TABLE `users` ADD COLUMN email VARCHAR(100) default null;
DELETE FROM `users` where `id` = 3
UPDATE `users` 
	SET `age` = 21, `email` = 'alice@example.com' 
	WHERE `id` = 1
UPDATE `users` 
	SET `email` = 'bob@example.com' 
	WHERE `id` = 2
INSERT INTO users (`id`, `age`, `name`, `email`) VALUES 
	(4, 25, 'David', 'david@example.com')


CREATE TABLE `new_table` (
  `id` INT NOT NULL,
  `title` VARCHAR(100) NOT NULL,
  PRIMARY KEY (`id`)
);

INSERT INTO `new_table` (`id`,`title`) VALUES
    (1,'hello'),
    (2,'world');
ALTER TABLE `composite_pk` ADD COLUMN `select` INT default 10;
DELETE FROM `composite_pk` where `a` = 1 AND `b` = 2
UPDATE `composite_pk` 
	SET `value` = 'one changed' 
	WHERE `a` = 1 AND `b` = 1
INSERT INTO composite_pk (`a`, `b`, `value`) VALUES 
	(1, 3, 'new row')


CREATE Table `create` (
  `insert` INT NOT NULL,
  `alert` INT DEFAULT 0,
  PRIMARY KEY (`insert`)
);

INSERT INTO `create` (`insert`,`alert`) VALUES (1, 2);

INSERT INTO `create` VALUES (3, 4);
DROP TABLE IF EXISTS `old_table`;
