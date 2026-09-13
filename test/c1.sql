/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sys_mdm002` (
  `typecd` varchar(100) NOT NULL DEFAULT '' COMMENT '类别编码',
  `dictcd` varchar(100) NOT NULL DEFAULT '' COMMENT '字典编码',
  `dictnm` varchar(300) DEFAULT NULL COMMENT '字典名称',
  `shortkey` varchar(100) DEFAULT NULL COMMENT '快捷码',
  `status` char(1) DEFAULT NULL COMMENT '状态',
  `mark` varchar(600) DEFAULT NULL COMMENT '说明',
  `createdt` datetime DEFAULT NULL COMMENT '创建时间',
  `updatedt` datetime DEFAULT NULL COMMENT '更新时间',
  `version` varchar(100) NOT NULL DEFAULT '' COMMENT '版本',
  PRIMARY KEY (`typecd`,`dictcd`,`version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='数据字典表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_mdm002`
--

LOCK TABLES `sys_mdm002` WRITE;
/*!40000 ALTER TABLE `sys_mdm002` DISABLE KEYS */;
INSERT INTO `sys_mdm002` VALUES 
  ('110000000000','110100000000','市辖区','SXQ','1',NULL,
  '2014-11-07 16:08:13','2014-11-07 16:08:13','1.0');
-- 表没有导致后面的数据不存在, 全都输出delete from .....
INSERT INTO `sys_mdm002` VALUES 
  ('110000000000','110200000000','县','X','1',NULL,
  '2014-11-07 16:08:13','2014-11-07 16:08:13','1.0'),
  ('110100000000','110101000000','东城区','DCQ','1',NULL,
  '2014-11-07 16:08:13','2014-11-07 16:08:13','1.0'),
  ('110100000000','110102000000','西城区','XCQ','1',NULL,
  '2014-11-07 16:08:13','2014-11-07 16:08:13','1.0');

