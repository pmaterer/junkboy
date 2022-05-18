package junkboy

import (
	"context"
	"net/http"
	"regexp"
	"strings"
)

type Router struct {
	pathPrefix string
	routes     []route
}

func NewRouter(pathPrefix string) *Router {
	return &Router{
		pathPrefix: strings.TrimSuffix(pathPrefix, "/"),
	}
}

func (rt *Router) AddRoute(methods []string, pattern string, handler http.HandlerFunc) {
	newRt := newRoute(methods, rt.pathPrefix+pattern, handler)
	rt.routes = append(rt.routes, newRt)
}

type ctxKey struct{}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var allow []string

	for _, route := range rt.routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
      for _, method := range route.methods {
        if r.Method != method {
          allow = append(allow, method)
          continue
        }
      }
      //if r.Method != route.method {
			//	allow = append(allow, route.method)
			//	continue
			//}

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
	methods  []string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func newRoute(methods []string, pattern string, handler http.HandlerFunc) route {
	return route{methods, regexp.MustCompile("^" + pattern + "$"), handler}
}

func getField(r *http.Request, index int) string {
	//nolint:errcheck // Need to figure out what error value can be.
	fields := r.Context().Value(ctxKey{}).([]string)

	return fields[index]
}
