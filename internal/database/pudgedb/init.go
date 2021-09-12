package pudgedb

import (
	"github.com/recoilme/pudge"
	_ "github.com/recoilme/pudge"
	"path/filepath"
	"sync"
)

type Db struct {
	dir          string
	fileMode     int
	dirMode      int
	syncInterval int
	storeMode    int

	tables map[string]*pudge.Db
	mu     sync.Mutex
}

func (d *Db) Close() error {

	d.tables = make(map[string]*pudge.Db)

	return pudge.CloseAll()
}

func (d *Db) GetPath(path string) string {
	return filepath.Join(d.dir, path)
}

func (d *Db) Table(name string) (*pudge.Db, error) {

	if table, ok := d.tables[name]; ok {
		return table, nil
	}

	open, err := pudge.Open(d.GetPath(name), &pudge.Config{
		FileMode:     d.fileMode,
		DirMode:      d.dirMode,
		SyncInterval: d.syncInterval,
		StoreMode:    d.storeMode,
	})

	if err != nil {
		return nil, err
	}
	d.mu.Lock()
	defer d.mu.Unlock()

	d.tables[name] = open

	return open, nil

}

func New(dir string) *Db {
	return &Db{
		dir:          dir,
		fileMode:     0666,
		dirMode:      0777,
		syncInterval: 60,
		storeMode:    0,
		tables:       make(map[string]*pudge.Db),
	}
}
