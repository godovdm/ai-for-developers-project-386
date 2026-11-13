package meeting

import (
	"strconv"
	"sync"
	"time"
)

// Repository stores meeting types.
type Repository interface {
	Create(Type) Type
	List() []Type
}

// Memory is a thread-safe in-memory meeting type repository.
type Memory struct {
	mu       sync.Mutex
	nextID   int
	items    []Type
	bookings []Booking
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

// Booking represents a reserved interval.
type Booking struct {
	Start         time.Time
	End           time.Time
	MeetingTypeID string
	GuestName     string
	GuestEmail    string
}

// Reserve stores an interval only if it does not overlap any reservation.
func (m *Memory) Reserve(booking Booking) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, existing := range m.bookings {
		if booking.Start.Before(existing.End) && existing.Start.Before(booking.End) {
			return false
		}
	}
	m.bookings = append(m.bookings, booking)
	return true
}

// Available reports whether an interval is currently free.
func (m *Memory) Available(start, end time.Time) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, booking := range m.bookings {
		if start.Before(booking.End) && booking.Start.Before(end) {
			return false
		}
	}
	return true
}

// Bookings returns a snapshot of all reserved intervals.
func (m *Memory) Bookings() []Booking {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]Booking(nil), m.bookings...)
}
