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

func TestUpcomingBookings(t *testing.T) {
	now := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)
	h := server.NewWithClock(func() time.Time { return now })
	post := func(path, body string) {
		t.Helper()
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequestWithContext(context.Background(), http.MethodPost, path, strings.NewReader(body)))
		if w.Code != http.StatusCreated {
			t.Fatalf("POST %s: %d %s", path, w.Code, w.Body.String())
		}
	}
	post("/api/v1/meeting-types", `{"id":"intro","name":"Intro","description":"","durationMinutes":60}`)
	post("/api/v1/meeting-types", `{"id":"other","name":"Other","description":"","durationMinutes":30}`)
	post("/api/v1/bookings", `{"meetingTypeId":"other","start":"2026-11-16T12:00:00+03:00","guestName":"B","guestEmail":"b@example.com"}`)
	post("/api/v1/bookings", `{"meetingTypeId":"intro","start":"2026-11-16T10:00:00+03:00","guestName":"A","guestEmail":"a@example.com"}`)
	get := func() []struct {
		MeetingType struct {
			Name string `json:"name"`
		} `json:"meetingType"`
		Start      string `json:"start"`
		GuestName  string `json:"guestName"`
		GuestEmail string `json:"guestEmail"`
	} {
		t.Helper()
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/bookings", nil))
		if w.Code != http.StatusOK {
			t.Fatalf("GET bookings: %d %s", w.Code, w.Body.String())
		}
		var result struct {
			Items []struct {
				MeetingType struct {
					Name string `json:"name"`
				} `json:"meetingType"`
				Start      string `json:"start"`
				GuestName  string `json:"guestName"`
				GuestEmail string `json:"guestEmail"`
			} `json:"items"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		return result.Items
	}
	items := get()
	if len(items) != 2 || items[0].MeetingType.Name != "Intro" || items[0].Start != "2026-11-16T10:00:00+03:00" || items[0].GuestName != "A" || items[0].GuestEmail != "a@example.com" || items[1].MeetingType.Name != "Other" {
		t.Fatalf("unexpected bookings: %+v", items)
	}
	now = time.Date(2026, 11, 16, 8, 0, 0, 0, time.UTC)
	if items = get(); len(items) != 1 || items[0].GuestName != "B" {
		t.Fatalf("completed booking remained: %+v", items)
	}
}
