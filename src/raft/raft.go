package raft

import (
	"fmt"
	"github.com/hhr12138/Konata/src/consts"
	"github.com/hhr12138/Konata/src/entity"
	"github.com/hhr12138/Konata/src/labrpc"
	"github.com/hhr12138/Konata/src/utils"
	"sync"
	"sync/atomic"
	"time"
)

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

// import "bytes"
// import "../labgob"



//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// in Lab 3 you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh; at that point you can add fields to
// ApplyMsg, but set CommandValid to false for these other uses.
//
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int
}

//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]
	dead      int32               // set by Kill()

	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.
	term    int          // 任期
	status consts.Status // 状态
	lockMap entity.Locks // 所有锁控制
	logs []*entity.Log // 日志
	voteFor int // 投票给谁
	electionTime atomic.Value // 选举超时
	applyCh chan ApplyMsg
	nextIndex []int
	matchIndex []int
	commitIndex int32
	lastApplied int32

}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {
	var term int
	var isleader bool
	// Your code here (2A).
	locks := utils.GetLockMap(consts.STATUS, consts.TERM)
	rf.lockMap.Lock(locks)
	defer rf.lockMap.Unlock(locks)
	term = rf.getTerm()
	isleader = rf.getStatus() == consts.LEADER
	return term, isleader
}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here (2C).
	// Example:
	// w := new(bytes.Buffer)
	// e := labgob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// data := w.Bytes()
	// rf.persister.SaveRaftState(data)
}


//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// Your code here (2C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := labgob.NewDecoder(r)
	// var xxx
	// var yyy
	// if d.Decode(&xxx) != nil ||
	//    d.Decode(&yyy) != nil {
	//   error...
	// } else {
	//   rf.xxx = xxx
	//   rf.yyy = yyy
	// }
}




//
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//
type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	term int
	candidateId int
	lastLogIndex int
	lastLogTerm int
}

//
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	// Your data here (2A).
	term int
	voteGranted bool
}

//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (2A, 2B).
}

type RequestAppendArgs struct {
	term int
	leaderId int
	prevLogIndex int
	prevLogTerm int
	entries []*entity.Log
	leaderCommit int32
}

type RequestAppendReply struct {
	term int
	success bool
}

func (rf *Raft) RequestAppendEntries(args *RequestAppendArgs,reply *RequestAppendReply, slaveId int) {
	locks := utils.GetLockMap(consts.LOG, consts.STATUS, consts.TERM, consts.NEXT_INDEX)
	rf.lockMap.Lock(locks)
	var(
		term = rf.getTerm()
		status = rf.getStatus()
		logStartIdx = rf.nextIndex[slaveId]
		logLen = rf.getLogLen()
		prevLogIndex = -1
		prevLogTerm = -1
	)
	if logStartIdx > logLen{
		utils.Printf(consts.ERROR,rf.getEndName(rf.me),term,logLen,"[RequestAppendEntries] 期望的nextIndex>主节点最大偏移量,slaveId=%v",slaveId)
		logStartIdx = logLen
	}
	appendLogs := rf.logs[logStartIdx:logLen]
	if logStartIdx > 0{
		prevLogIndex = rf.logs[logStartIdx-1].Index
		prevLogTerm = rf.logs[logStartIdx-1].Term
	}
	rf.lockMap.Unlock(locks)
	if rf.killed() || status != consts.LEADER{
		return
	}
	// appendLogs获取完毕,记得判空,然后接着写
	args.entries = appendLogs
	args.term = term
	args.leaderCommit = rf.getCommitIndex()
	args.leaderId = rf.me
	args.prevLogIndex = prevLogIndex
	args.prevLogTerm = prevLogTerm

	// 发送RPC
	success := rf.sendAppendEntries(slaveId, args, reply)

	// 失败直接返回
	if !success {
		return
	}
	// RPC响应检测
	success = rf.overdueRspCheck(reply.term, term)
	if !success{
		return
	}
}

