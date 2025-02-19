package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strings"

	"github.com/paskozdilar/tcontext"
)

// Define request data struct
type RequestData struct {
	RequestID string
	UserID    string
}

// Middleware to add request ID to the context
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := fmt.Sprintf("%d", rand.Int())
		ctx := tcontext.WithData(r.Context(), &RequestData{
			RequestID: requestID,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Handler to greet the user
func helloHandler(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimPrefix(r.URL.Path, "/")
	if userID == "" {
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}

	// NOTE:
	// Use tcontext.FromContext(r.Context()) if not sure if context is really
	// a tcontext.Context
	ctx := r.Context().(tcontext.Context[*RequestData])
	ctx.Data().UserID = userID

	// Extract request ID from context
	requestID := ctx.Data().RequestID

	// Log the request
	log.Printf("Request ID: %s, User ID: %s", requestID, userID)

	// Respond to the client
	fmt.Fprintf(w, "Hello, %s!\n", userID)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/hello/", http.StripPrefix("/hello", http.HandlerFunc(helloHandler)))

	// Wrap the mux with the request ID middleware
	handler := requestIDMiddleware(mux)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
