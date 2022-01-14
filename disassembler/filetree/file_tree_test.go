package filetree

import (
	"reflect"
	"testing"
)

func TestFileTree_FindNodeFromPath(t *testing.T) {
	type fields struct {
		Root *FileNode
	}
	type args struct {
		pathStr string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *FileNode
	}{
		{
			name: "Success",
			fields: fields{
				Root: &FileNode{
					Children: map[string]*FileNode{
						"a": {
							Children: map[string]*FileNode{
								"bb": {
									Children: map[string]*FileNode{
										"ccc": {
											Name: "ccc",
										},
									},
								},
							},
						},
					},
				},
			},
			args: args{
				pathStr: "a/bb/ccc/",
			},
			want: &FileNode{
				Name: "ccc",
			},
		},
		{
			name: "NotFound",
			fields: fields{
				Root: &FileNode{
					Children: map[string]*FileNode{
						"a": {},
						"b": {},
					},
				},
			},
			args: args{
				pathStr: "c/cc/ccc",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := &FileTree{
				Root: tt.fields.Root,
			}
			if got := tree.FindNodeFromPath(tt.args.pathStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindNodeFromPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
