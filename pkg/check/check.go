package check

type CheckRunner interface {
	Run(CheckInput) CheckResult
	Input() interface{}
}

type CheckInput struct {
	Input interface{}
}

type CheckResult struct {
	ResultCode CheckResultCode
}

type CheckResultCode int

const (
	CheckResultSuccess CheckResultCode = iota
	CheckResultFailure
	CheckResultError
)
