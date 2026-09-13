CREATE TABLE `users` (
  `id` INT NOT NULL,
  `name` VARCHAR(50) NOT NULL,
  `age` INT DEFAULT 0,
  PRIMARY KEY (`id`)
);
INSERT INTO `users` (`id`,`name`,`age`) VALUES
(1,'Alice',20),
(2,'Bob',30),
(3,'Carol',40);

CREATE TABLE `old_table` (
  `id` INT NOT NULL,
  `value` VARCHAR(50),
  PRIMARY KEY (`id`)
);
INSERT INTO `old_table` (`id`,`value`) VALUES
(1,'old'),
(2,'data');

CREATE TABLE `same_table` (
  `id` INT NOT NULL,
  `value` VARCHAR(50),
  PRIMARY KEY (`id`)
);
INSERT INTO `same_table` (`id`,`value`) VALUES
(1,'same'),
(2,'unchanged');

CREATE TABLE `composite_pk` (
  `a` INT NOT NULL,
  `b` INT NOT NULL,
  `value` VARCHAR(50),
  PRIMARY KEY (`a`,`b`)
);
INSERT INTO `composite_pk` (`a`,`b`,`value`) VALUES
(1,1,'one'),
(1,2,'two'),
(2,1,'three');