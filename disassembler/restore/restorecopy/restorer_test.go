package restorecopy

import "testing"

func Test_restorer_absPath(t *testing.T) {
	type fields struct {
		tmpWorkDir string
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Success1",
			fields: fields{
				tmpWorkDir: "/a/aa",
			},
			args: args{
				path: "/b/bb",
			},
			want: "/b/bb",
		},
		{
			name: "Success2",
			fields: fields{
				tmpWorkDir: "/a/aa",
			},
			args: args{
				path: "b/bb",
			},
			want: "/a/aa/b/bb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &restorer{
				tmpWorkDir: tt.fields.tmpWorkDir,
			}
			if got := r.absPath(tt.args.path); got != tt.want {
				t.Errorf("absPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
