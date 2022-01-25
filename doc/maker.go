package doc

type funcBytesMaker func([]byte) Doc

func (fbm funcBytesMaker) Make(v []byte) Doc {
	return fbm(v)
}
