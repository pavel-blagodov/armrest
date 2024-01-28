package utils

func Find[T any](slice []T, predicate func(T) bool) *T {
	for _, item := range slice {
		if predicate(item) {
			return &item
		}
	}
	return nil
}
