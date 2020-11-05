-- MySQL dump 10.13  Distrib 5.7.29, for Linux (x86_64)
--
-- Host: localhost    Database: pluto
-- ------------------------------------------------------
-- Server version	5.7.29

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `stock`
--

DROP TABLE IF EXISTS `stock`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `stock` (
  `deal_id` int(11) NOT NULL AUTO_INCREMENT,
  `deal_date` date NOT NULL COMMENT '成交日期',
  `stock_code` varchar(6) NOT NULL COMMENT '股票代码',
  `stock_name` varchar(16) NOT NULL COMMENT '股票名称',
  `operation` varchar(10) NOT NULL COMMENT '操作',
  `volume` smallint(5) NOT NULL COMMENT '成交数量',
  `balance` int(11) NOT NULL DEFAULT '0' COMMENT '变动后持股数量',
  `avg_price` decimal(4,2) NOT NULL COMMENT '成交均价',
  `turnover` decimal(7,2) unsigned NOT NULL COMMENT '成交金额',
  `amount` decimal(7,2) NOT NULL COMMENT '发生金额',
  `brokerage` decimal(4,2) unsigned NOT NULL COMMENT '佣金',
  `stamp_tax` decimal(4,2) unsigned NOT NULL COMMENT '印花税',
  `transfer_fee` decimal(3,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '过户费',
  PRIMARY KEY (`deal_id`),
  UNIQUE KEY `uniq_trade` (`deal_date`,`stock_code`,`operation`,`avg_price`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=348 DEFAULT CHARSET=utf8 COMMENT='股票交易记录';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-04-21  9:28:45
