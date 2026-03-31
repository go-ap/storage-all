//go:build storage_all || !(storage_boltdb || storage_fs || storage_badger || storage_sqlite || storage_pg)

package storage

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"

	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/filters"
	"github.com/go-ap/storage-conformance-suite/gen"
)

func setBenchId(it vocab.Item) {
	u, _ := gen.RootID.URL()
	u.Path = ""
	base := u.String()
	_ = vocab.OnObject(it, setId(vocab.IRI(base)))
}

func setId(base vocab.IRI) func(ob *vocab.Object) error {
	idMap := sync.Map{}
	return func(ob *vocab.Object) error {
		typ := ob.Type
		id := 1
		if latestId, ok := idMap.Load(typ); ok {
			id = latestId.(int) + 1
		}

		ob.ID = base.AddPath(strings.ToLower(typeAsString(typ))).AddPath(strconv.Itoa(id))
		idMap.Store(typ, id)
		return nil
	}
}

func typeAsString(typ vocab.Typer) string {
	if tt, ok := typ.(vocab.ActivityVocabularyType); ok {
		return string(tt)
	}
	if tt, ok := typ.(vocab.ActivityVocabularyTypes); ok && len(tt) > 0 {
		return string(tt[0])
	}
	return "unknown"
}

var collectionIRI = vocab.Inbox.Of(gen.Root).GetLink()

var results = make(map[vocab.Typer]int)

func populate(st Store, count int) error {
	defer func() {
		gen.SetItemID = gen.DefaultSetter
	}()
	gen.SetItemID = setBenchId

	if _, err := st.Save(gen.Root); err != nil {
		return err
	}
	results[gen.Root.GetType()]++

	col := gen.RandomCollection(gen.Root)
	_ = vocab.OnObject(col, func(ob *vocab.Object) error {
		ob.ID = collectionIRI
		return nil
	})
	if _, err := st.Create(col); err != nil {
		return err
	}
	results[col.GetType()]++

	var success int
	var failure int

	items := gen.RandomItemCollection(count)
	for _, it := range items {
		if _, err := st.Save(it); err != nil {
			failure++
			continue
		}
		results[it.GetType()]++
		success++
	}
	if err := st.AddTo(col.GetLink(), items...); err != nil {
		return err
	}
	if failure > 0 {
		return fmt.Errorf("failed to save %d items", failure)
	}

	return nil
}

func setup(b *testing.B, typ Type, count int) (FullStorage, error) {
	tempDir := b.TempDir()
	initFns := []InitFn{WithType(typ), UseIndex(true)}
	if typ == Postgres {
		conf := setupContainer(b)
		initFns = append(initFns, WithPath(conf.DSN()))
	} else {
		initFns = append(initFns, WithPath(tempDir))
	}

	err := Bootstrap(initFns...)
	if err != nil {
		return nil, err
	}

	st, err := New(initFns...)
	if err != nil {
		return nil, err
	}

	if err = st.Open(); err != nil {
		return nil, err
	}

	if err = populate(st, count); err != nil {
		return nil, err
	}

	return st, nil
}

func _init(b *testing.B, typ Type) (FullStorage, error) {
	return setup(b, typ, count)
}

const count = 2000

var checks = filters.Checks{filters.HasType(vocab.NoteType, vocab.ArticleType)}

func Benchmark_Load_BoltDB(b *testing.B) {
	st, err := _init(b, BoltDB)
	if err != nil {
		b.Fatalf("unable to initialize storage %s", err)
	}

	b.ResetTimer()
	for b.Loop() {
		_, err = st.Load(collectionIRI, checks...)
		if err != nil {
			b.Errorf("unable to load from storage %s", err)
		}
	}
}

func Benchmark_Load_Sqlite(b *testing.B) {
	st, err := _init(b, Sqlite)
	if err != nil {
		b.Fatalf("unable to initialize storage %s", err)
	}
	defer st.Close()

	b.ResetTimer()
	for b.Loop() {
		_, err = st.Load(collectionIRI, checks...)
		if err != nil {
			b.Errorf("unable to load from storage %s", err)
		}
	}
}

func Benchmark_Load_Badger(b *testing.B) {
	st, err := _init(b, Badger)
	if err != nil {
		b.Fatalf("unable to initialize storage %s", err)
	}
	defer st.Close()

	b.ResetTimer()
	for b.Loop() {
		_, err = st.Load(collectionIRI, checks...)
		if err != nil {
			b.Errorf("unable to load from storage %s", err)
		}
	}
}

func Benchmark_Load_FS(b *testing.B) {
	st, err := _init(b, FS)
	if err != nil {
		b.Fatalf("unable to initialize storage %s", err)
	}
	defer st.Close()

	b.ResetTimer()
	for b.Loop() {
		_, err = st.Load(collectionIRI, checks...)
		if err != nil {
			b.Errorf("unable to load from storage %s", err)
		}
	}
}

func Benchmark_Load_Postgres(b *testing.B) {
	st, err := _init(b, Postgres)
	if err != nil {
		b.Fatalf("unable to initialize storage %s", err)
	}
	defer st.Close()

	b.ResetTimer()
	for b.Loop() {
		_, err = st.Load(collectionIRI, checks...)
		if err != nil {
			b.Errorf("unable to load from storage %s", err)
		}
	}
}
