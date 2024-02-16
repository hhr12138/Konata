package consts

const (
	ELECTION_TIME  = 1200 // 选举超时, 单位ms
	HEARTBEAT_TIME = 5    // 从节点心跳检测时间,单位ms,不能过大,举个例子,a选举超时为140ms后,b选举超时为149ms后,但二者都在150ms开始新的检测,导致同时选举.因此检测时间过大会导致随机基本单元过大,同时选举概率增高
	HEART_TIME     = 120  // 主节点发送心跳时间,单位ms
	WAIT_VOTE_TIME = 1    // 等待选票睡眠时间,单位ms,同理,不能过大
)
