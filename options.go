package main

type OptionType struct {
	ShortName string
	LongName  string
	HasValue  bool
}

func newOptionType(shortName, longName string, hasValue bool) *OptionType {
	return &OptionType{
		ShortName: shortName,
		LongName:  longName,
		HasValue:  hasValue,
	}
}

var (
	optHelp      = newOptionType("-h", "--help", false)
	optVersion   = newOptionType("-v", "--version", false)
	optDirectory = newOptionType("-C", "--directory", true)
)

var optionTypes = []*OptionType{
	optHelp,
	optVersion,
	optDirectory,
}

var optionKeyMap = func() map[string]*OptionType {
	m := map[string]*OptionType{}
	for _, opt := range optionTypes {
		m[opt.ShortName] = opt
		m[opt.LongName] = opt
	}
	return m
}()

type Option struct {
	Type  *OptionType
	Value string
}

type Options []*Option
