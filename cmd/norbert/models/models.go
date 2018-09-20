package models

import (
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"io"
	"strconv"
	"time"

	"github.com/frankh/norbert/pkg/check"
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
	Cron   string `json:"cron"`

	Vars interface{} `json:"vars"`
}

type CheckResult struct {
	Id      string
	CheckId string

	ResultCode check.CheckResultCode
	ErrorMsg   string

	StartTime time.Time
	EndTime   time.Time
}

type Check struct {
	Name string `json:"name"`

	Service     string `json:"service"`
	CheckRunner string `json:"checkrunner"`
	Cron        string `json:"cron"`

	Severity Severity `json:"severity"`

	Vars interface{} `json:"vars"`
}

type CheckStatus int

const (
	_ CheckStatus = iota
	Ok
	Failed
	Initial
)

func (c *Check) Id() string {
	hash := fnv.New32()
	hash.Write([]byte(c.Name + c.Service))
	return hex.EncodeToString(hash.Sum(nil))
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

func (t *CheckStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	var err error
	*t, err = CheckStatusString(str)
	return err
}

func (t CheckStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(t.String()))
}
