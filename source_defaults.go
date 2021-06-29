package croconf

// TODO: make more flexible with callbacks, so that besides defaut values, we
// can use these for custom wrappers as well?

type defaultStringValue string

func (dsv defaultStringValue) BindStringValueTo(dest *string) func() error {
	return func() error {
		*dest = string(dsv)
		return nil
	}
}

func (dsv defaultStringValue) GetSource() Source {
	return nil
}

func DefaultStringValue(s string) StringValueBinding {
	return defaultStringValue(s)
}

type defaultInt64Value int64

func (div defaultInt64Value) BindInt64ValueTo(dest *int64) func() error {
	return func() error {
		*dest = int64(div)
		return nil
	}
}
func (div defaultInt64Value) GetSource() Source {
	return nil
}

func DefaultInt64Value(i int64) Int64ValueBinding {
	return defaultInt64Value(i)
}

//TODO: sources for the rest of the types
