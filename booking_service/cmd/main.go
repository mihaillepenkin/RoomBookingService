package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"booking_service/internal/config"
	"booking_service/internal/grpc"
	handler "booking_service/internal/infrastructure/http"
	"booking_service/internal/infrastructure/jwt"
	"booking_service/internal/infrastructure/postgres"
	"booking_service/internal/interfaces"
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
	schRepo := postgres.ScheduleRepository{DB: db}
	slotRepo := postgres.SlotRepository{DB: db}
	bookingRepo := postgres.BookingRepository{DB: db}
	roomServiceAddr := cfg.RoomAdress
    roomClient, err := grpc.NewRoomServiceClient(roomServiceAddr)
    if err != nil {
        log.Fatalf("failed to create room client: %v", err)
    }
    defer roomClient.Close() 
	jwtt := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.AccessExpiry)
    scheduleService := interfaces.NewScheduleService(&schRepo, roomClient, jwtt)
	slotService := interfaces.NewSlotService(&slotRepo, jwtt, &schRepo, roomClient)
	bookingService := interfaces.NewBookingService(&bookingRepo, jwtt, &schRepo)
	handl := handler.NewHandler(scheduleService, slotService, bookingService)
	m := mux.NewRouter()
	m.HandleFunc("/rooms/{roomId}/schedule/create", handl.CreateScheduleHandler).Methods("POST")
	m.HandleFunc("/rooms/{roomId}/slots/list", handl.GetAllSlotsHandler).Methods("GET")
	m.HandleFunc("/bookings/create", handl.CreateBookingHandler).Methods("POST")
	m.HandleFunc("/bookings/my", handl.GetMyBooksHandler).Methods("GET")
	m.HandleFunc("/bookings/{bookingId}/cancel", handl.CancelBookingHandler).Methods("POST")
	m.HandleFunc("/bookings/list", handl.GetAllBookingsHandler).Methods("GET")
	m.HandleFunc("/_info", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	srv := &http.Server{
		Addr:              ":" + strconv.Itoa(cfg.Port),
		Handler:           m,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("booking service start")

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