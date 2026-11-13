package server_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/godovdm/ai-for-developers-project-386/backend/internal/server"
)

func TestAvailabilityWindow(t *testing.T) {
	now := time.Date(2026, 11, 13, 20, 59, 0, 0, time.UTC)
	handler := server.NewWithClock(func() time.Time { return now })
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/meeting-types", strings.NewReader(`{"id":"intro","name":"Intro","description":"","durationMinutes":60}`)))
	if response.Code != http.StatusCreated {
		t.Fatal(response.Body.String())
	}
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/meeting-types/intro/availability", nil))
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `2026-11-16T10:00:00+03:00`) || strings.Contains(response.Body.String(), `2026-11-16T17:30:00+03:00`) {
		t.Fatalf("status %d: %s", response.Code, response.Body.String())
	}
}

func TestAvailabilityCalendarBoundaries(t *testing.T) {
	now := time.Date(2026, 11, 13, 20, 59, 0, 0, time.UTC)
	handler := server.NewWithClock(func() time.Time { return now })
	for _, body := range []string{`{"id":"short","name":"Short","description":"","durationMinutes":30}`, `{"id":"long","name":"Long","description":"","durationMinutes":120}`} {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/meeting-types", strings.NewReader(body)))
		if response.Code != http.StatusCreated {
			t.Fatal(response.Body.String())
		}
	}
	get := func(id string) struct {
		Days []struct {
			Date  string   `json:"date"`
			Slots []string `json:"slots"`
		} `json:"days"`
	} {
		t.Helper()
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/meeting-types/"+id+"/availability", nil))
		var result struct {
			Days []struct {
				Date  string   `json:"date"`
				Slots []string `json:"slots"`
			} `json:"days"`
		}
		if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		return result
	}
	short := get("short")
	if len(short.Days) != 14 || short.Days[0].Date != "2026-11-14" || short.Days[13].Date != "2026-11-27" || len(short.Days[0].Slots) != 0 || len(short.Days[1].Slots) != 0 {
		t.Fatalf("unexpected window: %+v", short.Days)
	}
	if len(short.Days[2].Slots) != 16 || short.Days[2].Slots[15] != "2026-11-16T17:30:00+03:00" {
		t.Fatalf("unexpected slots: %v", short.Days[2].Slots)
	}
	long := get("long")
	if len(long.Days[2].Slots) != 13 || long.Days[2].Slots[12] != "2026-11-16T16:00:00+03:00" {
		t.Fatalf("unexpected slots: %v", long.Days[2].Slots)
	}
	now = now.Add(time.Minute)
	shifted := get("short")
	if shifted.Days[0].Date != "2026-11-15" || shifted.Days[13].Date != "2026-11-28" {
		t.Fatalf("window did not shift: %+v", shifted.Days)
	}
}

func TestAvailabilityOutsideWindow(t *testing.T) {
	handler := server.NewWithClock(func() time.Time { return time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC) })
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/meeting-types", strings.NewReader(`{"id":"intro","name":"Intro","description":"","durationMinutes":30}`)))
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/meeting-types/intro/availability?from=2026-11-30&to=2026-12-01", nil))
	if response.Code != http.StatusOK || strings.TrimSpace(response.Body.String()) != `{"days":[]}` {
		t.Fatalf("out of window: %d %s", response.Code, response.Body.String())
	}
}
