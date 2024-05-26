package prtglaw

import (
	"encoding/json"
)

// Dump to JSON string
func Dumps(v interface{}, indent bool) (string, error) {
	var data []byte
	var err error

	if indent {
		data, err = json.MarshalIndent(v, "", "    ")
	} else {
		data, err = json.Marshal(v)
	}

	if err != nil {
		return "", err
	} else {
		return string(data), nil
	}
}

// Dump to JSON bytes
func Dumpb(v interface{}) ([]byte, error) {
	data, err := Dumps(v, false)
	return []byte(data), err
}
