package server_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/godovdm/ai-for-developers-project-386/backend/internal/server"
)

func TestV1MeetingTypes(t *testing.T) {
	handler := server.New()
	request := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/meeting-types", strings.NewReader(`{"id":"consultation","name":"Консультация","description":"","durationMinutes":30}`))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("create: status %d, body %s", response.Code, response.Body.String())
	}
	request = httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/meeting-types", nil)
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"id":"consultation"`) || !strings.Contains(response.Body.String(), `"items"`) {
		t.Fatalf("list: status %d, body %s", response.Code, response.Body.String())
	}
}

func TestV1MeetingTypeValidation(t *testing.T) {
	handler := server.New()
	cases := []struct{ name, body string }{
		{"missing id", `{"name":"N","description":"","durationMinutes":30}`},
		{"null description", `{"id":"valid","name":"N","description":null,"durationMinutes":30}`},
		{"unknown field", `{"id":"valid","name":"N","description":"","durationMinutes":30,"extra":true}`},
		{"invalid id", `{"id":"Bad-id","name":"N","description":"","durationMinutes":30}`},
		{"empty name", `{"id":"valid","name":"  ","description":"","durationMinutes":30}`},
		{"invalid duration", `{"id":"valid","name":"N","description":"","durationMinutes":45}`},
		{"invalid json", `{`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/meeting-types", strings.NewReader(tc.body))
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != http.StatusBadRequest || !strings.Contains(response.Body.String(), `"code":"validation_error"`) {
				t.Fatalf("status %d, body %s", response.Code, response.Body.String())
			}
		})
	}
	body := `{"id":"intro","name":" Name ","description":" Detail ","durationMinutes":60}`
	for i, expected := range []int{http.StatusCreated, http.StatusConflict} {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/meeting-types", strings.NewReader(body)))
		if response.Code != expected {
			t.Fatalf("post %d: status %d, body %s", i, response.Code, response.Body.String())
		}
	}
	for _, tc := range []struct {
		path     string
		status   int
		contains string
	}{
		{"/api/v1/meeting-types/intro", 200, `"name":"Name"`},
		{"/api/v1/meeting-types/missing", 404, `"code":"meeting_type_not_found"`},
	} {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequestWithContext(context.Background(), http.MethodGet, tc.path, nil))
		if response.Code != tc.status || !strings.Contains(response.Body.String(), tc.contains) {
			t.Fatalf("%s: status %d, body %s", tc.path, response.Code, response.Body.String())
		}
	}
}

func TestMeetingTypeVersionsShareStorage(t *testing.T) {
	handler := server.New()
	cases := []struct{ path, body string }{
		{"/api/meeting-types", `{"name":"Legacy","description":"Description","durationMinutes":30}`},
		{"/api/v1/meeting-types", `{"id":"new","name":"New","description":"","durationMinutes":60}`},
	}
	for _, tc := range cases {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequestWithContext(context.Background(), http.MethodPost, tc.path, strings.NewReader(tc.body)))
		if response.Code != http.StatusCreated {
			t.Fatalf("%s: status %d, body %s", tc.path, response.Code, response.Body.String())
		}
	}
	for _, tc := range cases {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequestWithContext(context.Background(), http.MethodGet, tc.path, nil))
		if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"name":"Legacy"`) || !strings.Contains(response.Body.String(), `"name":"New"`) {
			t.Fatalf("%s: status %d, body %s", tc.path, response.Code, response.Body.String())
		}
	}
}
