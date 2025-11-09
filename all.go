//go:build storage_all || !(storage_boltdb || storage_fs || storage_badger || storage_sqlite)

package storage

import "github.com/go-ap/errors"

const Default = FS

func New(initFns ...InitFn) (FullStorage, error) {
	opt, err := initConfig(initFns...)
	if err != nil {
		return nil, err
	}
	switch opt.Storage {
	case BoltDB:
		return getBoltStorage(opt)
	case Badger:
		return getBadgerStorage(opt)
	case Sqlite:
		return getSqliteStorage(opt)
	case FS:
		return getFsStorage(opt)
	}
	return nil, errors.NotImplementedf("Invalid storage type %s", opt.Storage)
}
