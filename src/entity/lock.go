package entity

import (
	"github.com/hhr12138/Konata/src/consts"
	"github.com/hhr12138/Konata/src/utils"
	"sync"
)

type Locks map[consts.LockName]*sync.Mutex

func (locks Locks) InitLocks(){
	utils.DPrintf("lock init")
	for _, name := range consts.LockOrder {
		if _, ok := locks[name]; ok {
			continue
		}
		locks[name] = new(sync.Mutex)
	}
}

func (locks Locks) Lock(targetLocks map[consts.LockName]struct{}) {
	// 如果死锁，可以简单的换成一个全局锁上锁
	if consts.IS_OVER_LOCK {
		locks[consts.OVER_LOCK].Lock()
	} else {
		for _, name := range consts.LockOrder {
			if _, ok := targetLocks[name]; !ok {
				continue
			}
			locks[name].Lock()
		}
	}
}

func (locks Locks) Unlock(targetLocks map[consts.LockName]struct{}) {
	// 如果死锁，可以简单的换成一个全局锁上锁
	if consts.IS_OVER_LOCK {
		locks[consts.OVER_LOCK].Unlock()
	} else {
		for i := len(consts.LockOrder) - 1; i >= 0; i-- {
			if _, ok := targetLocks[consts.LockOrder[i]]; !ok {
				continue
			}
			locks[consts.LockOrder[i]].Unlock()
		}
	}
}
