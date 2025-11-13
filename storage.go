package storage

import (
	"crypto"

	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/filters"
	"github.com/openshift/osin"
)

type ClientSaver interface {
	// UpdateClient updates the client (identified by its id) and replaces the values with the values of client.
	UpdateClient(c osin.Client) error
	// CreateClient stores the client in the database and returns an error, if something went wrong.
	CreateClient(c osin.Client) error
	// RemoveClient removes a client (identified by id) from the database. Returns an error if something went wrong.
	RemoveClient(id string) error
}

type ClientLister interface {
	// ListClients lists existing clients
	ListClients() ([]osin.Client, error)
	GetClient(id string) (osin.Client, error)
}

type Store interface {
	ReadStore
	WriteStore
	CollectionStore
}
type ReadStore interface {
	// Load returns an Item or an ItemCollection from an IRI
	// after filtering it through the FilterFn list of filtering functions. Eg ANY()
	Load(vocab.IRI, ...filters.Check) (vocab.Item, error)
}

// WriteStore saves ActivityStreams objects.
type WriteStore interface {
	// Save saves the incoming ActivityStreams Object, and returns it together with any properties
	// populated by the method's side effects. (eg, Published property can point to the current time, etc.).
	Save(vocab.Item) (vocab.Item, error)
	// Delete deletes completely from storage the ActivityStreams Object
	Delete(vocab.Item) error
}

// CollectionStore allows operations on ActivityStreams collections
type CollectionStore interface {
	// Create creates the "col" collection.
	Create(vocab.CollectionInterface) (vocab.CollectionInterface, error)
	// AddTo adds "it" element to the "col" collection.
	AddTo(vocab.IRI, ...vocab.Item) error
	// RemoveFrom removes "it" item from "col" collection
	RemoveFrom(vocab.IRI, ...vocab.Item) error
}
type KeyLoader interface {
	LoadKey(vocab.IRI) (crypto.PrivateKey, error)
}

type FullStorage interface {
	Open() error
	ClientSaver
	ClientLister
	Store
	KeyLoader
	MetadataStorage
	PasswordChanger
	osin.Storage
}

type PasswordChanger interface {
	PasswordSet(vocab.IRI, []byte) error
	PasswordCheck(vocab.IRI, []byte) error
}

type MetadataStorage interface {
	LoadMetadata(vocab.IRI, any) error
	SaveMetadata(vocab.IRI, any) error
}

type MimeTypeSaver interface {
	SaveNaturalLanguageValues(vocab.NaturalLanguageValues) error
	SaveMimeTypeContent(vocab.MimeType, vocab.NaturalLanguageValues) error
}
