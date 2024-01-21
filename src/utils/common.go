package utils

// 通用操作

type CASHandler func() bool

func OperationByCAS(nowTerm, oldTerm int, handler CASHandler) bool{
	if oldTerm != nowTerm{
		return false
	}
	return handler()
}