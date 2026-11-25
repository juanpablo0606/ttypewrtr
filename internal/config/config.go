package config

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Keys map[string]string `toml:"keys"`
}

func Default() Config {
	return Config{
		Keys: map[string]string{
			"h1":     "alt+1",
			"h2":     "alt+2",
			"h3":     "alt+3",
			"h4":     "alt+4",
			"h5":     "alt+5",
			"h6":     "alt+6",
			"bold":   "alt+b",
			"italic": "alt+i",
			"quote":  "alt+q",
			"list":   "alt+l",
		},
	}
}

func Path() (string, error) {
	p, err := os.UserConfigDir()
	if err != nil {
		return "", nil
	}
	return filepath.Join(p, "config.toml"), nil

}

func Load(path string) (Config, error) {
	cfg := Default()
	path, err := Path()
	if err != nil {
		return cfg, err
	}
	_, err = toml.DecodeFile(path, &cfg)
	if errors.Is(err, fs.ErrNotExist) {
		return Model{}, err
	} else if err != nil {
		return Model{}, err
	}

}