func (rf *Raft) AppendEntries(args *RequestAppendArgs, reply *RequestAppendReply) {
	term, _, _, _ := rf.getPrimeInfoInLock()
	// 初始化返回值
	reply.term = term
	reply.success = false
	// 过期请求,直接忽略
	if !rf.overdueReqCheck(args.term,term){
		return
	}
	// 更新超时时间
	// 日志判断&追加 2B
	// 修改commitIndex
	reply.success = true
	return
}

//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

// 发送appendEntries
func (rf *Raft) sendAppendEntries(server int, args *RequestAppendArgs, reply *RequestAppendReply) bool{
	ok := rf.peers[server].Call("Raft.AppendEntries",args,reply)
	return ok
}


//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election. even if the Raft instance has been killed,
// this function should return gracefully.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true

	// Your code here (2B).


	return index, term, isLeader
}

//
// the tester doesn't halt goroutines created by Raft after each test,
// but it does call the Kill() method. your code can use killed() to
// check whether Kill() has been called. the use of atomic avoids the
// need for a lock.
//
// the issue is that long-running goroutines use memory and may chew
// up CPU time, perhaps causing later tests to fail and generating
// confusing debug output. any goroutine with a long-running loop
// should call killed() to check whether it should stop.
//
func (rf *Raft) Kill() {
	atomic.StoreInt32(&rf.dead, 1)
	// Your code here, if desired.
}

func (rf *Raft) killed() bool {
	z := atomic.LoadInt32(&rf.dead)
	return z == 1
}

//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me

	// Your initialization code here (2A, 2B, 2C).
	rf.lockMap = make(map[consts.LockName]*sync.Mutex,len(consts.LockOrder))
	rf.nextIndex = make([]int,len(peers))
	rf.matchIndex = make([]int,len(peers))
	rf.logs = make([]*entity.Log,0)
	rf.electionTime.Store(time.Now())

	rf.lockMap.InitLocks()
	rf.applyCh = applyCh

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	// 开启心跳检测
	rf.heartBeat()
	return rf
}

// 通用操作
// 获取当前任期,偏移量,状态和自身编号等原信息, 不要提供无锁的getPrimeInfoInLock方法
func (rf *Raft) getPrimeInfoInLock() (term int, status consts.Status, offset,me int, ){
	locks := utils.GetLockMap(consts.TERM, consts.LOG, consts.STATUS)
	rf.lockMap.Lock(locks)
	defer rf.lockMap.Unlock(locks)
	term = rf.getTerm()
	status = rf.getStatus()
	offset = rf.getLogLen()
	me = rf.me
	return
}

// 过期RPC请求检查
func (rf *Raft) overdueReqCheck(argTerm, nowTerm int) bool{
	if argTerm < nowTerm {
		return false
	} else if argTerm > nowTerm{
		rf.updateTermByCASInLock(argTerm,nowTerm)
	}
	return true
}

// 过期RPC响应检查
func (rf *Raft) overdueRspCheck(rspTerm, nowTerm int) bool{
	if rspTerm < nowTerm{
		utils.Printf(consts.ERROR,rf.getEndName(rf.me),nowTerm,consts.NULL_OFFSET,"[overdueRspCheck] 响应任期小于预期任期,理论不可能,rspTerm=%v",rspTerm)
	} else if rspTerm > nowTerm{
		// 对方任期更高，回退为follower并修改任期
		locks := utils.GetLockMap(consts.TERM, consts.STATUS)
		utils.Printf(consts.INFO,rf.getEndName(rf.me),nowTerm,consts.NULL_OFFSET,"[overdueRspCheck] 收到高任期响应,回退为follower,修改任期为:%v",rspTerm)
		rf.lockMap.Lock(locks)
		defer rf.lockMap.Unlock(locks)
		if nowTerm != rf.term {
			utils.Printf(consts.WARN,rf.getEndName(rf.me),rf.term,consts.NULL_OFFSET,"[overdueRspCheck] 任期改变,忽略本次响应")
			return false
		}
		if !rf.updateStatus(consts.FOLLOWER){
			utils.Printf(consts.ERROR,rf.getEndName(rf.me),nowTerm,consts.NULL_OFFSET,"[overdueRspCheck] 非法的状态转换,当前状态为%v,目标状态为%v",rf.getStatus(),consts.FOLLOWER)
			return false
		}
		rf.updateTerm(rspTerm)
		return false
	}
	return true
}

