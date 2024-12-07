package opts

import "fmt"

func HelpItemsAndEnvVarMappings[T any](defaultOptions *T, defs []*Definition[T]) ([]string, []string) {
	optionItems := make([]string, len(defs))
	envVarItems := []string{}

	maxLongNameLength := 0
	for _, opt := range defs {
		if maxLongNameLength < len(opt.LongName) {
			maxLongNameLength = len(opt.LongName)
		}
	}
	indent := "  "

	longNameFormat := "%-" + fmt.Sprintf("%ds", maxLongNameLength)
	for i, opt := range defs {
		var item string
		if opt.ShortName == "" {
			item = fmt.Sprintf("%s    "+longNameFormat, indent, opt.LongName)
		} else {
			item = fmt.Sprintf("%s%s, "+longNameFormat, indent, opt.ShortName, opt.LongName)
		}
		item += " " + opt.GetHelp()
		if getter := opt.GetGetter(); getter != nil {
			item += fmt.Sprintf(" (default: %s)", getter(defaultOptions))
		}
		optionItems[i] = item

		if !opt.GetWithoutEnv() {
			envVarItems = append(envVarItems, fmt.Sprintf(longNameFormat+" %s", opt.LongName, opt.EnvKey()))
		}
	}
	return optionItems, envVarItems
}
