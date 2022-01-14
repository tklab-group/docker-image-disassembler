package utils

import (
	"reflect"
	"testing"
)

func TestCleanArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "Success",
			args: []string{"a a", " bb ", "     ccc"},
			want: []string{"a a", "bb", "ccc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CleanArgs(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CleanArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
