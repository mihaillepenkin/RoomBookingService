package main

import (
	"auth_service/internal/config"
	handler "auth_service/internal/infrastructure/http"
	"auth_service/internal/infrastructure/jwt"
	"auth_service/internal/infrastructure/postgres"
	"auth_service/internal/interfaces"
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func initDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DB.DSN())
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func main() {
	cfg, err := config.LoadConfig()
	if (err != nil) {
		log.Fatalf("Failed to load config: %v", err)
	}
	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	repo := postgres.UserRepoPG{DB: db}
	jwtt := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.AccessExpiry)
	authServ := interfaces.NewAuthUsecase(&repo, jwtt)
	handl := handler.NewAuthHandler(authServ)
	m := mux.NewRouter()
	m.HandleFunc("/register", handl.RegistrHandler).Methods("POST")
	m.HandleFunc("/login", handl.LoginHandler).Methods("POST")
	m.HandleFunc("/dummyLogin", handl.DummyLoginHandler).Methods("POST")

	srv := &http.Server{
		Addr:              ":" + strconv.Itoa(cfg.Port),
		Handler:           m,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
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