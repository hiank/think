package store

type easy[T ~string] struct {
	td Dictionary[T]
}

func (ed *easy[T]) Scan(k string, out any) (bool, error) {
	return ed.td.Scan(T(k), out)
}

func (ed *easy[T]) Set(k string, v any) error {
	return ed.td.Set(T(k), v)
}

func (ed *easy[T]) Del(k string, out ...any) error {
	return ed.td.Del(T(k), out...)
}

func (ed *easy[T]) Close() error {
	return ed.td.Close()
}

func ConvertoEasy[T ~string](td Dictionary[T]) EasyDictionary {
	return &easy[T]{td}
}
