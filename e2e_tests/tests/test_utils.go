package tests

// test_utils.go

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func MakeJSONRequest(t *testing.T, method, url string, requestBody any) (*http.Response, []byte) {
	t.Helper()

	var (
		bodyBytes []byte
		err       error
	)

	if requestBody != nil {
		bodyBytes, err = json.Marshal(requestBody)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	return resp, responseBody
}

func ParseJSONResponse(t *testing.T, body []byte, target any) {
	t.Helper()

	if err := json.Unmarshal(body, target); err != nil {
		t.Fatalf("failed to unmarshal response: %v. Body: %s", err, string(body))
	}
}

func AssertStatusCode(t *testing.T, resp *http.Response, expected int) {
	t.Helper()

	if resp.StatusCode != expected {
		t.Fatalf("expected status %d, got %d", expected, resp.StatusCode)
	}
}

func MakeQueryRequest(t *testing.T, method, baseURL string, queryParams map[string]string) (*http.Response, []byte) {
	t.Helper()

	url, err := url.Parse(baseURL)
	if err != nil {
		t.Fatalf("failed to parse URL: %v", err)
	}

	if queryParams != nil {
		query := url.Query()
		for key, value := range queryParams {
			query.Add(key, value)
		}

		url.RawQuery = query.Encode()
	}

	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	return resp, responseBody
}

func GetBoolPtr(b bool) *bool {
	return &b
}
