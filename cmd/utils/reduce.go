package utils

func Reduce[T, M any](s []T, f func(M, T, int) M, initValue M) M {
	acc := initValue
	for i, v := range s {
		acc = f(acc, v, i)
	}
	return acc
}
