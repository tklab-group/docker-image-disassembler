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

func Test_isExtraFormatCopy(t *testing.T) {
	type args struct {
		command string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ExtraFormatCopy",
			args: args{
				command: "/bin/sh -c #(nop) COPY file:50563a97010fd7ce1ceebd1fa4f4891ac3decdf428333fb2683696f4358af6c2 in /",
			},
			want: true,
		},
		{
			name: "ExtraFormatCMD",
			args: args{
				command: `/bin/sh -c #(nop)  CMD ["/hello"]`,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isExtraFormatCopy(tt.args.command); got != tt.want {
				t.Errorf("isExtraFormatCopy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDstPathFromExtraFormatCopy(t *testing.T) {
	type args struct {
		command string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				command: "/bin/sh -c #(nop) COPY file:50563a97010fd7ce1ceebd1fa4f4891ac3decdf428333fb2683696f4358af6c2 in /",
			},
			want:    "/",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDestPathFromExtraFormatCopy(tt.args.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDestPathFromExtraFormatCopy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getDestPathFromExtraFormatCopy() got = %v, want %v", got, tt.want)
			}
		})
	}
}
