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
