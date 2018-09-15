package check

type CheckRunner interface {
	Run(CheckInput) CheckResult

	// Return a struct containing the default input variables
	// required for the runner. The "vars" yaml field will be
	// deserialized into this struct.
	Vars() interface{}

	// Called at load time, to validate that e.g. required
	// environment variables are set, or other runtime checks.
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
