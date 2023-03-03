package easy

type Converter[T1, T2 any] interface {
	Convert(T1) T2
}

func TwoDimensinal[T1, T2 any](arr [][]T1, converter Converter[T1, T2]) (out [][]T2) {
	out = make([][]T2, len(arr))
	for i, v := range arr {
		out[i] = make([]T2, len(v))
		for ti, tv := range v {
			out[i][ti] = converter.Convert(tv)
		}
	}
	return
}
