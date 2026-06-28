package storage

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-ap/errors"
)

type BackendConfig struct {
	Enabled bool
	Host    string
	Port    int64
	User    string
	Pw      string
	Name    string
}

const (
	varEnv     = "%env%"
	varStorage = "%storage%"
	varHost    = "%host%"
)

const defaultDirPerm = os.ModeDir | os.ModePerm | 0700

func normalizeConfigPath(p string, o options) string {
	if len(p) == 0 {
		return p
	}
	if p[0] == '~' {
		p = os.Getenv("HOME") + p[1:]
	}
	if !filepath.IsAbs(p) {
		p, _ = filepath.Abs(p)
	}
	p = strings.ReplaceAll(p, varEnv, string(o.Env))
	p = strings.ReplaceAll(p, varStorage, string(o.Storage))
	p = strings.ReplaceAll(p, varHost, url.PathEscape(o.Hostname))
	return filepath.Clean(p)
}

func (o options) BaseStoragePath() (string, error) {
	o.StoragePath = normalizeConfigPath(o.StoragePath, o)
	fi, err := os.Stat(o.StoragePath)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(o.StoragePath, defaultDirPerm)
	}
	if err != nil {
		return "", err
	}
	fi, err = os.Stat(o.StoragePath)
	if err != nil {
		return "", err
	}
	if !fi.IsDir() {
		return "", errors.BadRequestf("path %s is invalid for storage", o.StoragePath)
	}
	return o.StoragePath, nil
}
