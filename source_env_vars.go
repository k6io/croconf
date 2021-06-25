package croconf

import (
	"encoding"
<<<<<<< HEAD
	"fmt"
	"strconv"
=======
>>>>>>> 6e1bf46 ([feat] basic types extension)
	"strings"
)

type SourceEnvVars struct {
	env map[string]string
	// TODO
}

func NewSourceFromEnv(environ []string) *SourceEnvVars {
	env := make(map[string]string, len(environ))
	for _, kv := range environ {
		k, v := parseEnvKeyValue(kv)
		env[k] = v
	}
	return &SourceEnvVars{env: env}
}

func (sev *SourceEnvVars) Initialize() error {
	return nil // TODO? maybe prefix handling?
}

func (sev *SourceEnvVars) GetName() string {
	return "environment variables" // TODO
}

func (sev *SourceEnvVars) From(name string) *envBinding {
	return &envBinding{
		source: sev,
		name:   name,
		lookup: func() (string, error) {
			val, ok := sev.env[name]
			if !ok {
				return "", ErrorMissing // TODO: better error message, e.g. 'field %s is not present in %s'?
			}
			return val, nil
		},
	}
}

type envBinding struct {
	source Source
	name   string
	lookup func() (string, error)
}

func (eb *envBinding) GetSource() Source {
	return eb.source
}

func (eb *envBinding) BindStringValueTo(dest *string) func() error {
	return func() error {
		val, err := eb.lookup()
		if err != nil {
			return err
		}
		*dest = val
		return nil
	}
}

func (eb *envBinding) BindIntValue() func(bitSize int) (int64, error) {
	return func(bitSize int) (int64, error) {
<<<<<<< HEAD
		val, err := eb.lookup()
		if err != nil {
			return 0, err
=======
		val, ok := eb.source.env[eb.name]
		if !ok {
			return 0, NewBindNameMissingError("BindIntValue", eb.name)
>>>>>>> 6e1bf46 ([feat] basic types extension)
		}
		intVal, err := parseInt(val, 10, bitSize)
		if err != nil {
			return 0, err.withFuncName("BindIntValue")
		}
		return intVal, nil
	}
}

func (eb *envBinding) BindUintValue() func(bitSize int) (uint64, error) {
	return func(bitSize int) (uint64, error) {
		strVal, ok := eb.source.env[eb.name]
		if !ok {
			return 0, NewBindNameMissingError("BindUintValue", eb.name)
		}
		intVal, err := parseUint(strVal, 10, bitSize)
		if err != nil {
			return 0, err.withFuncName("BindUintValue")
		}
		return intVal, nil
	}
}

func (eb *envBinding) BindFloatValue() func(bitSize int) (float64, error) {
	return func(bitSize int) (float64, error) {
		strVal, ok := eb.source.env[eb.name]
		if !ok {
			return 0, NewBindNameMissingError("BindFloatValue", eb.name)
		}
		val, err := parseFloat(strVal, bitSize)
		if err != nil {
			return 0, err.withFuncName("BindFloatValue")
		}
		return val, nil
	}
}

func (eb *envBinding) BindTextBasedValueTo(dest encoding.TextUnmarshaler) func() error {
	return func() error {
		val, err := eb.lookup()
		if err != nil {
			return err
		}

		return dest.UnmarshalText([]byte(val))
	}
}

func (eb *envBinding) BindArray() func() (Array, error) {
	return func() (Array, error) {
		val, err := eb.lookup()
		if err != nil {
			return nil, err
		}

		arr := strings.Split(val, ",") // TODO: figure out how to make the delimiter configurable

		return &envVarArray{eb: eb, array: arr}, nil
	}
}

type envVarArray struct {
	eb    *envBinding
	array []string
}

func (eva *envVarArray) Len() int {
	return len(eva.array)
}

func (eva *envVarArray) Element(elNum int) LazySingleValueBinding {
	name := fmt.Sprintf("%s[%d]", eva.eb.name, elNum)
	return &envBinding{
		source: eva.eb.source,
		name:   name,
		lookup: func() (string, error) {
			if elNum >= len(eva.array) {
				return "", fmt.Errorf("tried to access invalid element %s, array only has %d elements", name, elNum)
			}
			return eva.array[elNum], nil
		},
	}
}

func parseEnvKeyValue(kv string) (string, string) {
	if idx := strings.IndexRune(kv, '='); idx != -1 {
		return kv[:idx], kv[idx+1:]
	}
	return kv, ""
}
