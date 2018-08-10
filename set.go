package evidence

var exists = struct{}{}

type stringSet map[string]struct{}

func (s stringSet) Contains(value string) bool {
	_, c := s[value]
	return c
}
