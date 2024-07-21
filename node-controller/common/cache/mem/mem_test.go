package mem

import (
	"fmt"
	"testing"
)

func pointStr(s string) *string {
	return &s
}

func pointInt(s int) *int {
	return &s
}

func Test_copyValue(t *testing.T) {
	type tp struct {
		A string
	}
	type args struct {
		src interface{}
		to  interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "string",
			args: args{
				src: "",
				to:  "",
			},
			wantErr: true,
		},
		{
			name: "string",
			args: args{
				src: "123",
				to:  pointStr(""),
			},
			wantErr: false,
		},
		{
			name: "int",
			args: args{
				src: "123",
				to:  pointInt(0),
			},
			wantErr: true,
		},
		{
			name: "struct",
			args: args{
				src: tp{
					A: "11",
				},
				to: &tp{},
			},
			wantErr: false,
		},
		{
			name: "struct",
			args: args{
				src: &tp{
					A: "11",
				},
				to: &tp{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := copyValue(tt.args.src, tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("copyValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println(tt.args)
		})
	}
}
