package junkboy

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockAnchorService struct {
	AddAnchorFunc    func(a Anchor) (int, error)
	UpdateAnchorFunc func(a Anchor) error
	GetAnchorFunc    func(id int) (Anchor, error)
	GetAnchorsFunc   func() ([]Anchor, error)
	DeleteAnchorFunc func(id int) error
}

func (ar *mockAnchorService) AddAnchor(a Anchor) (int, error)  { return ar.AddAnchorFunc(a) }
func (ar *mockAnchorService) UpdateAnchor(a Anchor) error      { return ar.UpdateAnchorFunc(a) }
func (ar *mockAnchorService) GetAnchor(id int) (Anchor, error) { return ar.GetAnchorFunc(id) }
func (ar *mockAnchorService) GetAnchors() ([]Anchor, error)    { return ar.GetAnchorsFunc() }
func (ar *mockAnchorService) DeleteAnchor(id int) error        { return ar.DeleteAnchorFunc(id) }

var anchorJSON = []byte(`{"id":1,"url":"https://example.com"}`)

var anchorsJSON = []byte(`[{"id":2,"url":"https://example.com/a"},{"id":3,"url":"https://example.com/b"},{"id":4,"url":"https://example.com/c"}]`)

func TestAddAnchorHandler(t *testing.T) {
	tests := []struct {
		name           string
		reqBody        []byte
		contentType    string
		method         func(a Anchor) (int, error)
		expectedStatus int
		responseBody   string
	}{
		{
			name:           "Add anchor ok",
			reqBody:        anchorJSON,
			contentType:    "application/json",
			method:         func(a Anchor) (int, error) { return 1, nil },
			expectedStatus: http.StatusCreated,
			responseBody:   `{"id":1}`,
		},
		{
			name:           "Add anchor bad content-type",
			reqBody:        anchorJSON,
			contentType:    "application/junk",
			method:         func(a Anchor) (int, error) { return 1, nil },
			expectedStatus: http.StatusUnsupportedMediaType,
			responseBody:   `{"status":415,"message":"expected 'application/json' Content-Type, got 'application/junk'"}`,
		},
		{
			name:           "Add anchor bad body",
			reqBody:        []byte(`{"url":false}`),
			contentType:    "application/json",
			method:         func(a Anchor) (int, error) { return 1, nil },
			expectedStatus: http.StatusBadRequest,
			responseBody:   `{"status":400,"message":"json: cannot unmarshal bool into Go struct field Anchor.url of type string"}`,
		},
		{
			name:           "Add anchor internal error",
			reqBody:        anchorJSON,
			contentType:    "application/json",
			method:         func(a Anchor) (int, error) { return 0, errors.New("internal error") },
			expectedStatus: http.StatusInternalServerError,
			responseBody:   `{"status":500,"message":"internal error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &mockAnchorService{AddAnchorFunc: tt.method}
			anchorHandler := NewAnchorHTTPHandler(s)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(anchorHandler.addAnchorHandler)

			req, err := http.NewRequest(http.MethodPost, "/anchor", bytes.NewBuffer(tt.reqBody))
			assertNoError(t, err)
			req.Header.Add("Content-Type", tt.contentType)

			handler.ServeHTTP(rr, req)
			assertEqual(t, tt.expectedStatus, rr.Code)
			assertEqual(t, tt.responseBody, rr.Body.String())
		})
	}
}

func TestGetAnchorsHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         func() ([]Anchor, error)
		expectedStatus int
		responseBody   []byte
	}{
		{
			name:           "Get anchors",
			method:         func() ([]Anchor, error) { return testAnchors, nil },
			expectedStatus: http.StatusOK,
			responseBody:   anchorsJSON,
		},
		{
			name:           "Get anchors internal error",
			method:         func() ([]Anchor, error) { return nil, errors.New("internal server error") },
			expectedStatus: http.StatusInternalServerError,
			responseBody:   []byte(`{"status":500,"message":"internal server error"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &mockAnchorService{GetAnchorsFunc: tt.method}
			anchorHandler := NewAnchorHTTPHandler(s)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(anchorHandler.getAnchorsHandler)

			req, err := http.NewRequest(http.MethodGet, "/anchors", nil)
			assertNoError(t, err)

			handler.ServeHTTP(rr, req)

			assertEqual(t, tt.expectedStatus, rr.Code)
			assertBytesEqual(t, tt.responseBody, rr.Body.Bytes())
		})
	}
}

func TestGetAnchorHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         func(id int) (Anchor, error)
		expectedStatus int
		responseBody   []byte
		pathId         string
	}{
		{
			name:           "Get anchor ok",
			method:         func(id int) (Anchor, error) { return testAnchor, nil },
			expectedStatus: http.StatusOK,
			responseBody:   anchorJSON,
			pathId:         "1",
		},
		{
			name:           "Get anchor server error",
			method:         func(id int) (Anchor, error) { return testAnchor, errors.New("internal server error") },
			expectedStatus: http.StatusInternalServerError,
			responseBody:   []byte(`{"status":500,"message":"internal server error"}`),
			pathId:         "1",
		},
		{
			name:           "Get anchor bad id",
			method:         func(id int) (Anchor, error) { return testAnchor, errors.New("internal server error") },
			expectedStatus: http.StatusBadRequest,
			responseBody:   []byte(`{"status":400,"message":"invalid anchor id 'abc'"}`),
			pathId:         "abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &mockAnchorService{GetAnchorFunc: tt.method}
			anchorHandler := NewAnchorHTTPHandler(s)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(anchorHandler.getAnchorHandler)

			req, err := http.NewRequest(http.MethodGet, "/anchor/1", nil)
			assertNoError(t, err)
			ctx := context.WithValue(req.Context(), ctxKey{}, []string{tt.pathId})

			handler.ServeHTTP(rr, req.WithContext(ctx))

			assertEqual(t, tt.expectedStatus, rr.Code)
			assertBytesEqual(t, tt.responseBody, rr.Body.Bytes())
		})
	}
}

func TestUpdateAnchorHandler(t *testing.T) {
	tests := []struct {
		name           string
		reqBody        []byte
		contentType    string
		method         func(a Anchor) error
		responseBody   []byte
		expectedStatus int
	}{
		{
			name:           "Update anchor ok",
			reqBody:        anchorJSON,
			contentType:    "application/json",
			method:         func(a Anchor) error { return nil },
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Update anchor server error",
			reqBody:        anchorJSON,
			contentType:    "application/json",
			method:         func(a Anchor) error { return errors.New("internal server error") },
			expectedStatus: http.StatusInternalServerError,
			responseBody:   []byte(`{"status":500,"message":"internal server error"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &mockAnchorService{UpdateAnchorFunc: tt.method}
			anchorHandler := NewAnchorHTTPHandler(s)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(anchorHandler.updateAnchorHandler)

			req, err := http.NewRequest(http.MethodPut, "/anchor/1", bytes.NewBuffer(tt.reqBody))
			assertNoError(t, err)
			req.Header.Add("Content-Type", tt.contentType)

			handler.ServeHTTP(rr, req)
			assertEqual(t, tt.expectedStatus, rr.Code)
			assertBytesEqual(t, tt.responseBody, rr.Body.Bytes())
		})
	}
}

func TestDeleteAnchorHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         func(id int) error
		expectedStatus int
		responseBody   []byte
		pathId         string
	}{
		{
			name:           "Delete anchor ok",
			method:         func(id int) error { return nil },
			expectedStatus: http.StatusNoContent,
			pathId:         "1",
		},
		{
			name:           "Delete anchor server error",
			method:         func(id int) error { return errors.New("internal server error") },
			expectedStatus: http.StatusInternalServerError,
			pathId:         "1",
			responseBody:   []byte(`{"status":500,"message":"internal server error"}`),
		},
		{
			name:           "Delete anchor bad id",
			method:         func(id int) error { return nil },
			expectedStatus: http.StatusBadRequest,
			pathId:         "abc",
			responseBody:   []byte(`{"status":400,"message":"invalid anchor id 'abc'"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &mockAnchorService{DeleteAnchorFunc: tt.method}
			anchorHandler := NewAnchorHTTPHandler(s)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(anchorHandler.deleteAnchorHandler)

			req, err := http.NewRequest(http.MethodDelete, "/anchor/1", nil)
			assertNoError(t, err)
			ctx := context.WithValue(req.Context(), ctxKey{}, []string{tt.pathId})

			handler.ServeHTTP(rr, req.WithContext(ctx))

			assertEqual(t, tt.expectedStatus, rr.Code)
			assertBytesEqual(t, tt.responseBody, rr.Body.Bytes())
		})
	}
}
