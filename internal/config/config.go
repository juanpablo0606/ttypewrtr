package config

import (
	"errors"
	"fmt"
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
			"quit":   "ctrl+q",
			"save":   "ctrl+s",
		},
	}
}

func Path() (string, error) {
	p, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(p, "ttypewrtr", "config.toml"), nil

}

func Load(path string) (Config, error) {
	cfg := Default()
	md, err := toml.DecodeFile(path, &cfg)
	if errors.Is(err, fs.ErrNotExist) {
		return cfg, nil
	} else if err != nil {
		return Config{}, err
	}

	if undecoded := md.Undecoded(); len(undecoded) > 0 {
		return Config{}, fmt.Errorf("unkown setting %q in %s", undecoded[0].String(), path)
	}
	return cfg, nil
}
