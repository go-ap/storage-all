//go:build storage_sqlite

package storage

import (
	"github.com/go-ap/errors"
	storage "github.com/go-ap/storage-sqlite"
)

const Default = Sqlite

func New(initFns ...InitFn) (FullStorage, error) {
	opt, err := initConfig(initFns...)
	if err != nil {
		return nil, err
	}
	if opt.Storage != Default {
		return nil, errors.NotImplementedf("Invalid storage type %s expected %s", opt.Storage, Default)
	}
	return getSqliteStorage(opt)
}

func Clean(initFns ...InitFn) error {
	opt, err := initConfig(initFns...)
	if err != nil {
		return err
	}
	if opt.Storage != Default {
		return errors.Newf("invalid storage type %s, expected %s", opt.Storage, Default)
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
	if opt.Storage != Default {
		return errors.Newf("invalid storage type %s, expected %s", opt.Storage, Default)
	}
	conf, err := getSqliteConfig(opt)
	if err != nil {
		return err
	}
	return storage.Bootstrap(conf)
}
