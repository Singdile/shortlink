-- 修改 lurl_hash 字段类型为 BIGINT SIGNED，兼容 Go int64
ALTER TABLE `short_url_map` MODIFY COLUMN `lurl_hash` BIGINT SIGNED DEFAULT NULL COMMENT '长链接的 MurmurHash3 值，用于高效索引';
