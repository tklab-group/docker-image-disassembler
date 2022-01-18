package main

import (
	"bytes"
	"testing"
)

func Test_run(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				args: []string{"testdata/a.txt", "testdata/b.txt"},
			},
			wantOut: `file A: testdata/a.txt
file B: testdata/b.txt

それぞれのパスの数
	A: 3
	B: 4

一致しているパスの数: 2
データサイズも一致しているパスの数: 1`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			err := run(tt.args.args, out)
			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("run() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
