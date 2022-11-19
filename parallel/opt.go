package parallel

import "runtime"

func InitialOption(ops []OptionOp) *Option {
	opt := NewOption()
	for _, op := range ops {
		op(opt)
	}
	return opt
}

func NewOption() *Option {
	return &Option{workerNumber: runtime.GOMAXPROCS(0)}
}

type Option struct {
	workerNumber int

	// --- stream options start
	receiveDataExplicit bool
	ignoreResult        bool
	hashWorker          bool
}
type OptionOp func(opt *Option)

func WorkerNumber(num int) OptionOp {
	return func(opt *Option) {
		opt.workerNumber = num
	}
}

// IgnoreResult used in stream parallel, if you do not care about result, use it.
func IgnoreResult() OptionOp {
	return func(opt *Option) {
		opt.ignoreResult = true
	}
}

// ReceiveDataFromChan used in stream parallel, it means that you need receive data from chan not memory cached
func ReceiveDataFromChan() OptionOp {
	return func(opt *Option) {
		opt.receiveDataExplicit = true
	}
}

type KeyOwner interface {
	GetKey() string
}
