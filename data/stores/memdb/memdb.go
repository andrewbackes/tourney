// Package memdb is an in memory database.
package memdb

import (
	"github.com/andrewbackes/tourney/data/models"
	"sync"
)

// MemDB is an in memory database.
type MemDB struct {
	tournaments sync.Map
	games       sync.Map
	locks       sync.Map
	workers     sync.Map
}

// NewMemDB creates a new in memory database.
func NewMemDB() *MemDB {
	return &MemDB{
		tournaments: sync.Map{},
		locks:       sync.Map{},
		games:       sync.Map{},
		workers:     sync.Map{},
	}
}

func (m *MemDB) lock(id models.Id) {
	lock, exists := m.locks.Load(id)
	if !exists {
		panic("missing required element in map")
	}
	lock.(*sync.Mutex).Lock()
}

func (m *MemDB) unlock(id models.Id) {
	lock, exists := m.locks.Load(id)
	if !exists {
		panic("missing required element in map")
	}
	lock.(*sync.Mutex).Unlock()
}
