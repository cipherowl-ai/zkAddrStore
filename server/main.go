package main

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

var (
	filter *bloom.BloomFilter
	logger = log.New(os.Stdout, "BloomServer: ", log.LstdFlags)
)

type Response struct {
	Query   string `json:"query"`
	InSet   bool   `json:"in_set"`
	Message string `json:"message"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Println("No .env file found")
	}

	filename := os.Getenv("BLOOM_FILTER_FILE")
	if filename == "" {
		filename = "bloomfilter.gob"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := loadBloomFilter(filename); err != nil {
		logger.Fatalf("Failed to load Bloom filter: %v", err)
	}

	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.Handle("/check", rateLimitMiddleware(http.HandlerFunc(checkHandler))).Methods("GET")

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Printf("Starting server on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Could not listen on %s: %v\n", port, err)
		}
	}()

	gracefulShutdown(srv)
}

func loadBloomFilter(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file %s: %v", filename, err)
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&filter); err != nil {
		return fmt.Errorf("error decoding bloom filter: %v", err)
	}

	return nil
}

func checkHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("s")
	if query == "" {
		http.Error(w, `{"error": "Missing 's' parameter"}`, http.StatusBadRequest)
		return
	}

	inSet := filter.TestString(query)
	response := Response{
		Query: query,
		InSet: inSet,
		Message: fmt.Sprintf("The string '%s' is %s in the set.",
			query, map[bool]string{true: "possibly", false: "definitely not"}[inSet]),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(2, 5)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func gracefulShutdown(srv *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
	logger.Println("shutting down")
	os.Exit(0)
}
