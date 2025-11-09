//go:build storage_sqlite

package storage

import storage "github.com/go-ap/storage-sqlite"

const Default = Sqlite

func New(initFns ...InitFn) (FullStorage, error) {
	opt, err := initConfig(initFns...)
	if err != nil {
		return nil, err
	}
	return getSqliteStorage(opt)
}

func Clean(initFns ...InitFn) error {
	opt, err := initConfig(initFns...)
	if err != nil {
		return err
	}
	conf, err := getSqliteConfig(opt)
	if err != nil {
		return err
	}
	return storage.Clean(conf)
}

func Bootstrap(initFns ...InitFn) error {
	opt, err := initConfig(initFns...)
	if err != nil {
		return err
	}
	conf, err := getSqliteConfig(opt)
	if err != nil {
		return err
	}
	return storage.Bootstrap(conf)
}
