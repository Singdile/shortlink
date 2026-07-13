// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"short/internal/config"
	"short/model"
	"short/pkg/idgenerator"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config   config.Config
	UrlModel model.ShortUrlMapModel
	Sequence idgenerator.Generator
	BlackMap map[string]struct{}
}

func NewServiceContext(c config.Config) *ServiceContext {
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
		Config:   c,
		UrlModel: urlModel,
		Sequence: sequenceModel,
		BlackMap: blackMap,
	}
}
