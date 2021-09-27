package pudgedb

import (
	"github.com/elastic/go-ucfg"
	"github.com/recoilme/pudge"
	"path/filepath"
	"sync"
)

type Table struct {
	*pudge.Db
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

func newFromConfig(config *Config) *Db {
	return &Db{
		dir:          config.Path,
		fileMode:     config.FileMode,
		dirMode:      config.DirMode,
		syncInterval: int(config.SyncInterval.Seconds()),
		storeMode:    config.StoreMode,
		tables:       make(map[string]*Table),
	}
}

func New(cfg *ucfg.Config) (*Db, error) {

	config, err := Unpack(cfg)
	if err != nil {
		return nil, err
	}

	return newFromConfig(config), nil
}
