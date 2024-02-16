package entity

type Log struct {
	Term    int
	Index   int
	Command interface{}
}
