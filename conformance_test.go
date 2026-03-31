//go:build conformance

package storage

import (
	"testing"

	"git.sr.ht/~mariusor/lw"
	conformance "github.com/go-ap/storage-conformance-suite"
)

func initStorage(t *testing.T, typ Type) conformance.ActivityPubStorage {
	l := lw.Dev(lw.SetOutput(t.Output()))
	initFns := []InitFn{WithType(typ), WithPath(t.TempDir()), WithLogger(l)}
	if err := Bootstrap(initFns...); err != nil {
		t.Fatalf("unable to bootstrap storage: %s", err)
	}
	storage, err := New(initFns...)
	if err != nil {
		t.Fatalf("unable to initialize storage: %s", err)
	}
	return storage
}

func Test_Conformance_FS(t *testing.T) {
	conformance.Suite(
		conformance.TestActivityPub, conformance.TestMetadata,
		conformance.TestKey, conformance.TestOAuth, conformance.TestPassword,
	).Run(t, initStorage(t, FS))
}

func Test_Conformance_Sqlite(t *testing.T) {
	conformance.Suite(
		conformance.TestActivityPub, conformance.TestMetadata,
		conformance.TestKey, conformance.TestOAuth, conformance.TestPassword,
	).Run(t, initStorage(t, Sqlite))
}

func Test_Conformance_BoltDB(t *testing.T) {
	conformance.Suite(
		conformance.TestActivityPub, conformance.TestMetadata,
		conformance.TestKey, conformance.TestOAuth, conformance.TestPassword,
	).Run(t, initStorage(t, BoltDB))
}

func Test_Conformance_Badger(t *testing.T) {
	conformance.Suite(
		conformance.TestActivityPub, conformance.TestMetadata,
		conformance.TestKey, conformance.TestOAuth, conformance.TestPassword,
	).Run(t, initStorage(t, Badger))
}
