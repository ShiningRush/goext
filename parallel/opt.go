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

	// stream option
	ignoreResult bool
}
type OptionOp func(opt *Option)

func WorkerNumber(num int) OptionOp {
	return func(opt *Option) {
		opt.workerNumber = num
	}
}

func IgnoreResult() OptionOp {
	return func(opt *Option) {
		opt.ignoreResult = true
	}
}
