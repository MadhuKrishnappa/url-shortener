CREATE TABLE `short_url_mapping` (
  `id` int NOT NULL AUTO_INCREMENT,
  `short_url` varchar(100) DEFAULT NULL,
  `long_url` varchar(1000) DEFAULT NULL,
  `expiry_date` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `short_url_mapping_short_url_IDX` (`short_url`) USING BTREE,
  KEY `short_url_mapping_expiry_date_IDX` (`expiry_date`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1010 DEFAULT CHARSET=utf8mb3