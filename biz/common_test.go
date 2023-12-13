package biz

import "testing"

func TestStrComparion(t *testing.T) {
	type args struct {
		sType     string
		curStr    string
		targetStr string
	}
	tests := []struct {
		name    string
		args    args
		wantB   bool
		wantErr error
	}{
		// TODO: Add test cases.
		{"测试成功", args{"re", "成功", "成功|已存在|重复"},true, nil},
		{"测试已存在", args{"re", "已存在", "成功|已存在|重复"},true, nil},
		{"测试重复", args{"re", "重复", "成功|已存在|重复"},true, nil},
		{"测试未匹配到", args{"re", "失败", "成功|已存在|重复"}, false, nil},
		{"测试in", args{"in", "失败", "成功"}, false, nil},
		{"测试in", args{"in", "重复", "名称已重复"},  true, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotB, err := StrComparion(tt.args.sType, tt.args.curStr, tt.args.targetStr)
			if (err != nil)  {
				t.Errorf("StrComparion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotB != tt.wantB {
				t.Errorf("StrComparion() gotB = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}
