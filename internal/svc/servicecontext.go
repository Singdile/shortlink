// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"short/internal/config"
	"short/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config   config.Config
	UrlModel model.ShortUrlMapModel
	Sequence model.SequenceModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	urlModel := model.NewShortUrlMapModel(sqlx.NewMysql(c.ShortUrlDB.DSN))
	sequenceModel := model.NewSequenceModel(sqlx.NewMysql(c.SequenceDB.DSN))

	return &ServiceContext{
		Config:   c,
		UrlModel: urlModel,
		Sequence: sequenceModel,
	}
}
