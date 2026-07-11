package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ShortUrlMapModel = (*customShortUrlMapModel)(nil)

type (
	// ShortUrlMapModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShortUrlMapModel.
	ShortUrlMapModel interface {
		shortUrlMapModel
		withSession(session sqlx.Session) ShortUrlMapModel

		// FindByLurl 自定义方法，查找Lurl 对应的 surl
		FindByHashAndLurl(ctx context.Context, lurl string, hashvalue uint64) (*ShortUrlMap, error)
	}

	customShortUrlMapModel struct {
		*defaultShortUrlMapModel
	}
)

// NewShortUrlMapModel returns a model for the database table.
func NewShortUrlMapModel(conn sqlx.SqlConn) ShortUrlMapModel {
	return &customShortUrlMapModel{
		defaultShortUrlMapModel: newShortUrlMapModel(conn),
	}
}

func (m *customShortUrlMapModel) withSession(session sqlx.Session) ShortUrlMapModel {
	return NewShortUrlMapModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customShortUrlMapModel) FindByHashAndLurl(ctx context.Context, lurl string, hashvalue uint64) (*ShortUrlMap, error) {
	// MurmurHash3
	lurlHash := hashvalue

	// 查询数据库
	sqlStr := fmt.Sprintf("select %s from %s where `lurl_hash` = ? and `lurl` = ? limit 1", shortUrlMapRows, m.table)

	var resp ShortUrlMap
	err := m.conn.QueryRowCtx(ctx, &resp, sqlStr, lurlHash, lurl)

	// 返回查询结果
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}

}
