package entity

import (
	"github.com/hhr12138/Konata/src/consts"
	"github.com/hhr12138/Konata/src/utils"
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
			locks := Locks{}
			locks.InitLocks()
			utils.ShowLockLog(consts.INFO, "raft0", 1, 1, "[test]",tt.args.targetLocks)
			locks.Lock(tt.args.targetLocks)
			locks.Unlock(tt.args.targetLocks)
			utils.ShowUnlockLog(consts.INFO, "raft0", 1, 1, "[test]", tt.args.targetLocks)
		})
	}
}
