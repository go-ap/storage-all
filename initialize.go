package storage

import (
	"git.sr.ht/~mariusor/lw"
	"github.com/go-ap/errors"
	"github.com/go-ap/storage-badger"
	"github.com/go-ap/storage-boltdb"
	"github.com/go-ap/storage-fs"
	"github.com/go-ap/storage-sqlite"
)

type Type string

type options struct {
	Env          string
	Hostname     string
	Storage      Type
	StoragePath  string
	StorageCache bool
	UseIndex     bool
	Logger       lw.Logger
}

const (
	BoltDB = Type("boltdb")
	FS     = Type("fs")
	Badger = Type("badger")
	Sqlite = Type("sqlite")

	//Postgres = Type("postgres")
)

type InitFn func(t *options) error

func WithEnv(env string) InitFn {
	return func(o *options) error {
		o.Env = env
		return nil
	}
}

func WithHostname(hostname string) InitFn {
	return func(o *options) error {
		o.Hostname = hostname
		return nil
	}
}

func WithType(t Type) InitFn {
	return func(o *options) error {
		o.Storage = t
		return nil
	}
}

func WithPath(p string) InitFn {
	return func(o *options) error {
		o.StoragePath = p
		return nil
	}
}

func WithCache(enabled bool) InitFn {
	return func(o *options) error {
		o.StorageCache = enabled
		return nil
	}
}

func WithLogger(l lw.Logger) InitFn {
	return func(o *options) error {
		o.Logger = l
		return nil
	}
}

func initConfig(initFns ...InitFn) (options, error) {
	opt := options{
		Storage: Default,
		Logger:  lw.Nil(),
	}
	for _, fn := range initFns {
		err := fn(&opt)
		if err != nil {
			return opt, err
		}
	}
	return opt, nil
}

func getBadgerStorage(opt options) (FullStorage, error) {
	opt.Storage = Badger
	path, err := opt.BaseStoragePath()
	if err != nil {
		return nil, err
	}
	l := opt.Logger
	if l != nil {
		l = l.WithContext(lw.Ctx{"path": path, "storage": opt.Storage})
	}
	conf := badger.Config{
		Path:  path,
		LogFn: l.Debugf,
		ErrFn: l.Warnf,
	}
	db, err := badger.New(conf)
	if err != nil {
		return db, err
	}
	return db, nil
}

func getBoltStorage(opt options) (FullStorage, error) {
	opt.Storage = BoltDB
	path, err := opt.BaseStoragePath()
	if err != nil {
		return nil, err
	}
	l := opt.Logger
	if l != nil {
		l = l.WithContext(lw.Ctx{"path": path, "storage": opt.Storage})
	}
	db, err := boltdb.New(boltdb.Config{
		Path:  path,
		LogFn: l.Debugf,
		ErrFn: l.Warnf,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getFsStorage(opt options) (FullStorage, error) {
	opt.Storage = FS
	path, err := opt.BaseStoragePath()
	if err != nil {
		return nil, err
	}
	l := opt.Logger
	if l != nil {
		l = l.WithContext(lw.Ctx{"path": path, "storage": opt.Storage})
	}
	db, err := fs.New(fs.Config{
		Path:        path,
		CacheEnable: opt.StorageCache,
		Logger:      l,
		UseIndex:    opt.UseIndex,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getSqliteStorage(opt options) (FullStorage, error) {
	opt.Storage = Sqlite
	path, err := opt.BaseStoragePath()
	if err != nil {
		return nil, err
	}
	l := opt.Logger
	if l != nil {
		l = l.WithContext(lw.Ctx{"path": path, "storage": opt.Storage})
	}
	db, err := sqlite.New(sqlite.Config{
		Path:        path,
		CacheEnable: opt.StorageCache,
		LogFn:       l.Debugf,
		ErrFn:       l.Warnf,
	})

	if err != nil {
		return nil, errors.Annotatef(err, "unable to connect to sqlite storage")
	}
	return db, nil
}
