package meeting

import (
	"errors"
	"strings"
	"testing"
)

func TestDurationValid(t *testing.T) {
	t.Parallel()
	valid := []Duration{DurationThirty, DurationSixty, DurationNinety, DurationOneHundredTwenty}
	for _, duration := range valid {
		if !duration.Valid() {
			t.Errorf("Duration(%d).Valid() = false, want true", duration)
		}
	}
	invalid := []Duration{0, 15, 45, 121}
	for _, duration := range invalid {
		if duration.Valid() {
			t.Errorf("Duration(%d).Valid() = true, want false", duration)
		}
	}
}

func TestNewTypeValid(t *testing.T) {
	t.Parallel()
	meetingType, err := NewType("  Вводный звонок  ", "  Знакомство  ", DurationSixty)
	if err != nil {
		t.Fatalf("NewType() error = %v", err)
	}
	if meetingType.Name != "Вводный звонок" {
		t.Errorf("Name = %q, want trimmed name", meetingType.Name)
	}
	if meetingType.Description != "Знакомство" {
		t.Errorf("Description = %q, want trimmed description", meetingType.Description)
	}
	if meetingType.Duration != DurationSixty {
		t.Errorf("Duration = %d, want %d", meetingType.Duration, DurationSixty)
	}
}

func TestNewTypeInvalid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		typeName    string
		description string
		duration    Duration
		message     string
	}{
		{
			name:        "empty name",
			typeName:    "",
			description: "Описание",
			duration:    DurationThirty,
			message:     "Название не может быть пустым",
		},
		{
			name:        "whitespace name",
			typeName:    "   ",
			description: "Описание",
			duration:    DurationThirty,
			message:     "Название не может быть пустым",
		},
		{
			name:        "name too long",
			typeName:    strings.Repeat("а", 101),
			description: "Описание",
			duration:    DurationThirty,
			message:     "Название не должно быть длиннее 100 символов",
		},
		{
			name:        "empty description",
			typeName:    "Звонок",
			description: "  ",
			duration:    DurationThirty,
			message:     "Описание не может быть пустым",
		},
		{
			name:        "description too long",
			typeName:    "Звонок",
			description: strings.Repeat("а", 1001),
			duration:    DurationThirty,
			message:     "Описание не должно быть длиннее 1000 символов",
		},
		{
			name:        "invalid duration",
			typeName:    "Звонок",
			description: "Описание",
			duration:    Duration(45),
			message:     MessageInvalidDuration,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewType(test.typeName, test.description, test.duration)
			if err == nil {
				t.Fatal("NewType() error = nil, want validation error")
			}
			var validationError *ValidationError
			if !errors.As(err, &validationError) {
				t.Fatalf("NewType() error = %v, want *ValidationError", err)
			}
			if validationError.Message != test.message {
				t.Errorf("Message = %q, want %q", validationError.Message, test.message)
			}
			if err.Error() != test.message {
				t.Errorf("Error() = %q, want %q", err.Error(), test.message)
			}
		})
	}
}
