// Code generated by "enumer -type=CheckStatus -json -sql"; DO NOT EDIT.

package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

const _CheckStatusName = "OkFailedInitial"

var _CheckStatusIndex = [...]uint8{0, 2, 8, 15}

func (i CheckStatus) String() string {
	i -= 1
	if i < 0 || i >= CheckStatus(len(_CheckStatusIndex)-1) {
		return fmt.Sprintf("CheckStatus(%d)", i+1)
	}
	return _CheckStatusName[_CheckStatusIndex[i]:_CheckStatusIndex[i+1]]
}

var _CheckStatusValues = []CheckStatus{1, 2, 3}

var _CheckStatusNameToValueMap = map[string]CheckStatus{
	_CheckStatusName[0:2]:  1,
	_CheckStatusName[2:8]:  2,
	_CheckStatusName[8:15]: 3,
}

// CheckStatusString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func CheckStatusString(s string) (CheckStatus, error) {
	if val, ok := _CheckStatusNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to CheckStatus values", s)
}

// CheckStatusValues returns all values of the enum
func CheckStatusValues() []CheckStatus {
	return _CheckStatusValues
}

// IsACheckStatus returns "true" if the value is listed in the enum definition. "false" otherwise
func (i CheckStatus) IsACheckStatus() bool {
	for _, v := range _CheckStatusValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for CheckStatus
func (i CheckStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for CheckStatus
func (i *CheckStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("CheckStatus should be a string, got %s", data)
	}

	var err error
	*i, err = CheckStatusString(s)
	return err
}

func (i CheckStatus) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *CheckStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		bytes, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("value is not a byte slice")
		}

		str = string(bytes[:])
	}

	val, err := CheckStatusString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}
