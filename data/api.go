package dset

import "github.com/hiank/think/data/db"

const (
	TagUseMemory int8 = 1 << 0
	TagUseDisk   int8 = 1 << 1
)

type DBKey struct {
	Tag    int8
	Result string
}

//InTag check the given tag
func (hk *DBKey) InTag(want int8) bool {
	return (hk.Tag & want) == want
}

//IGamer gamer's information
type IGamer interface {
	GetUid() uint64
	GetToken() string
}

type BuildGamer func() IGamer

//IDataset data set
type IDataset interface {
	//GetGamer get gamer's information
	GetGamer(uid uint64) (IGamer, error)

	HGet(hkey *DBKey, fkey string) (db.IParser, error)
	HSet(hkey *DBKey, values ...interface{}) error
}
