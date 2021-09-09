package pudgedb

import (
	"github.com/recoilme/pudge"
	_ "github.com/recoilme/pudge"
	"path/filepath"
)

type Db struct {
	dir string
}

func (d Db) Close() error {
	return pudge.CloseAll()
}

func (d Db) GetPath(path string) string {
	return filepath.Join(d.dir, path)
}

func New(dir string) Db {
	return Db{
		dir: dir,
	}
}
