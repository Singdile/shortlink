package connect

import "testing"

func TestCheckURL(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		url     string
		want    bool
		wantErr bool
	}{
		{
			name:    "valid url",
			url:     "https://www.baidu.com",
			want:    true,
			wantErr: false,
		},
		{
			name:    "invalid url",
			url:     "http://invalid.url",
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := CheckURL(tt.url)
			if gotErr != nil { // gotErr is not nil, 表示请求失败
				if !tt.wantErr { //期望请求成功，但是实际请求失败
					t.Errorf("CheckURL() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr { //期望请求失败，但是实际请求成功
				t.Fatal("CheckURL() succeeded unexpectedly")
			}
			if got != tt.want { // got is not equal to want, 表示请求结果不符合预期
				t.Errorf("CheckURL() = %v, want %v", got, tt.want)
			}
		})
	}

}
