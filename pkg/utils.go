package prtglaw

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

// Dump to JSON string
func Dumps(v interface{}, indent bool, settings *Settings) (string, error) {
	var data []byte
	var err error

	if indent {
		data, err = json.MarshalIndent(v, "", "    ")
	} else {
		data, err = json.Marshal(v)
	}

	if err != nil {
		settings.Logger.Warn().Err(err).Msg("Failed to marshal object")
		return "", err
	} else {
		return string(data), nil
	}
}

// Dump to JSON bytes
func Dumpb(v interface{}, settings *Settings) ([]byte, error) {
	data, err := Dumps(v, false, settings)
	return []byte(data), err
}

// Dums but it will return an empty if there is an error in marshalling, logs the error
func DumpsForce(v interface{}, indent bool, settings *Settings) string {
	data, err := Dumps(v, indent, settings)
	if err != nil {
		log.Error().Str("Name", "DumpsForce")
		return ""
	}
	return data
}

func DumpbForce(v interface{}, settings *Settings) []byte {
	return []byte(DumpsForce(v, false, settings))
}
