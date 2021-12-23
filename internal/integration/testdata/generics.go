package testdata

type Restricted[T uint8 | ~uint16] interface {
	Accept(v T) bool
}

type Set[T any] interface {
	Contains() bool
	Add(v T)
	Remove(v T)
}

type NumericSets interface {
	//
	// TODO - stop from generating this
	//

	int8 | int16 | int64
}

type Map[K comparable, V any] interface {
	Contains(key K) bool
	Get(key K) (V, bool)
	Put(key K, value V) (oldValue V)
	Remove(key K)
}

type StringSetIntersector interface {
	// Intersect uses instantiated parametric types
	Intersect(s1, s2 Set[string])

	// Empty returns an instantiated parametric type
	Empty() Set[string]
}

type unexportedGeneric[T any] interface {
	Serialize() T
}
