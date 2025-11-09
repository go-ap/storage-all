//go:build storage_boltdb

package storage

const Default = BoltDB

func New(initFns ...InitFn) (FullStorage, error) {
	opt, err := initConfig(initFns...)
	if err != nil {
		return nil, err
	}
	return getBoltStorage(opt)
}
