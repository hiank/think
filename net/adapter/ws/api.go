package ws

type IStorage interface {
	GetUidByToken(token string) (uint64, bool)
}
