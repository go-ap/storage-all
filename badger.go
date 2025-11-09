//go:build storage_badger

package storage

const Default = Badger

func New(initFns ...InitFn) (FullStorage, error) {
	opt, err := initConfig(initFns...)
	if err != nil {
		return nil, err
	}
	return getBadgerStorage(opt)
}
