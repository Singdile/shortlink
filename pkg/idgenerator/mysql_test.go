package idgenerator

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

const dsn = "root:123456@tcp(127.0.0.1:13306)/url?parseTime=true"

func TestMysqlGenerator_Ping(t *testing.T) {
	generator, err := NewMysqlGenerator(&MySqlConfig{
		DSN: dsn,
	})
	if err != nil {
		t.Fatalf("failed to create MysqlGenerator: %v", err)
	}
	defer generator.Close()
	if err := generator.db.Ping(); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}
}

func TestMysqlGenerator_Next(t *testing.T) {
	generator, err := NewMysqlGenerator(&MySqlConfig{
		DSN: dsn,
	})
	if err != nil {
		t.Fatalf("failed to create MysqlGenerator: %v", err)
	}
	defer generator.Close()

	// 获取当前的最大ID
	maxID, err := generator.Next(context.Background())
	if err != nil {
		t.Fatalf("failed to get next ID: %v", err)
	}

	// 获取下面几个id,检测是否是递增的
	for i := 0; i < 5; i++ {
		id, err := generator.Next(context.Background())
		if err != nil {
			t.Fatalf("failed to get next ID: %v", err)
		}
		if id != maxID+uint64(i)+1 {
			t.Errorf("expected ID %d, got %d", maxID+uint64(i)+1, id)
		}
	}
}

// 并发测试，查看是否会重复
func TestMysqlGenerator_ConcurrentNext(t *testing.T) {
	generator, err := NewMysqlGenerator(&MySqlConfig{
		DSN: dsn,
	})
	if err != nil {
		t.Fatalf("failed to create MysqlGenerator: %v", err)
	}
	defer generator.Close()

	ids := make(chan uint64, 20)
	for i := 0; i < 20; i++ {
		go func() {
			id, err := generator.Next(context.Background())
			if err != nil {
				t.Errorf("failed to get next ID: %v", err)
				return
			}
			ids <- id
		}()
	}

	// 检查唯一性
	uniqueIDs := make(map[uint64]bool)
	for i := 0; i < 20; i++ {
		id := <-ids
		if uniqueIDs[id] {
			t.Errorf("duplicate ID found: %d", id)
		}
		uniqueIDs[id] = true
	}
}
