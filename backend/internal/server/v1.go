package server

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
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
