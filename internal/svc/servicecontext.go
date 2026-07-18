// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"short/internal/config"
	"short/model"
	"short/pkg/idgenerator"

	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/syncx"
)

type ServiceContext struct {
	Config       config.Config
	UrlModel     model.ShortUrlMapModel
	Sequence     idgenerator.Generator
	BlackMap     map[string]struct{} //过滤词汇表
	RedisClient  *redis.Redis
	SingleFlight syncx.SingleFlight
	LinkBloom    *bloom.Filter //短链接布隆过滤器
}

func NewServiceContext(c config.Config) *ServiceContext {
	rds := redis.MustNewRedis(c.RedisConf)
	urlModel := model.NewShortUrlMapModel(sqlx.NewMysql(c.ShortUrlDB.DSN))
	sequenceModel, err := idgenerator.NewMysqlGenerator(&c.SequenceDB) //发号器实实例
	if err != nil {
		panic(err)
	}

	blackMap := make(map[string]struct{})
	for _, v := range c.BlackList {
		blackMap[v] = struct{}{}
	}

	// 初始化bloom filter
	// 1.初始化对应的redis
	biredis := redis.MustNewRedis(c.BizRedis)
	// 2.初始化bloom 实例
	linkBloom := bloom.New(biredis, c.BloomFilter.Key, c.BloomFilter.Bits)

	return &ServiceContext{
		Config:       c,
		UrlModel:     urlModel,
		Sequence:     sequenceModel,
		BlackMap:     blackMap,
		RedisClient:  rds,
		LinkBloom:    linkBloom,
		SingleFlight: syncx.NewSingleFlight(),
	}
}

// loadDataToBloomFilter 加载已有的短链接数据到布隆过滤器
func loadDataToBloomFilter() {

}
