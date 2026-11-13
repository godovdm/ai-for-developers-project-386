// Package meeting holds the calendar meeting domain rules.
package meeting

import (
	"strings"
	"unicode/utf8"
)

// Duration is a meeting length in minutes.
type Duration int

const (
	// DurationThirty is a 30 minute meeting.
	DurationThirty Duration = 30
	// DurationSixty is a 60 minute meeting.
	DurationSixty Duration = 60
	// DurationNinety is a 90 minute meeting.
	DurationNinety Duration = 90
	// DurationOneHundredTwenty is a 120 minute meeting.
	DurationOneHundredTwenty Duration = 120
)

// MessageInvalidDuration describes the allowed meeting durations.
const MessageInvalidDuration = "Длительность должна быть 30, 60, 90 или 120 минут"

// Valid reports whether the duration is one of the supported values.
func (d Duration) Valid() bool {
	switch d {
	case DurationThirty, DurationSixty, DurationNinety, DurationOneHundredTwenty:
		return true
	default:
		return false
	}
}

// Type is a meeting type offered by the calendar owner.
type Type struct {
	ID          string
	Name        string
	Description string
	Duration    Duration
}

// ValidationError describes invalid meeting type data.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NewType validates and builds a meeting type.
func NewType(name, description string, duration Duration) (Type, error) {
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if name == "" {
		return Type{}, &ValidationError{Message: "Название не может быть пустым"}
	}
	if utf8.RuneCountInString(name) > 100 {
		return Type{}, &ValidationError{Message: "Название не должно быть длиннее 100 символов"}
	}
	if description == "" {
		return Type{}, &ValidationError{Message: "Описание не может быть пустым"}
	}
	if utf8.RuneCountInString(description) > 1000 {
		return Type{}, &ValidationError{Message: "Описание не должно быть длиннее 1000 символов"}
	}
	if !duration.Valid() {
		return Type{}, &ValidationError{Message: MessageInvalidDuration}
	}
	return Type{Name: name, Description: description, Duration: duration}, nil
}
