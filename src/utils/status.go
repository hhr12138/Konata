package utils

import "github.com/hhr12138/Konata/src/consts"

func CanTran(oldStatus, newStatus consts.Status) bool{
	if mp, ok := consts.StateTranMap[oldStatus]; ok{
		if mp[newStatus] {
			return true
		}
	}
	return false
}
