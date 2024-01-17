package utils

import (
	"github.com/hhr12138/Konata/src/consts"
	"testing"
)

func TestLock(t *testing.T) {
	type args struct {
		targetLocks map[consts.LockName]struct{}
	}
	// mock 上锁顺序
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "success test",
			args: args{
				targetLocks: map[consts.LockName]struct{}{
					consts.TERM: struct {
					}{},
					consts.LOG: struct {
					}{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ShowLockLog(consts.INFO, "raft0", 1, 1, tt.args.targetLocks)
			Lock(tt.args.targetLocks)
			Unlock(tt.args.targetLocks)
			ShowUnlockLog(consts.INFO, "raft0", 1, 1, tt.args.targetLocks)
		})
	}
}
