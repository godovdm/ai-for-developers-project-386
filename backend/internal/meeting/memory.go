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
	for {
		m.nextID++
		meetingType.ID = strconv.Itoa(m.nextID)
		found := false
		for _, item := range m.items {
			if item.ID == meetingType.ID {
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	m.items = append(m.items, meetingType)
	return meetingType
}

// CreateWithID stores a meeting type if its identifier is unique.
func (m *Memory) CreateWithID(meetingType Type) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, item := range m.items {
		if item.ID == meetingType.ID {
			return false
		}
	}
	m.items = append(m.items, meetingType)
	return true
}

// List returns the meeting types in creation order.
func (m *Memory) List() []Type {
	m.mu.Lock()
	defer m.mu.Unlock()
	items := make([]Type, len(m.items))
	copy(items, m.items)
	return items
}
