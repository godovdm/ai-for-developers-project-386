package meeting

import (
	"strconv"
	"sync"
)

// Repository stores meeting types.
type Repository interface {
	Create(Type) Type
	List() []Type
}

// Memory is a thread-safe in-memory meeting type repository.
type Memory struct {
	mu     sync.Mutex
	nextID int
	items  []Type
}

// NewMemory returns an empty in-memory repository.
func NewMemory() *Memory {
	return &Memory{}
}

// Create assigns an identifier and stores the meeting type.
func (m *Memory) Create(meetingType Type) Type {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nextID++
	meetingType.ID = strconv.Itoa(m.nextID)
	m.items = append(m.items, meetingType)
	return meetingType
}

// List returns the meeting types in creation order.
func (m *Memory) List() []Type {
	m.mu.Lock()
	defer m.mu.Unlock()
	items := make([]Type, len(m.items))
	copy(items, m.items)
	return items
}
