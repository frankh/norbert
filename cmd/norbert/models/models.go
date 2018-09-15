package models

import (
	"fmt"
	"io"
	"strconv"
)

type Severity int

const (
	_ Severity = iota
	Info
	Error
	Critical
)

type CheckRunner struct {
	Name   string `json:"name"`
	Plugin string `json:"plugin"`

	Vars interface{} `json:"vars"`
}

type Check struct {
	Service     string `json:"service"`
	CheckRunner string `json:"checkrunner"`

	Severity Severity `json:"severity"`

	Vars interface{} `json:"vars"`
}

type Service struct {
	Name string `json:"name"`
	Url  string `json:"url"`

	Vars interface{} `json:"vars"`
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
