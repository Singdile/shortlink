-- 回滚 lurl_hash 字段类型为 BIGINT UNSIGNED
ALTER TABLE `short_url_map` MODIFY COLUMN `lurl_hash` BIGINT UNSIGNED DEFAULT NULL COMMENT '长链接的 MurmurHash3 值，用于高效索引';
