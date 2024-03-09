package consts

// 是否开启磁盘持久化
const DiskPersistence = true

type FieldName string

const (
	PersistTerm    FieldName = "term"
	PersistVoteFor FieldName = "voteFor"
	PersistLogs    FieldName = "logs"
)
