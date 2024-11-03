package main

func splitStringsInto2(args []string, fn func(string) bool) ([]string, []string) {
	var a, b []string
	for _, arg := range args {
		if fn(arg) {
			a = append(a, arg)
		} else {
			b = append(b, arg)
		}
	}
	return a, b
}
