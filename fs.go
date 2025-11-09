//go:build storage_fs

package storage

const Default = FS

func New(initFns ...InitFn) (FullStorage, error) {
	opt, err := initConfig(initFns...)
	if err != nil {
		return nil, err
	}
	return getFsStorage(opt)
}
