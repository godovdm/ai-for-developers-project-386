package server_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/godovdm/ai-for-developers-project-386/backend/internal/server"
)

func TestBooking(t *testing.T) {
	h := server.NewWithClock(func() time.Time { return time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC) })
	post := func(path, body string) *httptest.ResponseRecorder {
		t.Helper()
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequestWithContext(context.Background(), http.MethodPost, path, strings.NewReader(body)))
		return w
	}
	post("/api/v1/meeting-types", `{"id":"intro","name":"Intro","description":"","durationMinutes":60}`)
	w := post("/api/v1/bookings", `{"meetingTypeId":"intro","start":"2026-11-16T10:00:00+03:00","guestName":"Anna","guestEmail":"anna@example.com"}`)
	if w.Code != http.StatusCreated || !strings.Contains(w.Body.String(), "Intro") || !strings.Contains(w.Body.String(), "2026-11-16T10:00:00+03:00") {
		t.Fatalf("booking: %d %s", w.Code, w.Body.String())
	}
	w = post("/api/v1/bookings", `{"meetingTypeId":"intro","start":"2026-11-16T10:30:00+03:00","guestName":"Anna","guestEmail":"anna@example.com"}`)
	if w.Code != http.StatusConflict {
		t.Fatalf("overlap: %d %s", w.Code, w.Body.String())
	}
}

func TestBookingBoundariesAndConcurrency(t *testing.T) {
	h := server.NewWithClock(func() time.Time { return time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC) })
	post := func(path, body string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequestWithContext(context.Background(), http.MethodPost, path, strings.NewReader(body)))
		return w
	}
	for _, body := range []string{`{"id":"long","name":"Long","description":"","durationMinutes":60}`, `{"id":"short","name":"Short","description":"","durationMinutes":30}`} {
		if w := post("/api/v1/meeting-types", body); w.Code != 201 {
			t.Fatal(w.Body.String())
		}
	}
	book := func(id, start, name, email string) *httptest.ResponseRecorder {
		body := fmt.Sprintf(`{"meetingTypeId":%q,"start":%q,"guestName":%q,"guestEmail":%q}`, id, start, name, email)
		return post("/api/v1/bookings", body)
	}
	for _, tc := range []struct{ start, name, email string }{
		{"2026-11-16T10:00:00+03:00", "", "a@example.com"}, {"2026-11-16T10:00:00+03:00", "A", "bad"},
		{"2026-11-13T10:00:00+03:00", "A", "a@example.com"}, {"2026-11-14T10:00:00+03:00", "A", "a@example.com"},
		{"2026-11-16T09:30:00+03:00", "A", "a@example.com"}, {"2026-11-16T17:30:00+03:00", "A", "a@example.com"},
		{"2026-11-16T10:15:00+03:00", "A", "a@example.com"},
	} {
		if w := book("long", tc.start, tc.name, tc.email); w.Code != 400 {
			t.Errorf("invalid %v: %d", tc, w.Code)
		}
	}
	if w := book("long", "2026-11-16T10:00:00+03:00", "A", "a@example.com"); w.Code != 201 {
		t.Fatal(w.Body.String())
	}
	for _, start := range []string{"2026-11-16T10:00:00+03:00", "2026-11-16T10:30:00+03:00"} {
		if w := book("short", start, "A", "a@example.com"); w.Code != 409 || !strings.Contains(w.Body.String(), "booking_conflict") {
			t.Fatalf("conflict: %d %s", w.Code, w.Body.String())
		}
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/meeting-types/short/availability", nil))
	if strings.Contains(w.Body.String(), "2026-11-16T10:30:00+03:00") {
		t.Fatal("reserved slot still available")
	}
	if w := book("short", "2026-11-16T11:00:00+03:00", "A", "a@example.com"); w.Code != 201 {
		t.Fatalf("adjacent: %d", w.Code)
	}
	var wg sync.WaitGroup
	codes := make(chan int, 2)
	for range 2 {
		wg.Add(1)
		go func() { defer wg.Done(); codes <- book("long", "2026-11-17T10:00:00+03:00", "A", "a@example.com").Code }()
	}
	wg.Wait()
	close(codes)
	successes, conflicts := 0, 0
	for code := range codes {
		if code == 201 {
			successes++
		}
		if code == 409 {
			conflicts++
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("concurrent successes=%d conflicts=%d", successes, conflicts)
	}
}
