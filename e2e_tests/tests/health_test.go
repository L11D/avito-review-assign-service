package tests

import (
	"net/http"
	"os"
	"testing"
)

func TestHealth(t *testing.T) {
	url := os.Getenv("API_URL") + "/health"

	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
