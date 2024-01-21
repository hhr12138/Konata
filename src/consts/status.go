package consts

type Status int

const(
	NULL_STATUS Status = iota-1
	FOLLOWER
	CANDIDATE
	LEADER
)

// 允许的状态转移
var StateTranMap = map[Status]map[Status]bool{
	FOLLOWER: {
		FOLLOWER: true,
		CANDIDATE: true,
	},
	CANDIDATE: {
		FOLLOWER: true,
		LEADER: true,
	},
	LEADER: {
		FOLLOWER: true,
	},
}

var StateToString = map[Status]string{
	FOLLOWER: "follower",
	CANDIDATE: "candidate",
	LEADER: "leader",
}

func (status Status) String() string{
	if name,ok := StateToString[status]; ok{
		return name
	}
	return "undefined"
}