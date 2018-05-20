package consts

const (
	JudgePD  = 0  //Pending
	JudgeRJ  = 1  //Running & judging
	JudgeCE  = 2  //Compile Error
	JudgeAC  = 3  //Accepted
	JudgeRE  = 4  //Runtime Error
	JudgeWA  = 5  //Wrong Answer
	JudgeTLE = 6  //Time Limit Exceeded
	JudgeMLE = 7  //Memory Limit Exceeded
	JudgeOLE = 8  //Output Limit Exceeded
	JudgePE  = 9  //Presentation Error
	JudgeNA  = 10 //System Error
	JudgeRPD = 11 //Rejudge Pending
)

const STD_MB = 1048576
const STD_T_LIM = 2
const STD_F_LIM = (STD_MB << 5)
const STD_M_LIM = (STD_MB << 7)

const (
	RLIMIT_CPU        = iota // CPU time in sec
	RLIMIT_FSIZE             // Maximum filesize
	RLIMIT_DATA              // max data size
	RLIMIT_STACK             // max stack size
	RLIMIT_CORE              // max core file size
	RLIMIT_RSS               // max resident set size
	RLIMIT_NPROC             // max number of processes
	RLIMIT_NOFILE            // max number of open files
	RLIMIT_MEMLOCK           // max locked-in-memory address space
	RLIMIT_AS                // address space limit
	RLIMIT_LOCKS             // maximum file locks held
	RLIMIT_SIGPENDING        // max number of pending signals
	RLIMIT_MSGQUEUE          // maximum bytes in POSIX mqueues
	RLIMIT_NICE              // max nice prio allowed to raise to
	RLIMIT_RTPRIO            // maximum realtime priority
	RLIMIT_RTTIME            // timeout for RT tasks in us
)
