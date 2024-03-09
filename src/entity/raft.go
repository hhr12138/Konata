package entity

type PersistRaft struct {
	Term    int    `json:"term"`
	VoteFor int    `json:"vote_for"`
	Logs    []*Log `json:"logs"`
}
