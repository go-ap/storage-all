//go:build storage_badger

package storage

import (
	"github.com/go-ap/errors"
	storage "github.com/go-ap/storage-badger"
)

const Default = Badger

func New(initFns ...InitFn) (FullStorage, error) {
	opt, err := initConfig(initFns...)
	if err != nil {
		return nil, err
	}
	if opt.Storage != Default {
		return nil, errors.NotImplementedf("Invalid storage type %s expected %s", opt.Storage, Default)
	}
	return getBadgerStorage(opt)
}

func Clean(initFns ...InitFn) error {
	opt, err := initConfig(initFns...)
	if err != nil {
		return err
	}
	if opt.Storage != Default {
		return errors.Newf("invalid storage type %s, expected %s", opt.Storage, Default)
	}
	conf, err := getBadgerConfig(opt)
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
	conf, err := getBadgerConfig(opt)
	if err != nil {
		return err
	}
	return storage.Bootstrap(conf)
}
