//go:build storage_sqlite

package storage

const Default = Sqlite

func New(initFns ...InitFn) (FullStorage, error) {
	opt, err := initConfig(initFns...)
	if err != nil {
		return nil, err
	}
	return getSqliteStorage(opt)
}
