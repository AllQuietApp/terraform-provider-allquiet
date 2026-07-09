package provider

import (
	"net/http"
	"testing"
)

func TestIsCreateSuccessStatus(t *testing.T) {
	tests := []struct {
		statusCode int
		want       bool
	}{
		{http.StatusOK, true},
		{http.StatusCreated, true},
		{http.StatusNoContent, false},
		{http.StatusNotFound, false},
		{http.StatusInternalServerError, false},
	}

	for _, tt := range tests {
		if got := isCreateSuccessStatus(tt.statusCode); got != tt.want {
			t.Errorf("isCreateSuccessStatus(%d) = %v, want %v", tt.statusCode, got, tt.want)
		}
	}
}

func TestEnsureCreateSuccess(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "https://api.allquiet.app/api/public/v1/team", nil)

	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{"200 OK", http.StatusOK, false},
		{"201 Created", http.StatusCreated, false},
		{"204 No Content", http.StatusNoContent, true},
		{"404 Not Found", http.StatusNotFound, true},
		{"500 Internal Server Error", http.StatusInternalServerError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Request:    req,
			}

			err := ensureCreateSuccess(resp, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("ensureCreateSuccess() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsDeleteSuccessStatus(t *testing.T) {
	tests := []struct {
		statusCode int
		want       bool
	}{
		{http.StatusOK, true},
		{http.StatusNoContent, true},
		{http.StatusCreated, false},
		{http.StatusNotFound, false},
		{http.StatusInternalServerError, false},
	}

	for _, tt := range tests {
		if got := isDeleteSuccessStatus(tt.statusCode); got != tt.want {
			t.Errorf("isDeleteSuccessStatus(%d) = %v, want %v", tt.statusCode, got, tt.want)
		}
	}
}

func TestEnsureDeleteSuccess(t *testing.T) {
	req, _ := http.NewRequest(http.MethodDelete, "https://api.allquiet.app/api/public/v1/team/123", nil)

	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{"200 OK", http.StatusOK, false},
		{"204 No Content", http.StatusNoContent, false},
		{"201 Created", http.StatusCreated, true},
		{"404 Not Found", http.StatusNotFound, true},
		{"500 Internal Server Error", http.StatusInternalServerError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Request:    req,
			}

			err := ensureDeleteSuccess(resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("ensureDeleteSuccess() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
