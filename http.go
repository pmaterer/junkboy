package junkboy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"strings"
	"time"
)

type LoggingMiddleware struct {
	handler http.Handler
}

func NewLoggingMiddleware(handler http.Handler) *LoggingMiddleware {
	return &LoggingMiddleware{handler}
}

func (h *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	h.handler.ServeHTTP(w, r)
	log.Printf("%s %s in %s", r.Method, r.RequestURI, time.Since(start))
}

type CorsMiddleware struct {
  handler http.Handler
}

func NewCorsMiddleware(handler http.Handler) *CorsMiddleware {
  return &CorsMiddleware{handler}
}


func (h *CorsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Access-Control-Allow-Headers", "*")
  h.handler.ServeHTTP(w, r)
}

func readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		// We triage the error.
		var syntaxError *json.SyntaxError

		var unmarshalTypeError *json.UnmarshalTypeError

		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		// In some circumstances `Decode()` may return an `io.ErrUnexpectedEOF` error
		// for syntax errors in the JSON. See: https://github.com/golang/go/issues/25956
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		// JSON value is the wrong type for the target.
		case errors.As(err, &unmarshalTypeError):
			return fmt.Errorf("body contains incorrect JSON type field %q (at character %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)

		// This is returned if the request body is empty.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		// This will be returned if we pass a non-nil pointer to `Decode()`.
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	// Ensure request body contains a single JSON value.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func logerr(n int, err error) {
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	logerr(w.Write(js))
}

func writeError(w http.ResponseWriter, status int, message string) {
	errorResponse := ErrorResponse{
		Status:  status,
		Message: message,
	}

	writeJSON(w, status, errorResponse)
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func contentTypeIsValid(w http.ResponseWriter, r *http.Request, expectedContentType string) bool {
	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)

	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return false
	}

	if mediaType != expectedContentType {
		writeError(w, http.StatusUnsupportedMediaType, fmt.Sprintf("expected '%s' Content-Type, got '%s'", expectedContentType, mediaType))
		return false
	}

	return true
}
