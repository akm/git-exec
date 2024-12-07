package main

import "strings"

type OptionType struct {
	envKeyPrefix string
	ShortName    string
	LongName     string
	HasValue     bool
	SetFunc      func(*Options, string)
	WithoutEnv   bool
}

func newOptionType(envKeyPrefix, shortName, longName string, hasValue bool, setFunc func(*Options, string)) *OptionType {
	return &OptionType{
		envKeyPrefix: envKeyPrefix,
		ShortName:    shortName,
		LongName:     longName,
		HasValue:     hasValue,
		SetFunc:      setFunc,
	}
}

func (o *OptionType) setValue(opts *Options, value string) {
	o.SetFunc(opts, value)
}

func (o *OptionType) envKey() string {
	return o.envKeyPrefix + strings.ToUpper(strings.ReplaceAll(strings.TrimLeft(o.LongName, "-"), "-", "_"))
}

func (o *OptionType) withoutEnv() *OptionType {
	o.WithoutEnv = true
	return o
}
