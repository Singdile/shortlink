// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"short/pkg/idgenerator"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	ShortUrlDB ShortUrlDB
	SequenceDB idgenerator.MySqlConfig
}

type ShortUrlDB struct {
	DSN string
}
