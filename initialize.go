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

func UseIndex(enabled bool) InitFn {
	return func(o *options) error {
		o.UseIndex = enabled
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

func getBadgerConfig(opt options) (badger.Config, error) {
	opt.Storage = Badger
	path, err := opt.BaseStoragePath()
	if err != nil {
		return badger.Config{}, err
	}
	l := opt.Logger
	if l != nil {
		l = l.WithContext(lw.Ctx{"path": path, "storage": opt.Storage})
	}
	return badger.Config{
		Path:  path,
		LogFn: l.Debugf,
		ErrFn: l.Warnf,
	}, nil
}

func getBadgerStorage(opt options) (FullStorage, error) {
	conf, err := getBadgerConfig(opt)
	if err != nil {
		return nil, err
	}
	db, err := badger.New(conf)
	if err != nil {
		return db, err
	}
	return db, nil
}

func getBoltConfig(opt options) (boltdb.Config, error) {
	opt.Storage = BoltDB
	path, err := opt.BaseStoragePath()
	if err != nil {
		return boltdb.Config{}, err
	}
	l := opt.Logger
	if l != nil {
		l = l.WithContext(lw.Ctx{"path": path, "storage": opt.Storage})
	}
	return boltdb.Config{
		Path:  path,
		LogFn: l.Debugf,
		ErrFn: l.Warnf,
	}, nil
}

func getBoltStorage(opt options) (FullStorage, error) {
	conf, err := getBoltConfig(opt)
	if err != nil {
		return nil, err
	}
	db, err := boltdb.New(conf)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getFsConfig(opt options) (fs.Config, error) {
	opt.Storage = FS
	path, err := opt.BaseStoragePath()
	if err != nil {
		return fs.Config{}, err
	}
	l := opt.Logger
	if l != nil {
		l = l.WithContext(lw.Ctx{"path": path, "storage": opt.Storage})
	}
	return fs.Config{
		Path:                     path,
		Logger:                   l,
		EnableCache:              opt.StorageCache,
		EnableIndex:              opt.UseIndex,
		EnableOptimizedFiltering: true,
	}, nil
}

func getFsStorage(opt options) (FullStorage, error) {
	conf, err := getFsConfig(opt)
	if err != nil {
		return nil, err
	}
	db, err := fs.New(conf)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getSqliteConfig(opt options) (sqlite.Config, error) {
	opt.Storage = Sqlite
	path, err := opt.BaseStoragePath()
	if err != nil {
		return sqlite.Config{}, err
	}
	l := opt.Logger
	if l != nil {
		l = l.WithContext(lw.Ctx{"path": path, "storage": opt.Storage})
	}
	return sqlite.Config{
		Path:        path,
		CacheEnable: opt.StorageCache,
		LogFn:       l.Debugf,
		ErrFn:       l.Warnf,
	}, nil
}
func getSqliteStorage(opt options) (FullStorage, error) {
	conf, err := getSqliteConfig(opt)
	if err != nil {
		return nil, err
	}
	db, err := sqlite.New(conf)
	if err != nil {
		return nil, errors.Annotatef(err, "unable to connect to sqlite storage")
	}
	return db, nil
}
