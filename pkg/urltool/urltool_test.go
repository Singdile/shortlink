package urltool_test

import (
	"short/pkg/urltool"
	"testing"
)

func TestGetBaseUrl(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		longurl string
		want    string
		wantErr bool
	}{
		{
			name:    "valid url",
			longurl: "https://www.example.com/path/to/resource",
			want:    "resource",
			wantErr: false,
		},
		{
			name:    "invalid url",
			longurl: "xxx/12345",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := urltool.GetBaseUrl(tt.longurl)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetBaseUrl() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetBaseUrl() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				if got != tt.want {
					t.Errorf("GetBaseUrl() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
