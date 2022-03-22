package workflow

import (
	"reflect"
	"testing"
)

func Test_getResults(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		want    []string
		wantErr bool
	}{
		{
			name:    "empty line",
			body:    "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "actual line",
			body:    "var c1Arr = new Array(\"\",\"to test sb./sth.\",\"to test sth.\");",
			want:    []string{"to test sb./sth.", "to test sth."},
			wantErr: false,
		},
		{
			name:    "actual line",
			body:    "var c1Arr = new Array(\"\",\"to test sb./sth.\",\"to test sth.\");",
			want:    []string{"to test sb./sth.", "to test sth."},
			wantErr: false,
		},
		{
			name:    "results have comma",
			body:    "var c1Arr = new Array(\"\",\"to test, sb./sth.\",\"to test sth.\");",
			want:    []string{"to test, sb./sth.", "to test sth."},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getResults(1, tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("getResults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getResults() got = %v, want %v", got, tt.want)
			}
		})
	}
}
