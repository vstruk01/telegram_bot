package utils

type Set map[string]struct{}

func (s Set) Check(key string) bool {
	_, ok := s[key]
	return ok
}

func (s Set) Set(key string) {
	s[key] = struct{}{}
}
