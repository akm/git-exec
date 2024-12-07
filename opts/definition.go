package opts

import "strings"

type Definition[T any] struct {
	envKeyPrefix string
	ShortName    string
	LongName     string
	HasValue     bool
	Setter       func(*T, string)
	getter       func(*T) string
	help         string
	withoutEnv   bool
}

func NewDefinition[T any](envKeyPrefix, shortName, longName string, hasValue bool, setFunc func(*T, string)) *Definition[T] {
	return &Definition[T]{
		envKeyPrefix: envKeyPrefix,
		ShortName:    shortName,
		LongName:     longName,
		HasValue:     hasValue,
		Setter:       setFunc,
	}
}

func (o *Definition[_]) EnvKey() string {
	return o.envKeyPrefix + strings.ToUpper(strings.ReplaceAll(strings.TrimLeft(o.LongName, "-"), "-", "_"))
}

func (o *Definition[T]) WithoutEnv() *Definition[T] {
	o.withoutEnv = true
	return o
}
func (o *Definition[_]) GetWithoutEnv() bool {
	return o.withoutEnv
}

func (o *Definition[T]) Getter(getter func(*T) string) *Definition[T] {
	o.getter = getter
	return o
}
func (o *Definition[T]) GetGetter() func(*T) string {
	return o.getter
}

func (o *Definition[T]) Help(help string) *Definition[T] {
	o.help = help
	return o
}
func (o *Definition[T]) GetHelp() string {
	return o.help
}