// 任期相关
// 获取任期
func (rf *Raft) getTerm() int{
	return rf.term
}

// 上锁获取任期
func (rf *Raft) getTermInLock() int{
	locks := utils.GetLockMap(consts.TERM)
	rf.lockMap.Lock(locks)
	defer rf.lockMap.Unlock(locks)
	return rf.getTerm()
}

// 上锁修改任期和voteFor
func (rf *Raft) updateTermByCASInLock(target int, old int) bool{
	locks := utils.GetLockMap(consts.TERM,consts.LOG)
	rf.lockMap.Lock(locks)
	//var(
		//offset = len(rf.logs)
		//term = rf.getTerm()
	//)
	defer func() {
		rf.lockMap.Unlock(locks)
		//utils.ShowUnlockLog(consts.INFO,rf.getEndName(rf.me),term,offset,"[modTermByCASInLock]",locks)
	}()

	//utils.ShowLockLog(consts.INFO,rf.getEndName(rf.me),term,offset,"[modTermByCASInLock]",locks)
	return rf.updateTermByCAS(target,old)
}

func (rf *Raft) updateTermByCAS(target int , old int) bool{
	if rf.term == old{
		rf.updateTerm(target)
		return true
	}
	return false
}

func (rf *Raft) updateTerm(target int){
	rf.term = target
	rf.voteFor = consts.NULL_CAN
}

// 日志相关
func (rf *Raft) getLogLenInLock() int{
	var(
		locks = utils.GetLockMap(consts.LOG)
	)
	rf.lockMap.Lock(locks)
	defer rf.lockMap.Unlock(locks)
	return rf.getLogLen()
}

func (rf *Raft) getLogLen() int{
	return len(rf.logs)
}
//status相关
func (rf *Raft) getStatus() consts.Status{
	return rf.status
}

//状态转换
func (rf *Raft) updateStatusByCASInLock(targetStatus consts.Status, oldTerm int) bool{
	locks := utils.GetLockMap(consts.STATUS,consts.TERM)
	rf.lockMap.Lock(locks)
	defer rf.lockMap.Unlock(locks)
	if rf.getTerm() == oldTerm{
		return rf.updateStatus(targetStatus)
	}
	return false
}


func (rf *Raft) updateStatus(targetStatus consts.Status) bool{
	status := rf.getStatus()
	if utils.CanTran(status,targetStatus){
		rf.status = targetStatus
		return true
	}
	return false
}

// peers相关
func (rf *Raft) getEndName(idx int) string {
	if idx > len(rf.peers){
		return "undefined"
	}
	return fmt.Sprintf("raft%v",idx)
}

// 选举相关

// 心跳检测
func (rf *Raft) heartBeat() {
	for {
		if rf.killed() {
			return
		}
		// 不是leader才检测
		term, status, _, _ := rf.getPrimeInfoInLock()
		if status != consts.LEADER{
			// 检测不通过开始发起选举
			if !rf.overtimeCheck() {
				rf.startVote(term)
			}
		}
		time.Sleep(time.Millisecond*consts.HEARTBEAT_TIME)
	}
}

