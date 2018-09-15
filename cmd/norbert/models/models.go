package models

import (
	"fmt"
	"io"
	"strconv"
)

type Severity int

const (
	_ Severity = iota
	Warning
	Error
	Critical
)

type CheckRunner struct {
	Name   string `json:"name"`
	Plugin string `json:"plugin"`

	Vars map[string]interface{} `json:"vars"`
}

type Service struct {
	Name string

	Source string
}

type Check struct {
	CheckRunner string

	Service  Service
	Severity Severity

	Source string
}

func (t *Severity) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	var err error
	*t, err = SeverityString(str)
	return err
}

func (t Severity) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(t.String()))
}
