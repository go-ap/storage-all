//go:build storage_all || !(storage_boltdb || storage_fs || storage_badger || storage_sqlite)

package storage

import (
	"github.com/go-ap/errors"
	badger "github.com/go-ap/storage-badger"
	boltdb "github.com/go-ap/storage-boltdb"
	fs "github.com/go-ap/storage-fs"
	sqlite "github.com/go-ap/storage-sqlite"
)

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

func Clean(initFns ...InitFn) error {
	opt, err := initConfig(initFns...)
	if err != nil {
		return err
	}
	switch opt.Storage {
	case BoltDB:
		conf, err := getBoltConfig(opt)
		if err != nil {
			return err
		}
		return boltdb.Clean(conf)
	case Badger:
		conf, err := getBadgerConfig(opt)
		if err != nil {
			return err
		}
		return badger.Clean(conf)
	case Sqlite:
		conf, err := getSqliteConfig(opt)
		if err != nil {
			return err
		}
		return sqlite.Clean(conf)
	case FS:
		conf, err := getFsConfig(opt)
		if err != nil {
			return err
		}
		return fs.Clean(conf)
	}
	return errors.NotImplementedf("Invalid storage type %s", opt.Storage)
}

func Bootstrap(initFns ...InitFn) error {
	opt, err := initConfig(initFns...)
	if err != nil {
		return err
	}
	switch opt.Storage {
	case BoltDB:
		conf, err := getBoltConfig(opt)
		if err != nil {
			return err
		}
		return boltdb.Bootstrap(conf)
	case Badger:
		conf, err := getBadgerConfig(opt)
		if err != nil {
			return err
		}
		return badger.Bootstrap(conf)
	case Sqlite:
		conf, err := getSqliteConfig(opt)
		if err != nil {
			return err
		}
		return sqlite.Bootstrap(conf)
	case FS:
		conf, err := getFsConfig(opt)
		if err != nil {
			return err
		}
		return fs.Bootstrap(conf)
	}
	return errors.NotImplementedf("Invalid storage type %s", opt.Storage)
}
