package util

type CompareFunc[T any] func(a, b T) bool

func CloneSlice[T any](origin []T) []T {
	clone := make([]T, len(origin))
	copy(clone, origin)
	return clone
}

func RemoveFromSlice[T any](origin []T, target T, compare CompareFunc[T]) []T {
	var clone []T
	for _, item := range origin {
		if !compare(item, target) {
			clone = append(clone, item)
		}
	}
	return clone
}
