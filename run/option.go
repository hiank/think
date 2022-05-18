package run

type Option[T any] interface {
	Apply(T)
}

type FuncOption[T any] func(T)

func (fo FuncOption[T]) Apply(v T) {
	fo(v)
}
