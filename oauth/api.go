package oauth

type IAuther interface {
	Auth(token string) (uid uint64, err error)
}
