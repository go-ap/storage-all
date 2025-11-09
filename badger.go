//go:build storage_badger

package storage

import storage "github.com/go-ap/storage-badger"

const Default = Badger

func New(initFns ...InitFn) (FullStorage, error) {
	opt, err := initConfig(initFns...)
	if err != nil {
		return nil, err
	}
	return getBadgerStorage(opt)
}

func Clean(initFns ...InitFn) error {
	opt, err := initConfig(initFns...)
	if err != nil {
		return err
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
	conf, err := getBadgerConfig(opt)
	if err != nil {
		return err
	}
	return storage.Bootstrap(conf)
}
