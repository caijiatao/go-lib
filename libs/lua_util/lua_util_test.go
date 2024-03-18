package lua_util

import (
	"testing"
)

func TestCheckLua(t *testing.T) {
	type args struct {
		script string
		opts   []CheckOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "正常脚本",
			args: args{
				script: `function truncateString(value)
    return string.sub(value, 1, 5)
end`,
			},
			wantErr: false,
		},
		{
			name: "错误脚本",
			args: args{
				script: "1111",
			},
			wantErr: true,
		},
		{
			name: "方法不存在",
			args: args{
				script: `function truncateString(value)
    return string.sub(value, 1, 5)
end`,
				opts: []CheckOption{WithCheckConfigFuncName("truncateString")},
			},
			wantErr: false,
		},
		{
			name: "方法不存在",
			args: args{
				script: `function truncateString(value)
    return string.sub(value, 1, 5)
end`,
				opts: []CheckOption{WithCheckConfigFuncName("truncateString1")},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Check(tt.args.script, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("ValidateLuaSyntax() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetLuaTypeValue(t *testing.T) {
	GetLuaTypeValue(nil)
}
