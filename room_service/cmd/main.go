package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"room_service/internal/config"
	"room_service/internal/grpc"
	handler "room_service/internal/infrastructure/http"
	"room_service/internal/infrastructure/jwt"
	"room_service/internal/infrastructure/postgres"
	"room_service/internal/interfaces"
	"strconv"
	"syscall"
	"time"

	roomv1 "room_service/proto/room/v1"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	grpcс "google.golang.org/grpc"
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
	repo := postgres.RoomRepository{DB: db}
	grpcServer := grpcс.NewServer()
	myGrpcServer := grpc.NewRoomServer(&repo)
	roomv1.RegisterRoomServiceServer(grpcServer, myGrpcServer)
	go func() {
		lis, err := net.Listen("tcp", ":" + cfg.RoomAdress)
		if err != nil {
			log.Fatalf("failed to listen gRPC: %v", err)
		}
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()
	jwtt := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.AccessExpiry)
	roomServ := interfaces.NewRoomUsecase(jwtt, &repo)
	handl := handler.NewHandler(roomServ)
	m := mux.NewRouter()
	m.HandleFunc("/rooms/create", handl.CreateHandler).Methods("POST")
	m.HandleFunc("/rooms/list", handl.GetListHandler).Methods("GET")
	m.HandleFunc("/rooms/{id}", handl.GetHandler).Methods("GET")

	srv := &http.Server{
		Addr:              ":" + strconv.Itoa(cfg.Port),
		Handler:           m,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}


	log.Println("room service start")

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