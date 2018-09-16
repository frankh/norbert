// Code generated by "enumer -type=CheckResultCode -json"; DO NOT EDIT.

package check

import (
	"encoding/json"
	"fmt"
)

const _CheckResultCodeName = "CheckResultSuccessCheckResultFailureCheckResultError"

var _CheckResultCodeIndex = [...]uint8{0, 18, 36, 52}

func (i CheckResultCode) String() string {
	if i < 0 || i >= CheckResultCode(len(_CheckResultCodeIndex)-1) {
		return fmt.Sprintf("CheckResultCode(%d)", i)
	}
	return _CheckResultCodeName[_CheckResultCodeIndex[i]:_CheckResultCodeIndex[i+1]]
}

var _CheckResultCodeValues = []CheckResultCode{0, 1, 2}

var _CheckResultCodeNameToValueMap = map[string]CheckResultCode{
	_CheckResultCodeName[0:18]:  0,
	_CheckResultCodeName[18:36]: 1,
	_CheckResultCodeName[36:52]: 2,
}

// CheckResultCodeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func CheckResultCodeString(s string) (CheckResultCode, error) {
	if val, ok := _CheckResultCodeNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to CheckResultCode values", s)
}

// CheckResultCodeValues returns all values of the enum
func CheckResultCodeValues() []CheckResultCode {
	return _CheckResultCodeValues
}

// IsACheckResultCode returns "true" if the value is listed in the enum definition. "false" otherwise
func (i CheckResultCode) IsACheckResultCode() bool {
	for _, v := range _CheckResultCodeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for CheckResultCode
func (i CheckResultCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for CheckResultCode
func (i *CheckResultCode) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("CheckResultCode should be a string, got %s", data)
	}

	var err error
	*i, err = CheckResultCodeString(s)
	return err
}