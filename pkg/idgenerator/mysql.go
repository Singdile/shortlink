package idgenerator

import (
	"context"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/jmoiron/sqlx"
)

// 定义Mysql生成器类型
type MysqlGenerator struct {
	// 数据库连接
	db    *sqlx.DB
	stmt  *sqlx.Stmt
	table string
	mu    sync.Mutex
}

// MySqlConfig 配置 MySQL 发号器
type MySqlConfig struct {
	DSN             string        // 数据源名称
	Table           string        // 表名
	MaxOpenConns    int           // 最大打开连接数
	MaxIdleConns    int           // 最大空闲连接数
	ConnMaxLifetime time.Duration // 连接最大生命周期，单位分钟，比如30m,表示30分钟
}

// NewMysqlGenerator 创建一个新的 MySQL 发号器
func NewMysqlGenerator(cfg *MySqlConfig) (*MysqlGenerator, error) {
	// 创建数据库连接
	db, err := sqlx.Open("mysql", cfg.DSN)
	if err != nil {
		return nil, err
	}

	// 设置默认值
	if cfg.Table == "" {
		cfg.Table = "sequence"
	}
	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = 10
	}

	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = 5
	}
	if cfg.ConnMaxLifetime == 0 {
		cfg.ConnMaxLifetime = 30 * time.Minute
	}

	// 配置连接池参数
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// 测试连接
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 预编译 SQL 语句
	// stmt 能加快该连接的后续相同的sql语句的执行速度
	query := fmt.Sprintf("replace into %s (stub) values ('a')", cfg.Table)
	stmt, err := db.Preparex(query)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	// 返回 MysqlGenerator 实例
	return &MysqlGenerator{
		db:    db,
		stmt:  stmt,
		table: cfg.Table,
	}, nil
}

// Next 获取下一个发号器的值
func (g *MysqlGenerator) Next(ctx context.Context) (uint64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// 执行预编译的REPLACE语句
	result, err := g.stmt.ExecContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to execute replace statement: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return uint64(id), nil
}

// 关闭数据库连接和预编译语句
func (g *MysqlGenerator) Close() error {
	if g.db != nil {
		return g.db.Close()
	}
	return nil
}
