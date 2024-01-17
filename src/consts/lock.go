package consts

type LockName int

const IS_OVER_LOCK = false

const (
	TERM LockName = iota
	STATUS
	LOG
	OVER_LOCK
)

var lockNameToString = map[LockName]string{
	STATUS:    "status",
	TERM:      "term",
	LOG:       "log",
	OVER_LOCK: "overLock",
}

func (lockName LockName) String() string {
	if name, ok := lockNameToString[lockName]; ok {
		return name
	}
	return "undefined"
}
