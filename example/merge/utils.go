package merge

func rangeTwodimensional[T any](s [][]T, f func(ia, ib int, v T)) {
	for ia, arr := range s {
		for ib, v := range arr {
			f(ia, ib, v)
		}
	}
}
