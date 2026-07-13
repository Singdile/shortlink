package base62

import "testing"

func TestUint2string(t *testing.T) {
	tests := []struct {
		name    string
		num     uint64
		want    string
		wantErr bool
	}{
		{
			name: "test-62",
			num:  77,
			want: "1F",
		},
		{
			name: "test-100",
			num:  100,
			want: "1c",
		},
		{
			name: "test-123456789",
			num:  123456789,
			want: "8M0kX",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Uint2string(tt.num)
			if got != tt.want {
				t.Errorf("Uint2string() = %v, want %v", got, tt.want)
			}
		})

	}
}

func TestString2Uint(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want uint64
	}{
		{
			name: "test-1F",
			str:  "1F",
			want: 77,
		},
		{
			name: "test-1c",
			str:  "1c",
			want: 100,
		},
		{
			name: "test-8M0kX",
			str:  "8M0kX",
			want: 123456789,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := String2Uint(tt.str); got != tt.want {
				t.Errorf("String2Uint() = %v, want %v", got, tt.want)
			}
		})
	}

}
