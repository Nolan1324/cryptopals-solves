package slicex

func Repeat[T any](val T, len int) []T {
	s := make([]T, len)
	for i := range s {
		s[i] = val
	}
	return s
}