// 发起选举
func (rf *Raft) startVote(oldTerm int){
	var (
		term int
		status consts.Status
		locks = utils.GetLockMap(consts.STATUS, consts.TERM)
		successVoteCnt int32 = 0
		failVoteCnt int32 = 0
		slaveCnt = len(rf.peers)
		targetCnt int32 = int32(slaveCnt)/2+1
		timeout time.Time
	)
	rf.lockMap.Lock(locks)
	status = rf.getStatus()
	term = rf.getTerm()
	// 发生状态改变,返回
	if oldTerm != term{
		return
	}
	if !rf.updateStatus(consts.CANDIDATE) {
		utils.Printf(consts.ERROR,rf.getEndName(rf.me),term,consts.NULL_OFFSET,"[startVote] 非法的状态转换,当前状态为%v,目标状态为%v",status,consts.CANDIDATE)
		rf.lockMap.Unlock(locks)
		return
	}
	// 增加任期
	rf.updateTerm(oldTerm+1)
	// 为自己投票
	rf.updateVoteFor(rf.me)
	// 更新选举超时
	timeout = rf.updateElectionTime()
	rf.lockMap.Unlock(locks)
	// 并行开始请求选票
	// 选票计数
	for i := 0; i < slaveCnt; i++ {
		go func() {
			var (
				args = &RequestVoteArgs{}
				reply = &RequestVoteReply{}
			)
			rf.RequestVote(args,reply)
			// 收到选票+1
			if reply.voteGranted{
				atomic.AddInt32(&successVoteCnt,1)
				return
			}
			// 过期任期检测
			rf.overdueRspCheck(reply.term,term)
			atomic.AddInt32(&failVoteCnt,1)
		}()
	}
	for {
		locks = utils.GetLockMap(consts.TERM)
		rf.lockMap.Lock(locks)
		// 当前rf死亡/任期改变/超时/出现结果就立即返回
		if rf.killed() || rf.getTerm() != term || time.Now().Before(timeout) || atomic.LoadInt32(&successVoteCnt) >= targetCnt || atomic.LoadInt32(&failVoteCnt) >= targetCnt{
			rf.lockMap.Unlock(locks)
			break
		}
		rf.lockMap.Unlock(locks)
		time.Sleep(consts.WAIT_VOTE_TIME*time.Millisecond)
	}
	if atomic.LoadInt32(&successVoteCnt) < targetCnt {
		return
	}
	// 成功当选, 晋升为leader并立即发送心跳
	success := rf.promoteLeader(term)
	if success{
		rf.initLeader()
	} else {
		// 回退为follower
		rf.rollbackFollower(term)
	}
}

// 尝试晋升为leader
func (rf *Raft) promoteLeader(oldTerm int) bool{
	locks := utils.GetLockMap(consts.TERM, consts.STATUS)
	rf.lockMap.Lock(locks)
	defer rf.lockMap.Unlock(locks)
	term := rf.getTerm()
	utils.Printf(consts.INFO,rf.getEndName(rf.me),term,consts.NULL_OFFSET,"[promoteLeader] 开始晋升leader")
	success := utils.OperationByCAS(term,oldTerm, func() bool{
		if !rf.updateStatus(consts.LEADER) {
			return false
		}
		return true
	})
	if !success{
		utils.Printf(consts.WARN,rf.getEndName(rf.me),term,consts.NULL_OFFSET,"[promoteLeader] 晋升失败,当前状态为%v",rf.status.String())
	}
	return success
}

// 回退为follower
func (rf *Raft) rollbackFollower(oldTerm int) bool{
	locks := utils.GetLockMap(consts.TERM,consts.STATUS)
	rf.lockMap.Lock(locks)
	defer rf.lockMap.Unlock(locks)
	term := rf.getTerm()
	utils.Printf(consts.INFO,rf.getEndName(rf.me),term,consts.NULL_OFFSET,"[rollbackFollower] 选举失败,回退为follower")
	success := utils.OperationByCAS(term,oldTerm, func() bool{
		if !rf.updateStatus(consts.FOLLOWER) {
			return false
		}
		return true
	})
	if !success{
		utils.Printf(consts.WARN,rf.getEndName(rf.me),term,consts.NULL_OFFSET,"[rollbackFollower] 回退失败,当前状态为%v",rf.status.String())
	}
	return success
}

// leader初始化
func (rf *Raft) initLeader(){
	// 立即发送心跳并定期发送

}

func (rf *Raft) getElectionTime() time.Time{
	t := rf.electionTime.Load().(time.Time)
	return t
}

// 更新选举超时
func (rf *Raft) updateElectionTime() time.Time{
	now := time.Now()
	now.Add(time.Duration(consts.ELECTION_TIME)*time.Millisecond)
	rf.electionTime.Store(now)
	return now
}

// 超时检测, true通过检测,false不通过
func (rf *Raft) overtimeCheck() bool{
	t := rf.getElectionTime()
	return time.Now().Before(t)
}

func (rf *Raft) updateVoteFor(idx int){
	rf.voteFor = idx
}

// commitIndex
func (rf *Raft) getCommitIndex() int32{
	return atomic.LoadInt32(&rf.commitIndex)
}
