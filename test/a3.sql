ALTER TABLE `users` ADD COLUMN email VARCHAR(100) default null;
DELETE FROM users where `id`=1
DELETE FROM users where `id`=2
DELETE FROM users where `id`=3
INSERT INTO users (`id`,`name`,`age`,`email`) VALUES 
	(1,'Alice',21,'alice@example.com')
INSERT INTO users (`id`,`name`,`age`,`email`) VALUES 
	(2,'Bob',30,'bob@example.com')
INSERT INTO users (`id`,`name`,`age`,`email`) VALUES 
	(4,'David',25,'david@example.com')


CREATE TABLE `new_table` (
  `id` INT NOT NULL,
  `title` VARCHAR(100) NOT NULL,
  PRIMARY KEY (`id`)
);

INSERT INTO `new_table` (`id`,`title`) VALUES
(1,'hello'),
(2,'world');
DELETE FROM same_table where `id`=1
DELETE FROM same_table where `id`=2
INSERT INTO same_table (`id`,`value`) VALUES 
	(1,'same')
INSERT INTO same_table (`id`,`value`) VALUES 
	(2,'unchanged')
DELETE FROM composite_pk where `a`=1 AND `b`=1
DELETE FROM composite_pk where `a`=1 AND `b`=2
DELETE FROM composite_pk where `a`=2 AND `b`=1
INSERT INTO composite_pk (`a`,`b`,`value`) VALUES 
	(1,1,'one changed')
INSERT INTO composite_pk (`a`,`b`,`value`) VALUES 
	(1,3,'new row')
INSERT INTO composite_pk (`a`,`b`,`value`) VALUES 
	(2,1,'three')


CREATE Table `create` (
  `insert` INT NOT NULL,
  `alert` INT DEFAULT 0,
  PRIMARY KEY (`insert`)
);

INSERT INTO `create` (`insert`,`alert`) VALUES (1, 2);

INSERT INTO `create` VALUES (3, 4);
DROP TABLE IF EXISTS old_table;
