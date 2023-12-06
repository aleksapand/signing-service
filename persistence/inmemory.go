package persistence

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/uuid"
	"sync"
)

var db *InMemoryDB

// InMemoryDB is a struct that holds a devices map.
type InMemoryDB struct {
	data map[uuid.UUID]*domain.SignatureDevice
	mu   sync.RWMutex
}

// GetInMemoryDB returns the instance of InMemoryDB
func GetInMemoryDB() *InMemoryDB {

	db = &InMemoryDB{
		data: make(map[uuid.UUID]*domain.SignatureDevice),
	}
	return db
}

// Set sets a device in the in-memory database.
func (db *InMemoryDB) Set(key uuid.UUID, device *domain.SignatureDevice) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = device
}

// Get retrieves the device associated with the specified key.
func (db *InMemoryDB) Get(key uuid.UUID) (*domain.SignatureDevice, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	value, ok := db.data[key]
	return value, ok
}

// GetAll retrieves all devices.
func (db *InMemoryDB) GetAll() []*domain.SignatureDevice {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var allDevices []*domain.SignatureDevice
	for _, value := range db.data {
		allDevices = append(allDevices, value)
	}

	return allDevices
}
