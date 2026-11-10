package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestSPAServesIndexForClientRoutes(t *testing.T) {
	dir := t.TempDir()
	index := "<!doctype html><title>Calendar</title>"
	writeFile(t, filepath.Join(dir, "index.html"), index)
	writeFile(t, filepath.Join(dir, "assets", "app.js"), "console.log('app')")

	handler := SPA(dir)

	for _, target := range []string{"/", "/owner", "/meeting-types/42"} {
		rec := serve(handler, target)
		if rec.Code != http.StatusOK {
			t.Errorf("%s status = %d, want %d", target, rec.Code, http.StatusOK)
		}
		if rec.Body.String() != index {
			t.Errorf("%s body = %q, want index.html", target, rec.Body.String())
		}
	}

	rec := serve(handler, "/assets/app.js")
	if rec.Code != http.StatusOK {
		t.Errorf("asset status = %d, want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != "console.log('app')" {
		t.Errorf("asset body = %q, want app.js", rec.Body.String())
	}
}

func TestSPAReturnsNotFoundForMissingAssets(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "index.html"), "<!doctype html>")

	rec := serve(SPA(dir), "/assets/missing.js")
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func serve(handler http.Handler, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequestWithContext(context.Background(), http.MethodGet, target, nil))
	return rec
}

func writeFile(t *testing.T, name, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(name), 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(name, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}
