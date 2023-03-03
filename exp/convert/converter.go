package convert

type Converter[T1, T2 any] interface {
	Convert(T1) T2
}

type ConverterFunc[T1, T2 any] func(T1) T2

func (cf ConverterFunc[T1, T2]) Convert(v T1) T2 {
	return cf(v)
}
