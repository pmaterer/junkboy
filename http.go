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

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
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

func (r *Router) AddRoute(method, pattern string, handler http.HandlerFunc) {
	route := newRoute(method, r.pathPrefix+pattern, handler)
	r.routes = append(r.routes, route)
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
	fields := r.Context().Value(ctxKey{}).([]string)
	return fields[index]
}

func checkContentType(w http.ResponseWriter, r *http.Request, expectedContentType string) {
	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if mediaType != expectedContentType {
		writeError(w, http.StatusUnsupportedMediaType, fmt.Sprintf("expected '%s' Content-Type, got '%s'", expectedContentType, mediaType))
		return
	}
	return
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
