package main

import (
	"fmt"
	"log"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request URI: %s", r.RequestURI)
		next.ServeHTTP(w, r)
	}
}

func applyMiddlewares(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

type Router struct {
	routes map[string]http.HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]http.HandlerFunc),
	}
}

func (r *Router) HandleFunc(pattern string, handler http.HandlerFunc) {
	r.routes[pattern] = handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if handler, ok := r.routes[req.URL.Path]; ok {
		handler(w, req)
	} else {
		http.NotFound(w, req)
	}
}

func main() {
	router := NewRouter()

	router.HandleFunc("/", applyMiddlewares(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	}, loggingMiddleware))

	router.HandleFunc("/about", applyMiddlewares(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is the about page!")
	}, loggingMiddleware))

	http.ListenAndServe(":8080", router)
}
