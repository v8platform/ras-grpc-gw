package pudgedb

import (
	"github.com/recoilme/pudge"
	_ "github.com/recoilme/pudge"
)

type Db struct {
	dir string
}

func (d Db) Close() error {
	return pudge.CloseAll()
}

func New(dir string) Db {
	return Db{
		dir: dir,
	}
}
