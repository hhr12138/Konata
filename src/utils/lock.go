package utils

import "github.com/hhr12138/Konata/src/consts"

func ShowLockLog(level consts.LogLevel, serviceId string, term int, offset int, funcName string, targetLocks map[consts.LockName]struct{}) {
	for _, name := range consts.LockOrder {
		if _, ok := targetLocks[name]; !ok {
			continue
		}
		Printf(level, serviceId, term, offset, "%v %v lock", funcName, name.String())
	}
}

func ShowUnlockLog(level consts.LogLevel, serviceId string, term int, offset int, funcName string, targetLocks map[consts.LockName]struct{}) {
	for i := len(consts.LockOrder) - 1; i >= 0; i-- {
		if _, ok := targetLocks[consts.LockOrder[i]]; !ok {
			continue
		}
		Printf(level, serviceId, term, offset, "%v %v unlock", funcName, consts.LockOrder[i].String())
	}
}

func GetLockMap(args ...consts.LockName) map[consts.LockName]struct{}{
	mp := make(map[consts.LockName]struct{},len(args))
	for _,name := range args{
		if _,ok := mp[name]; ok{
			continue
		}
		mp[name] = struct{}{}
	}
	return mp
}
