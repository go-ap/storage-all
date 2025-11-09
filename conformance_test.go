//go:build conformance

package storage

import (
	"testing"

	"git.sr.ht/~mariusor/lw"
	conformance "github.com/go-ap/storage-conformance-suite"
)

func initStorage(t *testing.T) conformance.ActivityPubStorage {
	l := lw.Dev(lw.SetOutput(t.Output()))
	initFns := []InitFn{WithPath(t.TempDir()), WithLogger(l)}
	if err := Bootstrap(initFns...); err != nil {
		t.Fatalf("unable to bootstrap storage: %s", err)
	}
	storage, err := New(initFns...)
	if err != nil {
		t.Fatalf("unable to initialize storage: %s", err)
	}
	return storage
}

func Test_Conformance(t *testing.T) {
	conformance.Suite(
		conformance.TestActivityPub, conformance.TestMetadata,
		conformance.TestKey, conformance.TestOAuth, conformance.TestPassword,
	).Run(t, initStorage(t))
}
