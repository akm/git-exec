package opts

import "strings"

type Definition[T any] struct {
	envKeyPrefix string
	ShortName    string
	LongName     string
	HasValue     bool
	SetFunc      func(*T, string)
	withoutEnv   bool
}

func NewDefinition[T any](envKeyPrefix, shortName, longName string, hasValue bool, setFunc func(*T, string)) *Definition[T] {
	return &Definition[T]{
		envKeyPrefix: envKeyPrefix,
		ShortName:    shortName,
		LongName:     longName,
		HasValue:     hasValue,
		SetFunc:      setFunc,
	}
}

func (o *Definition[T]) EnvKey() string {
	return o.envKeyPrefix + strings.ToUpper(strings.ReplaceAll(strings.TrimLeft(o.LongName, "-"), "-", "_"))
}

func (o *Definition[T]) WithoutEnv() *Definition[T] {
	o.withoutEnv = true
	return o
}
func (o *Definition[T]) GetWithoutEnv() bool {
	return o.withoutEnv
}