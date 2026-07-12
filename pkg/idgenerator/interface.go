package idgenerator

import "context"

type Generator interface {
	// Next 获取下一个发号器的值
	Next(context.Context) (uint64, error)

	// Close 关闭相关资源
	Close() error
}
