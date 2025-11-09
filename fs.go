//go:build storage_fs

package storage

import storage "github.com/go-ap/storage-fs"

const Default = FS

func New(initFns ...InitFn) (FullStorage, error) {
	opt, err := initConfig(initFns...)
	if err != nil {
		return nil, err
	}
	return getFsStorage(opt)
}

func Clean(initFns ...InitFn) error {
	opt, err := initConfig(initFns...)
	if err != nil {
		return err
	}
	conf, err := getFsConfig(opt)
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
	conf, err := getFsConfig(opt)
	if err != nil {
		return err
	}
	return storage.Bootstrap(conf)
}
