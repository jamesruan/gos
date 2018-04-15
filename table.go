package gos

import (
	"github.com/Workiva/go-datastructures/trie/ctrie"
	spooky "github.com/dgryski/go-spooky"
	"hash"
	"sync"
)

// keyHash use spooky 32 hash for ctrie
func keyHash() hash.Hash32 {
	s1 := uint64(0)
	s2 := uint64(0)
	return spooky.New(s1, s2)
}

var gos_table_pool *sync.Pool

func init() {
	gos_table_pool = &sync.Pool{
		New: func() interface{} {
			return new(Table)
		},
	}
}

type Table struct {
	name string
	trie *ctrie.Ctrie
}

func newTable(name string) *Table {
	t := gos_table_pool.Get().(*Table)
	t.name = name
	t.trie = ctrie.New(keyHash)
	return t
}

// Size returns the current size of table
func (t Table) Size() uint {
	return t.trie.Size()
}

// Get loads the value for the key, returns nil if not existed
func (t Table) Get(key []byte) (got interface{}, ok bool) {
	return t.trie.Lookup(key)
}

// Put stores the value for the key, replacing the old value if the key existed
func (t Table) Put(key []byte, value interface{}) {
	t.trie.Insert(key, value)
}

// Delete removes a key, returns the removed
func (t Table) Delete(key []byte) (deleted interface{}, ok bool) {
	return t.trie.Remove(key)
}

// Copy returns a copy of current table
// Modification to the copy has no effect on original table
func (t Table) Copy() *Table {
	return &Table{
		name: t.name,
		trie: t.trie.Snapshot(),
	}
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// It is promised that every key/value pair is visited once. Any modification with the table during the range will not visible to f.
func (t Table) Range(f func(key []byte, value interface{}) bool) {
	trie := t.trie.ReadOnlySnapshot()
	cancel := make(chan struct{})
	entry := trie.Iterator(cancel)

	stop := false
	for e := range entry {
		if !stop {
			if f(e.Key, e.Value) {
				close(cancel)
				stop = true
			}
		}
	}
}
