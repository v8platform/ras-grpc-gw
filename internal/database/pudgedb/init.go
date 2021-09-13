package pudgedb

import (
	"github.com/recoilme/pudge"
	"path/filepath"
	"sync"
)

type Table struct {
	*pudge.Db
}

func (t *Table) FindOne() {

}

type Db struct {
	dir          string
	fileMode     int
	dirMode      int
	syncInterval int
	storeMode    int

	tables map[string]*Table
	mu     sync.Mutex
}

func (d *Db) Close() error {

	d.tables = make(map[string]*Table)

	return pudge.CloseAll()
}

func (d *Db) GetPath(path string) string {
	return filepath.Join(d.dir, path)
}

func (d *Db) Table(name string) (*Table, error) {

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

	table := &Table{open}

	d.tables[name] = table

	return table, nil

}

func New(dir string) *Db {
	return &Db{
		dir:          dir,
		fileMode:     0666,
		dirMode:      0777,
		syncInterval: 60,
		storeMode:    0,
		tables:       make(map[string]*Table),
	}
}
