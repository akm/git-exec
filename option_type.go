package main

import "strings"

type OptionType[T any] struct {
	envKeyPrefix string
	ShortName    string
	LongName     string
	HasValue     bool
	SetFunc      func(*T, string)
	WithoutEnv   bool
}

func newOptionType[T any](envKeyPrefix, shortName, longName string, hasValue bool, setFunc func(*T, string)) *OptionType[T] {
	return &OptionType[T]{
		envKeyPrefix: envKeyPrefix,
		ShortName:    shortName,
		LongName:     longName,
		HasValue:     hasValue,
		SetFunc:      setFunc,
	}
}

func (o *OptionType[T]) envKey() string {
	return o.envKeyPrefix + strings.ToUpper(strings.ReplaceAll(strings.TrimLeft(o.LongName, "-"), "-", "_"))
}

func (o *OptionType[T]) withoutEnv() *OptionType[T] {
	o.WithoutEnv = true
	return o
}
