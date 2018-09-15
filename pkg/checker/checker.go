package checker

type Checker interface {
	Run(CheckInput) CheckResult
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
