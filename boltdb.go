//go:build storage_boltdb

package storage

import storage "github.com/go-ap/storage-boltdb"

const Default = BoltDB

func New(initFns ...InitFn) (FullStorage, error) {
	opt, err := initConfig(initFns...)
	if err != nil {
		return nil, err
	}
	return getBoltStorage(opt)
}

func Clean(initFns ...InitFn) error {
	opt, err := initConfig(initFns...)
	if err != nil {
		return err
	}
	conf, err := getBoltConfig(opt)
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
	conf, err := getBoltConfig(opt)
	if err != nil {
		return err
	}
	return storage.Bootstrap(conf)
}
