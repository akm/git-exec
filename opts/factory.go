package opts

import "fmt"

type Factory[T any] struct {
	envKeyPrefix string

	boolToString func(bool) string
	quote        func(string) string
}

func NewFactory[T any](envKeyPrefix string) *Factory[T] {
	return &Factory[T]{
		envKeyPrefix: envKeyPrefix,
		boolToString: func(b bool) string {
			if b {
				return "true"
			} else {
				return "false"
			}
		},
		quote: func(s string) string { return fmt.Sprintf("%q", s) },
	}
}

func (f *Factory[T]) Bool(shortName, longName string, help string, getter func(*T) bool, setter func(*T)) *Definition[T] {
	var actualGetter func(o *T) string
	if getter != nil {
		actualGetter = func(o *T) string { return f.boolToString(getter(o)) }
	}
	return NewDefinition(f.envKeyPrefix, shortName, longName, false,
		func(o *T, v string) { setter(o) }).
		Getter(actualGetter).Help(help)
}

func (f *Factory[T]) String(shortName, longName string, help string, getter func(*T) string, setter func(*T, string)) *Definition[T] {
	var actualGetter func(o *T) string
	if getter != nil {
		actualGetter = func(o *T) string { return f.quote(getter(o)) }
	}
	return NewDefinition(f.envKeyPrefix, shortName, longName, true, setter).
		Getter(actualGetter).Help(help)
}
