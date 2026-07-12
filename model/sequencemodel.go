package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SequenceModel = (*customSequenceModel)(nil)

type (
	// SequenceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSequenceModel.
	SequenceModel interface {
		sequenceModel
		withSession(session sqlx.Session) SequenceModel

		// 获取下一个发号器的值
		Next(context.Context) (uint64, error)
	}

	customSequenceModel struct {
		*defaultSequenceModel
	}
)

// NewSequenceModel returns a model for the database table.
func NewSequenceModel(conn sqlx.SqlConn) SequenceModel {
	return &customSequenceModel{
		defaultSequenceModel: newSequenceModel(conn),
	}
}

func (m *customSequenceModel) withSession(session sqlx.Session) SequenceModel {
	return NewSequenceModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customSequenceModel) Next(ctx context.Context) (uint64, error) {
	query := fmt.Sprintf("replace into %s (stub) values ('a')", m.table)
	ret, err := m.conn.ExecCtx(ctx, query)

	if err != nil {
		return 0, fmt.Errorf("failed to replace into sequence: %w", err)
	}

	id, err := ret.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return uint64(id), nil
}
