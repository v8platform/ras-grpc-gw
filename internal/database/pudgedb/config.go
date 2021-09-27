package pudgedb

import (
	"github.com/elastic/go-ucfg"
	"time"
)

var defaultConfig = Config{
	FileMode:     0666,
	DirMode:      0777,
	SyncInterval: 60,
	StoreMode:    0,
}

type Config struct {
	Path         string        `config:"path, required"`
	FileMode     int           `config:"file_mode"`     // 0666
	DirMode      int           `config:"dir_mode"`      // 0777
	SyncInterval time.Duration `config:"sync_interval"` // in seconds
	StoreMode    int           `config:"store_mode"`    // 0 - file first, 2 - memory first(with persist on close), 2 - with empty file - memory without persist

}

func Unpack(cfg *ucfg.Config) (*Config, error) {

	config := defaultConfig
	err := cfg.Unpack(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
