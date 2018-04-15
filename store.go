package gos

import (
	"sync"
)

// the global store
var gos_store *sync.Map

func init() {
	gos_store = new(sync.Map)
}

// Get returns a table with name, creates one if not existed.
func Get(name string) *Table {
	if v, ok := gos_store.Load(name); !ok {
		t := newTable(name)
		gos_store.Store(name, t)
		return t
	} else {
		return v.(*Table)
	}
}

// Delete removes the table if it existed.
func Delete(t *Table) {
	if t != nil {
		gos_store.Delete(t.name)
		gos_table_pool.Put(t)
	}
}
