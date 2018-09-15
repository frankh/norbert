package check

type CheckRunner interface {
	Run(CheckInput) CheckResult
	Input() interface{}
	Validate() error
}

type CheckInput struct {
	Vars interface{}
}

type CheckResult struct {
	ResultCode CheckResultCode
	Error      error
}

type CheckResultCode int

const (
	CheckResultSuccess CheckResultCode = iota
	CheckResultFailure
	CheckResultError
)
