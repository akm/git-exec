package opts

import "strings"

type OptionType[T any] struct {
	envKeyPrefix string
	ShortName    string
	LongName     string
	HasValue     bool
	SetFunc      func(*T, string)
	withoutEnv   bool
}

func NewOptionType[T any](envKeyPrefix, shortName, longName string, hasValue bool, setFunc func(*T, string)) *OptionType[T] {
	return &OptionType[T]{
		envKeyPrefix: envKeyPrefix,
		ShortName:    shortName,
		LongName:     longName,
		HasValue:     hasValue,
		SetFunc:      setFunc,
	}
}

func (o *OptionType[T]) EnvKey() string {
	return o.envKeyPrefix + strings.ToUpper(strings.ReplaceAll(strings.TrimLeft(o.LongName, "-"), "-", "_"))
}

func (o *OptionType[T]) WithoutEnv() *OptionType[T] {
	o.withoutEnv = true
	return o
}
func (o *OptionType[T]) GetWithoutEnv() bool {
	return o.withoutEnv
}
