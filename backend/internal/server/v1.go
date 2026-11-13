package server

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/mail"
	"regexp"
	"slices"
	"strings"
	"time"
	_ "time/tzdata" // Embed IANA time zone data for deployments without system zoneinfo.
	"unicode/utf8"

	"github.com/godovdm/ai-for-developers-project-386/backend/internal/api/generated"
	"github.com/godovdm/ai-for-developers-project-386/backend/internal/meeting"
)

var validID = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,62}[a-z0-9])?$`)

func (s server) ListMeetingTypesV1(w http.ResponseWriter, _ *http.Request) {
	items := make([]generated.MeetingTypeV1, 0)
	for _, item := range s.v1.List() {
		items = append(items, toV1(item))
	}
	writeJSON(w, http.StatusOK, generated.MeetingTypeListV1{Items: items})
}
func (s server) GetMeetingTypeV1(w http.ResponseWriter, _ *http.Request, id string) {
	for _, item := range s.v1.List() {
		if item.ID == id {
			writeJSON(w, http.StatusOK, toV1(item))
			return
		}
	}
	writeV1Error(w, http.StatusNotFound, "meeting_type_not_found", "Тип встречи не найден")
}
func (s server) CreateMeetingTypeV1(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("close request body: %v", err)
		}
	}()
	var req struct {
		ID              *string         `json:"id"`
		Name            *string         `json:"name"`
		Description     *string         `json:"description"`
		DurationMinutes json.RawMessage `json:"durationMinutes"`
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxRequestBodySize))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeV1Error(w, 400, "validation_error", "Некорректный запрос")
		return
	}
	var trailing json.RawMessage
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) || req.ID == nil || req.Name == nil || req.Description == nil {
		writeV1Error(w, 400, "validation_error", "Некорректный запрос")
		return
	}
	name, description := strings.TrimSpace(*req.Name), strings.TrimSpace(*req.Description)
	duration, ok := parseDuration(req.DurationMinutes)
	if !validID.MatchString(*req.ID) {
		writeV1Error(w, 400, "validation_error", "ID должен содержать от 1 до 64 строчных латинских букв, цифр или дефисов и начинаться и заканчиваться буквой или цифрой")
		return
	}
	if name == "" || utf8.RuneCountInString(name) > 100 {
		writeV1Error(w, 400, "validation_error", "Название должно содержать от 1 до 100 символов")
		return
	}
	if utf8.RuneCountInString(description) > 2000 {
		writeV1Error(w, 400, "validation_error", "Описание не должно быть длиннее 2000 символов")
		return
	}
	if !ok {
		writeV1Error(w, 400, "validation_error", meeting.MessageInvalidDuration)
		return
	}
	item := meeting.Type{ID: *req.ID, Name: name, Description: description, Duration: duration}
	if !s.v1.CreateWithID(item) {
		writeV1Error(w, 409, "meeting_type_exists", "Тип встречи с таким id уже существует")
		return
	}
	writeJSON(w, http.StatusCreated, toV1(item))
}
func writeV1Error(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, generated.ApiErrorV1{Code: code, Message: message})
}

func toV1(item meeting.Type) generated.MeetingTypeV1 {
	return generated.MeetingTypeV1{Id: item.ID, Name: item.Name, Description: item.Description, DurationMinutes: generated.MeetingDuration(item.Duration)}
}

func (s server) GetAvailability(w http.ResponseWriter, _ *http.Request, id string, params generated.GetAvailabilityParams) {
	var duration meeting.Duration
	for _, item := range s.v1.List() {
		if item.ID == id {
			duration = item.Duration
			break
		}
	}
	if duration == 0 {
		writeV1Error(w, http.StatusNotFound, "meeting_type_not_found", "Тип встречи не найден")
		return
	}
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		writeV1Error(w, http.StatusInternalServerError, "timezone_error", "Не удалось определить часовой пояс")
		return
	}
	now := s.now().In(location)
	first := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, location)
	days := make([]generated.AvailabilityDay, 0, 14)
	for i := range 14 {
		date := first.AddDate(0, 0, i)
		dateString := date.Format(time.DateOnly)
		if (params.From != nil && dateString < *params.From) || (params.To != nil && dateString > *params.To) {
			continue
		}
		slots := make([]string, 0)
		if date.Weekday() != time.Saturday && date.Weekday() != time.Sunday {
			for minute := 10 * 60; minute+int(duration) <= 18*60; minute += 30 {
				start := date.Add(time.Duration(minute) * time.Minute)
				if s.v1.Available(start, start.Add(time.Duration(duration)*time.Minute)) {
					slots = append(slots, start.Format(time.RFC3339))
				}
			}
		}
		days = append(days, generated.AvailabilityDay{Date: date.Format(time.DateOnly), Slots: slots})
	}
	writeJSON(w, http.StatusOK, generated.Availability{Days: days})
}

func (s server) CreateBooking(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("close request body: %v", err)
		}
	}()
	var req generated.CreateBookingRequest
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxRequestBodySize))
	decoder.DisallowUnknownFields()
	var trailing json.RawMessage
	if decoder.Decode(&req) != nil || !errors.Is(decoder.Decode(&trailing), io.EOF) {
		writeV1Error(w, http.StatusBadRequest, "validation_error", "Некорректный запрос")
		return
	}
	address, err := mail.ParseAddress(req.GuestEmail)
	if strings.TrimSpace(req.GuestName) == "" || err != nil || address.Address != req.GuestEmail {
		writeV1Error(w, http.StatusBadRequest, "validation_error", "Укажите имя и корректный email")
		return
	}
	var item meeting.Type
	for _, candidate := range s.v1.List() {
		if candidate.ID == req.MeetingTypeId {
			item = candidate
			break
		}
	}
	if item.ID == "" {
		writeV1Error(w, http.StatusNotFound, "meeting_type_not_found", "Тип встречи не найден")
		return
	}
	start, err := time.Parse(time.RFC3339, req.Start)
	if err != nil {
		writeV1Error(w, http.StatusBadRequest, "validation_error", "Некорректное время")
		return
	}
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		writeV1Error(w, http.StatusInternalServerError, "timezone_error", "Не удалось определить часовой пояс")
		return
	}
	start = start.In(location)
	now := s.now().In(location)
	first := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, location)
	minute := start.Hour()*60 + start.Minute()
	if start.Before(first) || !start.Before(first.AddDate(0, 0, 14)) || start.Weekday() == time.Saturday || start.Weekday() == time.Sunday || minute < 600 || minute+int(item.Duration) > 1080 || minute%30 != 0 || start.Second() != 0 || start.Nanosecond() != 0 {
		writeV1Error(w, http.StatusBadRequest, "validation_error", "Недопустимое время записи")
		return
	}
	if !s.v1.Reserve(meeting.Booking{Start: start, End: start.Add(time.Duration(item.Duration) * time.Minute), MeetingTypeID: item.ID, GuestName: strings.TrimSpace(req.GuestName), GuestEmail: req.GuestEmail}) {
		writeV1Error(w, http.StatusConflict, "booking_conflict", "Это время уже занято")
		return
	}
	writeJSON(w, http.StatusCreated, generated.BookingConfirmation{MeetingType: toV1(item), Start: start.Format(time.RFC3339)})
}

func (s server) ListUpcomingBookings(w http.ResponseWriter, _ *http.Request) {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		writeV1Error(w, http.StatusInternalServerError, "timezone_error", "Не удалось определить часовой пояс")
		return
	}
	bookings := s.v1.Bookings()
	slices.SortFunc(bookings, func(a, b meeting.Booking) int { return a.Start.Compare(b.Start) })
	items := make([]generated.UpcomingBooking, 0)
	now := s.now()
	types := s.v1.List()
	for _, booking := range bookings {
		if !booking.End.After(now) {
			continue
		}
		for _, item := range types {
			if item.ID == booking.MeetingTypeID {
				items = append(items, generated.UpcomingBooking{MeetingType: toV1(item), Start: booking.Start.In(location).Format(time.RFC3339), GuestName: booking.GuestName, GuestEmail: booking.GuestEmail})
				break
			}
		}
	}
	writeJSON(w, http.StatusOK, generated.UpcomingBookingList{Items: items})
}
