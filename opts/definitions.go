package opts

type Definitions[T any] []*Definition[T]

func BuildKeyMap[T any](defs Definitions[T]) map[string]*Definition[T] {
	m := make(map[string]*Definition[T], len(defs))
	for _, def := range defs {
		m[def.LongName] = def
		m[def.ShortName] = def
	}
	return m
}
