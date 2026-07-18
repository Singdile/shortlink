// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package main

import (
	"flag"
	"fmt"
	"os"

	"short/internal/config"
	"short/internal/handler"
	"short/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/shortlink-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 环境变量覆盖
	if dsn := os.Getenv("MYSQL_URL"); dsn != "" {
		c.ShortUrlDB.DSN = dsn
		c.SequenceDB.DSN = dsn
	}
	if host := os.Getenv("REDIS_URL"); host != "" {
		c.RedisConf.Host = host
		c.BizRedis.Host = host
	}

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
