// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"short/pkg/idgenerator"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	ShortUrlDB  ShortUrlDB
	SequenceDB  idgenerator.MySqlConfig
	BlackList   []string
	Domain      string          //短域名
	RedisConf   redis.RedisConf // 对应 yaml 中的 RedisConf
	BloomFilter BloomFilter
	BizRedis    redis.RedisConf
}

type ShortUrlDB struct {
	DSN string
}

// bloom filter config
type BloomFilter struct {
	Key  string
	Bits uint
}
