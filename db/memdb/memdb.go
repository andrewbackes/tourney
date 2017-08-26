// Package memdb is an in memory database.
package memdb

import (
	"sync"
)

// MemDB is an in memory database.
type MemDB struct {
	tournaments sync.Map
	games       sync.Map
}

// NewMemDB creates a new in memory database.
func NewMemDB() *MemDB {
	return &MemDB{
		tournaments: sync.Map{},
		games:       sync.Map{},
	}
}
