package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func reverseProxy(target string) *httputil.ReverseProxy {
    targetURL, err := url.Parse(target)
    if err != nil {
        log.Fatalf("Invalid target URL: %v", err)
    }
    return httputil.NewSingleHostReverseProxy(targetURL)
}

func main() {
    authService := getEnv("AUTH_SERVICE_ADDR", "http://auth-service:8081")
    roomService := getEnv("ROOM_SERVICE_ADDR", "http://room-service:8082")
    bookingService := getEnv("BOOKING_SERVICE_ADDR", "http://booking-service:8083")
    r := mux.NewRouter()
    r.Path("/register").Handler(reverseProxy(authService))
    r.Path("/login").Handler(reverseProxy(authService))
    r.Path("/rooms/{roomId}/slots/list").Handler(reverseProxy(bookingService))
    r.Path("/rooms/{roomId}/schedule/create").Handler(reverseProxy(bookingService))
    r.Path("/rooms/{roomId}/slots").Handler(reverseProxy(bookingService))
    r.PathPrefix("/bookings").Handler(reverseProxy(bookingService))
    r.Path("/dummyLogin").Handler(reverseProxy(authService))
    r.Path("/rooms/list").Handler(reverseProxy(roomService))
    r.Path("/rooms/{id}").Handler(reverseProxy(roomService))
    r.Path("/rooms/create").Handler(reverseProxy(roomService))
    r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })
    r.HandleFunc("/_info", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")
    srv := &http.Server{
        Addr:              ":8080",
        Handler:           r,
        ReadHeaderTimeout: 5 * time.Second,
        ReadTimeout:       15 * time.Second,
        WriteTimeout:      15 * time.Second,
        IdleTimeout:       60 * time.Second,
    }

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Gateway error: %v", err)
        }
    }()

    quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped")
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}