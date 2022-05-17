package junkboy

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type LoggingMiddleware struct {
	handler http.Handler
}

func (h *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	h.handler.ServeHTTP(w, r)
	log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
}

func NewLoggingMiddleware(handler http.Handler) *LoggingMiddleware {
	return &LoggingMiddleware{handler}
}

// type LogWriter struct {
// 	http.ResponseWriter
// }

// func (w LogWriter) Write(p []byte) (n int, err error) {
// 	n, err = w.ResponseWriter.Write(p)
// 	if err != nil {
// 		log.Printf("Write failed: %v", err)
// 	}
// 	return
// }

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

type Router struct {
	pathPrefix string
	routes     []route
}

func NewRouter(pathPrefix string) *Router {
	return &Router{
		pathPrefix: strings.TrimSuffix(pathPrefix, "/"),
	}
}

func (rt *Router) AddRoute(method, pattern string, handler http.HandlerFunc) {
	newRt := newRoute(method, rt.pathPrefix+pattern, handler)
	rt.routes = append(rt.routes, newRt)
}

type ctxKey struct{}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var allow []string

	for _, route := range rt.routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}

			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.handler(w, r.WithContext(ctx))

			return
		}
	}

	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)

		return
	}

	http.NotFound(w, r)
}

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

func getField(r *http.Request, index int) string {
	//nolint:errcheck // Need to figure out what error value can be.
	fields := r.Context().Value(ctxKey{}).([]string)

	return fields[index]
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
