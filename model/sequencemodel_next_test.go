package model

import (
	"context"
	"testing"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func Test_customSequenceModel_Next(t *testing.T) {
	// 使用硬编码的 DSN，不依赖配置文件
	// 或者使用环境变量
	dsn := "root:123456@tcp(127.0.0.1:13306)/url?parseTime=true"
	
	conn := sqlx.NewMysql(dsn)
	m := NewSequenceModel(conn)

	// 获取当前ID作为基准
	ctx := context.Background()
	quence, err := m.FindOneByStub(ctx, "a")
	if err != nil {
		t.Fatalf("Failed to find sequence: %v", err)
	}
	baseID := quence.Id

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "first call",
			wantErr: false,
		},
		{
			name:    "second call",
			wantErr: false,
		},
		{
			name:    "third call",
			wantErr: false,
		},
	}

	// 记录上一次的ID
	lastID := baseID

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := m.Next(ctx)
			
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Next() failed: %v", gotErr)
				}
				return
			}
			
			if tt.wantErr {
				t.Fatal("Next() succeeded unexpectedly")
			}

			// 验证ID递增
			if got <= lastID {
				t.Errorf("Next() = %v, should be greater than lastID %v", got, lastID)
			}

			// 验证ID递增量为1
			if got != lastID+1 {
				t.Errorf("Next() = %v, expected %v (increment by 1)", got, lastID+1)
			}

			t.Logf("✓ %s: got ID %d, previous was %d", tt.name, got, lastID)
			lastID = got
		})
	}
}

// 测试并发安全性
func Test_customSequenceModel_Next_Concurrent(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:13306)/url?parseTime=true"
	conn := sqlx.NewMysql(dsn)
	m := NewSequenceModel(conn)

	ctx := context.Background()
	
	// 获取基准ID
	quence, _ := m.FindOneByStub(ctx, "a")
	baseID := quence.Id

	// 并发调用
	concurrency := 10
	ids := make(chan uint64, concurrency)
	
	for i := 0; i < concurrency; i++ {
		go func() {
			id, err := m.Next(ctx)
			if err != nil {
				t.Errorf("Concurrent Next() failed: %v", err)
				ids <- 0
				return
			}
			ids <- id
		}()
	}

	// 收集结果
	receivedIDs := make(map[uint64]bool)
	for i := 0; i < concurrency; i++ {
		id := <-ids
		if id == 0 {
			continue
		}
		
		// 检查是否有重复ID
		if receivedIDs[id] {
			t.Errorf("Duplicate ID detected: %d", id)
		}
		receivedIDs[id] = true
		
		// 验证ID大于基准
		if id <= baseID {
			t.Errorf("ID %d should be greater than base ID %d", id, baseID)
		}
	}

	t.Logf("✓ Concurrent test passed: %d unique IDs generated", len(receivedIDs))
}
