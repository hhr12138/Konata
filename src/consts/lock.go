package consts

type LockName int

const IS_OVER_LOCK = false

const (
	TERM LockName = iota
	STATUS
	LOG
	NEXT_INDEX
	MATCH_INDEX
	LAST_APPLIED

	OVER_LOCK
)

var lockNameToString = map[LockName]string{
	STATUS:       "status",
	TERM:         "term",
	LOG:          "log",
	OVER_LOCK:    "overLock",
	NEXT_INDEX:   "nextIndex",
	MATCH_INDEX:  "matchIndex",
	LAST_APPLIED: "lastApplied",
}

// 上锁顺序数组
var LockOrder = []LockName{TERM, STATUS, LOG, NEXT_INDEX, MATCH_INDEX, LAST_APPLIED}

func (lockName LockName) String() string {
	if name, ok := lockNameToString[lockName]; ok {
		return name
	}
	return "undefined"
}
