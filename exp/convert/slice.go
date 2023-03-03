package convert

func Slice[T1, T2 any](arr []T1, converter Converter[T1, T2]) (out []T2) {
	out = make([]T2, len(arr))
	for i, v := range arr {
		out[i] = converter.Convert(v)
	}
	return
}
