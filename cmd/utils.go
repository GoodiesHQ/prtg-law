package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

// to be run when setting a target value to either a provided fallback or the zero-value
func setDefaultFunc[T any](target *T, fallback *T) func() {
	return func() {
		if fallback != nil && target != nil {
			*target = *fallback
		}
	}
}

// to be run when getting a target value as either a provided fallback or the zero-value
func getDefaultFunc[T any](fallback *T) func() T {
	return func() T {
		var zero T
		if fallback == nil {
			return zero
		}
		return *fallback
	}
}

func Parse[T any](value string, fallback *T) (T, error) {
	var zero T

	invalid := func(typeName string) error {
		return fmt.Errorf("error parsing value '%s' as type %s", value, typeName)
	}

	getDefault := getDefaultFunc(fallback)

	switch any(zero).(type) {
	case int:
		if result, err := strconv.Atoi(value); err != nil {
			return getDefault(), invalid("int")
		} else {
			return any(result).(T), nil
		}
	case uint16:
		if result, err := strconv.ParseUint(value, 10, 16); err != nil {
			return getDefault(), invalid("uint16")
		} else {
			return any(result).(T), nil
		}
	case float32:
		if result, err := strconv.ParseFloat(value, 32); err != nil {
			return zero, invalid("float32")
		} else {
			return any(result).(T), nil
		}
	case float64:
		if result, err := strconv.ParseFloat(value, 64); err != nil {
			return zero, invalid("float64")
		} else {
			return any(result).(T), nil
		}
	case string:
		return any(value).(T), nil
	default:
		return zero, fmt.Errorf("invalid type")
	}
}

func EnvLookup[T any](target *T, key string, fallback *T) error {
	setDefault := setDefaultFunc(target, fallback)

	if value, found := os.LookupEnv(key); !found {
		setDefault()
		return fmt.Errorf("environment variable '%s' not found", key)
	} else {
		parsed, err := Parse(value, fallback)
		if err != nil {
			setDefault()
			return fmt.Errorf("unable to parse environment variable '%v'", target)
		}
		*target = parsed
		return nil
	}
}

func EnvLookupMust[T any](target *T, key string) {
	err := EnvLookup(target, key, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Environment Variable Failure")
		// log.Fatalf("%v", err)
	}
}
