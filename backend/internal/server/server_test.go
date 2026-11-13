package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/godovdm/ai-for-developers-project-386/backend/internal/api/generated"
)

func TestHealth(t *testing.T) {
	rec := request(t, New(), http.MethodGet, "/api/health", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var body generated.HealthStatus
	decode(t, rec, &body)
	if body.Status != "ok" {
		t.Errorf("status = %q, want ok", body.Status)
	}
}

func TestListMeetingTypesEmpty(t *testing.T) {
	rec := request(t, New(), http.MethodGet, "/api/meeting-types", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != "[]" {
		t.Errorf("body = %s, want []", got)
	}
}

func TestCreateAndListMeetingTypes(t *testing.T) {
	handler := New()
	rec := request(t, handler, http.MethodPost, "/api/meeting-types",
		`{"name":"  Вводный звонок  ","description":"  Знакомство  ","durationMinutes":60}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusCreated, rec.Body.String())
	}
	if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", contentType)
	}

	var created generated.MeetingType
	decode(t, rec, &created)
	if created.Id == "" {
		t.Error("id is empty")
	}
	if created.Name != "Вводный звонок" {
		t.Errorf("name = %q, want trimmed name", created.Name)
	}
	if created.Description != "Знакомство" {
		t.Errorf("description = %q, want trimmed description", created.Description)
	}
	if created.DurationMinutes != 60 {
		t.Errorf("durationMinutes = %v, want 60", created.DurationMinutes)
	}

	rec = request(t, handler, http.MethodGet, "/api/meeting-types", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d", rec.Code, http.StatusOK)
	}
	var list []generated.MeetingType
	decode(t, rec, &list)
	if len(list) != 1 {
		t.Fatalf("list length = %d, want 1", len(list))
	}
	if list[0] != created {
		t.Errorf("list[0] = %+v, want %+v", list[0], created)
	}
}

func TestCreateMeetingTypeRejectsInvalidData(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "empty name", body: `{"name":"","description":"Описание","durationMinutes":30}`},
		{name: "whitespace name", body: `{"name":"   ","description":"Описание","durationMinutes":30}`},
		{name: "name too long", body: bodyWith(t, map[string]any{
			"name":            strings.Repeat("а", 101),
			"description":     "Описание",
			"durationMinutes": 30,
		})},
		{name: "empty description", body: `{"name":"Звонок","description":"","durationMinutes":30}`},
		{name: "description too long", body: bodyWith(t, map[string]any{
			"name":            "Звонок",
			"description":     strings.Repeat("а", 1001),
			"durationMinutes": 30,
		})},
		{name: "invalid duration", body: `{"name":"Звонок","description":"Описание","durationMinutes":45}`},
		{name: "fractional duration", body: `{"name":"Звонок","description":"Описание","durationMinutes":30.5}`},
		{name: "tiny fractional duration", body: `{"name":"Звонок","description":"Описание","durationMinutes":30.0000001}`},
		{name: "sub-float fractional duration", body: `{"name":"Звонок","description":"Описание","durationMinutes":30.0000000000000001}`},
		{name: "small exponent duration", body: `{"name":"Звонок","description":"Описание","durationMinutes":1e1}`},
		{name: "fractional exponent duration", body: `{"name":"Звонок","description":"Описание","durationMinutes":4.5e1}`},
		{name: "string duration", body: `{"name":"Звонок","description":"Описание","durationMinutes":"60"}`},
		{name: "boolean duration", body: `{"name":"Звонок","description":"Описание","durationMinutes":true}`},
		{name: "null duration", body: `{"name":"Звонок","description":"Описание","durationMinutes":null}`},
		{name: "missing duration", body: `{"name":"Звонок","description":"Описание"}`},
		{name: "trailing garbage", body: `{"name":"Звонок","description":"Описание","durationMinutes":60} garbage`},
		{name: "two json objects", body: `{"name":"Звонок","description":"Описание","durationMinutes":60}{"name":"Второй","description":"Описание","durationMinutes":30}`},
		{name: "broken json", body: `{"name":`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := New()
			rec := request(t, handler, http.MethodPost, "/api/meeting-types", test.body)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusBadRequest, rec.Body.String())
			}
			var apiError generated.ApiError
			decode(t, rec, &apiError)
			if apiError.Message == "" {
				t.Error("error message is empty")
			}

			rec = request(t, handler, http.MethodGet, "/api/meeting-types", "")
			if got := strings.TrimSpace(rec.Body.String()); got != "[]" {
				t.Errorf("state changed: list body = %s, want []", got)
			}
		})
	}
}

func TestCreateMeetingTypeAcceptsValidDurations(t *testing.T) {
	tests := []struct {
		name     string
		duration string
		want     int
	}{
		{name: "30", duration: "30", want: 30},
		{name: "60", duration: "60", want: 60},
		{name: "90", duration: "90", want: 90},
		{name: "120", duration: "120", want: 120},
		{name: "60.0", duration: "60.0", want: 60},
		{name: "6e1", duration: "6e1", want: 60},
	}

	handler := New()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := `{"name":"Звонок","description":"Описание","durationMinutes":` + test.duration + `}`
			rec := request(t, handler, http.MethodPost, "/api/meeting-types", body)
			if rec.Code != http.StatusCreated {
				t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusCreated, rec.Body.String())
			}
			var created generated.MeetingType
			decode(t, rec, &created)
			if created.DurationMinutes != generated.MeetingDuration(test.want) {
				t.Errorf("durationMinutes = %v, want %d", created.DurationMinutes, test.want)
			}
		})
	}

	rec := request(t, handler, http.MethodGet, "/api/meeting-types", "")
	var list []generated.MeetingType
	decode(t, rec, &list)
	if len(list) != len(tests) {
		t.Errorf("list length = %d, want %d", len(list), len(tests))
	}
}

func TestCreateMeetingTypeKeepsStateOnInvalidDuration(t *testing.T) {
	handler := New()
	rec := request(t, handler, http.MethodPost, "/api/meeting-types",
		`{"name":"Первый","description":"Описание","durationMinutes":30}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("seed status = %d, want %d", rec.Code, http.StatusCreated)
	}

	for _, body := range []string{
		`{"name":"Второй","description":"Описание","durationMinutes":30.0000001}`,
		`{"name":"Второй","description":"Описание","durationMinutes":30.0000000000000001}`,
		`{"name":"Второй","description":"Описание","durationMinutes":30.5}`,
		`{"name":"Второй","description":"Описание","durationMinutes":"60"}`,
		`{"name":"Второй","description":"Описание","durationMinutes":true}`,
		`{"name":"Второй","description":"Описание","durationMinutes":null}`,
		`{"name":"Второй","description":"Описание"}`,
		`{"name":"Второй","description":"Описание","durationMinutes":60} garbage`,
		`{"name":"Второй","description":"Описание","durationMinutes":60}{"name":"Третий","description":"Описание","durationMinutes":30}`,
	} {
		rec = request(t, handler, http.MethodPost, "/api/meeting-types", body)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusBadRequest, rec.Body.String())
		}
		var apiError generated.ApiError
		decode(t, rec, &apiError)
		if apiError.Message == "" {
			t.Error("error message is empty")
		}
	}

	rec = request(t, handler, http.MethodGet, "/api/meeting-types", "")
	var list []generated.MeetingType
	decode(t, rec, &list)
	if len(list) != 1 {
		t.Fatalf("list length = %d, want 1", len(list))
	}
	if list[0].Name != "Первый" {
		t.Errorf("list[0].Name = %q, want Первый", list[0].Name)
	}
}

func request(t *testing.T, handler http.Handler, method, target, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequestWithContext(context.Background(), method, target, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func decode(t *testing.T, rec *httptest.ResponseRecorder, target any) {
	t.Helper()
	if err := json.NewDecoder(rec.Body).Decode(target); err != nil {
		t.Fatalf("decode response %q: %v", rec.Body.String(), err)
	}
}

func bodyWith(t *testing.T, value any) string {
	t.Helper()
	body, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}
	return string(body)
}
