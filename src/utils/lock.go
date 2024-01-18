package utils

import (
	"github.com/hhr12138/Konata/src/consts"
	"sync"
)

// 上锁顺序数组
var lockOrder = []consts.LockName{consts.TERM, consts.STATUS, consts.LOG}

var lockMap = map[consts.LockName]*sync.Mutex{}

func init() {
	DPrintf("utils\\lock.go init")
	for _, name := range lockOrder {
		if _, ok := lockMap[name]; ok {
			continue
		}
		lockMap[name] = new(sync.Mutex)
	}
}

func Lock(targetLocks map[consts.LockName]struct{}) {
	// 如果死锁，可以简单的换成一个全局锁上锁
	if consts.IS_OVER_LOCK {
		lockMap[consts.OVER_LOCK].Lock()
	} else {
		for _, name := range lockOrder {
			if _, ok := targetLocks[name]; !ok {
				continue
			}
			lockMap[name].Lock()
		}
	}
}

func Unlock(targetLocks map[consts.LockName]struct{}) {
	// 如果死锁，可以简单的换成一个全局锁上锁
	if consts.IS_OVER_LOCK {
		lockMap[consts.OVER_LOCK].Unlock()
	} else {
		for i := len(lockOrder) - 1; i >= 0; i-- {
			if _, ok := targetLocks[lockOrder[i]]; !ok {
				continue
			}
			lockMap[lockOrder[i]].Unlock()
		}
	}
}

func ShowLockLog(level consts.LogLevel, serviceId string, term int, offset int, targetLocks map[consts.LockName]struct{}) {
	for _, name := range lockOrder {
		if _, ok := targetLocks[name]; !ok {
			continue
		}
		Printf(level, serviceId, term, offset, "%v lock", name.String())
	}
}

func ShowUnlockLog(level consts.LogLevel, serviceId string, term int, offset int, targetLocks map[consts.LockName]struct{}) {
	for i := len(lockOrder) - 1; i >= 0; i-- {
		if _, ok := targetLocks[lockOrder[i]]; !ok {
			continue
		}
		Printf(level, serviceId, term, offset, "%v unlock", lockOrder[i].String())
	}
}
