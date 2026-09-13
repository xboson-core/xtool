CREATE DATABASE  `test` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `test`;
--
-- Table structure for table `users`
--
DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` INT NOT NULL,
  `age` INT DEFAULT 0,
  `name` VARCHAR(50) NOT NULL,
  `email` VARCHAR(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
);
/*!40101 SET character_set_client = @saved_cs_client */;

INSERT INTO `users` (`id`,`name`,`age`,`email`) VALUES
    (1,'Alice',21,'alice@example.com'),
    (2,'Bob',30,'bob@example.com'),
    (4,'David',25,'david@example.com');

CREATE TABLE `new_table` (
  `id` INT NOT NULL,
  `title` VARCHAR(100) NOT NULL,
  PRIMARY KEY (`id`)
);
INSERT INTO `new_table` (`id`,`title`) VALUES
    (1,'hello'),
    (2,'world');

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
  `select` INT default 10,
  `value` VARCHAR(50),
  PRIMARY KEY (`a`,`b`)
);
INSERT INTO `composite_pk` (`a`,`b`,`value`) VALUES
    (1,1,'one changed'),
    (1,3,'new row'),
    (2,1,'three');

CREATE Table `create` (
  `insert` INT NOT NULL,
  `alert` INT DEFAULT 0,
  PRIMARY KEY (`insert`)
);
INSERT INTO `create` (`insert`,`alert`) VALUES (1, 2);
INSERT INTO `create` VALUES (3, 4);