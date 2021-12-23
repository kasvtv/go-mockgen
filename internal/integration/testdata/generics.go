package testdata

type Set[T any] interface {
	Contains() bool
	Add(v T)
	Remove(v T)
}

type StringSetIntersector interface {
	// Intersect uses instantiated parametric types
	Intersect(s1, s2 Set[string])

	// Empty returns an instantiated parametric type
	Empty() Set[string]
}

type Map[K comparable, V any] interface {
	Contains(key K) bool
	Get(key K) (V, bool)
	Put(key K, value V) (oldValue V)
	Remove(key K)
}
