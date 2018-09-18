package check

import (
	"fmt"
	"io"
	"strconv"
)

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
	Success CheckResultCode = iota
	Failure
	Error
)

func (t *CheckResultCode) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	var err error
	*t, err = CheckResultCodeString(str)
	return err
}

func (t CheckResultCode) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(t.String()))
}
