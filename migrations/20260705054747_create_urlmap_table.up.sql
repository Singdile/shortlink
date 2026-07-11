-- 创建长短链接map表
CREATE TABLE `short_url_map` (
 `id` bigint unsigned not null AUTO_INCREMENT COMMENT '主键',
 `create_at` DATETIME not null DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
 `create_by` varchar(64) not null DEFAULT '' COMMENT '创建者',
 `is_del` tinyint UNSIGNED not null DEFAULT '0' COMMENT '是否删除:0正常1删除',

 `lurl` varchar(2048) DEFAULT NULL COMMENT '长链接',
 `surl` varchar(11) DEFAULT NULL COMMENT '短链接',
 `lurl_hash` bigint unsigned DEFAULT NULL COMMENT '长链接的 MurmurHash3 值，用于高效索引',
 PRIMARY KEY (`id`),
 KEY `idx_lurl_hash` (`lurl_hash`),
 INDEX `idx_is_del` (`is_del`),
 UNIQUE KEY `uk_surl`(`surl`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT = 'url_map';
