// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"short/internal/config"
	"short/model"
	"short/pkg/idgenerator"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/syncx"
)

type ServiceContext struct {
	Config       config.Config
	UrlModel     model.ShortUrlMapModel
	Sequence     idgenerator.Generator
	BlackMap     map[string]struct{}
	RedisClient  *redis.Redis
	SingleFlight syncx.SingleFlight
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

	return &ServiceContext{
		Config:       c,
		UrlModel:     urlModel,
		Sequence:     sequenceModel,
		BlackMap:     blackMap,
		RedisClient:  rds,
		SingleFlight: syncx.NewSingleFlight(),
	}
}
