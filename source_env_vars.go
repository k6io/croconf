package croconf

import (
	"encoding"
	"strconv"
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

func (sev *SourceEnvVars) GetName() string {
	return "environment variables" // TODO
}

func (sev *SourceEnvVars) From(name string) *envBinding {
	return &envBinding{
		source: sev,
		name:   name,
	}
}

type envBinding struct {
	source *SourceEnvVars
	name   string
}

func (eb *envBinding) GetSource() Source {
	return eb.source
}

func (eb *envBinding) BindStringValueTo(dest *string) func() error {
	return func() error {
		val, ok := eb.source.env[eb.name]
		if !ok {
			return ErrorMissing // TODO: better error message, e.g. 'field %s is not present in %s'?
		}
		*dest = val
		return nil
	}
}

func (eb *envBinding) BindInt64ValueTo(dest *int64) func() error {
	return func() error {
		val, ok := eb.source.env[eb.name]
		if !ok {
			return ErrorMissing // TODO: better error message, e.g. 'field %s is not present in %s'?
		}
		intVal, err := strconv.ParseInt(val, 10, 64) // TODO: use a custom function with better error message
		if err != nil {
			return err
		}
		*dest = intVal
		return nil
	}
}

func (eb *envBinding) BindTextBasedValueTo(dest encoding.TextUnmarshaler) func() error {
	return func() error {
		val, ok := eb.source.env[eb.name]
		if !ok {
			return ErrorMissing // TODO: better error message, e.g. 'field %s is not present in %s'?
		}

		return dest.UnmarshalText([]byte(val))
	}
}

func (eb *envBinding) BindValue(dest interface{}) func() error {
	// TODO: call the text binder insted to re-write the logic
	return func() error {
		val, ok := eb.source.env[eb.name]
		if !ok {
			return ErrorMissing // TODO: better error message, e.g. 'field %s is not present in %s'?
		}
		tdest, ok := dest.(encoding.TextUnmarshaler)
		if !ok {
			return ErrorMissing // TODO: better error message, e.g. 'field %s is not present in %s'?
		}
		return tdest.UnmarshalText([]byte(val))
	}
}

func parseEnvKeyValue(kv string) (string, string) {
	if idx := strings.IndexRune(kv, '='); idx != -1 {
		return kv[:idx], kv[idx+1:]
	}
	return kv, ""
}
