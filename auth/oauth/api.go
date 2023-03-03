package oauth

type Auther interface {
	Auth(token string) (uid uint64, err error)
}
