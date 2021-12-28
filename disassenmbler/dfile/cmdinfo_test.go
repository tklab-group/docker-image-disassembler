package dfile

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCmdInfos(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []*CmdInfo
	}{
		{
			name: "Success",
			args: args{
				s: "apt-get update     && apt-get install -y      tzdata      wget=1.21-1ubuntu3     && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime     && rm -rf /usr/share/zoneinfo",
			},
			want: []*CmdInfo{
				{
					MainCmd:  "apt-get",
					Args:     []string{"update"},
					Original: "apt-get update",
				},
				{
					MainCmd:  "apt-get",
					Args:     []string{"install", "-y", "tzdata", "wget=1.21-1ubuntu3"},
					Original: "apt-get install -y      tzdata      wget=1.21-1ubuntu3",
				},
				{
					MainCmd:  "cp",
					Args:     []string{"/usr/share/zoneinfo/Asia/Tokyo", "/etc/localtime"},
					Original: "cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime",
				},
				{
					MainCmd:  "rm",
					Args:     []string{"-rf", "/usr/share/zoneinfo"},
					Original: "rm -rf /usr/share/zoneinfo",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCmdInfos(tt.args.s)
			assert.Equal(t, tt.want, got)
		})
	}
}
