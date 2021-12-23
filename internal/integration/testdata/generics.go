package testdata

type Set[T any] interface {
	Contains() bool
	Add(v T)
	Remove(v T)
}
