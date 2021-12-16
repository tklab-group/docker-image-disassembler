package pkginfo

import (
	"reflect"
	"regexp"
	"testing"
)

func Test_extractNamedGroup(t *testing.T) {
	type args struct {
		reg *regexp.Regexp
		s   string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "Success",
			args: args{
				reg: regexp.MustCompile(`(?P<a>a+)(b+)(?P<c>c*).+`),
				s:   "aaabbddd",
			},
			want: map[string]string{
				"a": "aaa",
				"c": "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractNamedGroup(tt.args.reg, tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractNamedGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
