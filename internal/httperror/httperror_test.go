package utility

import (
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	err := New(http.StatusBadRequest, "bad request")
	if err == nil {
		t.Fatal("New returned nil")
	}
	if err.Code != http.StatusBadRequest {
		t.Errorf("expected code %d, got %d", http.StatusBadRequest, err.Code)
	}
	if err.Error() != "bad request" {
		t.Errorf("expected message 'bad request', got '%s'", err.Error())
	}
}

func TestNewWithDifferentCodes(t *testing.T) {
	tests := []struct {
		code    int
		message string
	}{
		{http.StatusNotFound, "not found"},
		{http.StatusInternalServerError, "internal error"},
		{http.StatusForbidden, "forbidden"},
		{http.StatusUnauthorized, "unauthorized"},
	}

	for _, tt := range tests {
		err := New(tt.code, tt.message)
		if err.Code != tt.code {
			t.Errorf("expected code %d, got %d", tt.code, err.Code)
		}
		if err.Error() != tt.message {
			t.Errorf("expected message '%s', got '%s'", tt.message, err.Error())
		}
	}
}

func TestNewEmptyMessage(t *testing.T) {
	err := New(http.StatusOK, "")
	if err == nil {
		t.Fatal("New returned nil")
	}
	if err.Error() != "" {
		t.Errorf("expected empty message, got '%s'", err.Error())
	}
}

func TestNotFound(t *testing.T) {
	err := NotFound("user not found")
	if err == nil {
		t.Fatal("NotFound returned nil")
	}
	if err.Code != http.StatusNotFound {
		t.Errorf("expected code %d, got %d", http.StatusNotFound, err.Code)
	}
	if err.Error() != "user not found" {
		t.Errorf("expected message 'user not found', got '%s'", err.Error())
	}
}

func TestNotFoundIsHTTPError(t *testing.T) {
	err := NotFound("test")
	if _, ok := interface{}(err).(*HTTPError); !ok {
		t.Error("NotFound did not return an *HTTPError")
	}
}

func TestErrorImplementsErrorInterface(t *testing.T) {
	err := New(http.StatusTeapot, "I'm a teapot")
	var e error = err
	if e.Error() != "I'm a teapot" {
		t.Errorf("expected 'I'm a teapot', got '%s'", e.Error())
	}
}
