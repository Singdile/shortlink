package murmur3

import (
	"testing"
)

func TestHash(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected uint64
	}{
		{
			name:     "空字符串",
			input:    "",
			expected: 0,
		},
		{
			name:     "简单字符串",
			input:    "hello",
			expected: 14688674573012802306, // 实际运行结果
		},
		{
			name:     "URL字符串",
			input:    "https://example.com/path?key=value",
			expected: 2233349381266069907, // 实际运行结果
		},
		{
			name:     "中文",
			input:    "你好世界",
			expected: 12314529345063139381, // 实际运行结果
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Hash(tt.input)
			if got != tt.expected {
				t.Errorf("Hash(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

// TestHashConsistency 测试一致性：相同输入必须产生相同输出
func TestHashConsistency(t *testing.T) {
	input := "https://example.com/test"

	// 调用多次，确保结果一致
	results := make([]uint64, 10)
	for i := 0; i < 10; i++ {
		results[i] = Hash(input)
	}

	// 所有结果应该相同
	for i := 1; i < 10; i++ {
		if results[i] != results[0] {
			t.Errorf("Hash 一致性失败: 第0次=%v, 第%d次=%v", results[0], i, results[i])
		}
	}
}

// TestHashDifferentInputs 测试分布性：不同输入应该产生不同输出
func TestHashDifferentInputs(t *testing.T) {
	inputs := []string{
		"https://example.com/a",
		"https://example.com/b",
		"https://example.com/c",
		"https://example.org/a", // 不同域名
		"https://example.com/A", // 大小写不同
	}

	results := make(map[uint64]string)
	for _, input := range inputs {
		hash := Hash(input)
		if existing, ok := results[hash]; ok {
			t.Errorf("哈希冲突: %q 和 %q 产生了相同的哈希值 %v", existing, input, hash)
		}
		results[hash] = input
	}
}
