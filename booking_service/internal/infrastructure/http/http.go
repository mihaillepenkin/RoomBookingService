package http

import (
	"booking_service/internal/interfaces"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type BookingHandler struct {
	sch     interfaces.ScheduleService
	slot    interfaces.SlotService
	booking interfaces.BookingService
}

func NewHandler(sch interfaces.ScheduleService, slot interfaces.SlotService, booking interfaces.BookingService) *BookingHandler {
	return &BookingHandler{sch: sch, slot: slot, booking: booking}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *BookingHandler) CreateScheduleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeError(w, 401, "you are not authtorized")
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		writeError(w, 401, "you are not authtorized")
		return
	}
	token := strings.TrimSpace(parts[1])
	vars := mux.Vars(r)
	id := vars["roomId"]
	if id == "" {
		writeError(w, 400, "bad request, room id is required")
		return
	}
	var input interfaces.CreateScheduleDTO
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&input); err != nil {
		writeError(w, 400, "bad request")
		return
	}
	input.RoomID = id
	output := (h.sch).CreateSchedule(r.Context(), input, token)
	writeJSON(w, output.Status, output.Data)
}

func (h *BookingHandler) GetAllSlotsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeError(w, 401, "you are not authtorized")
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		writeError(w, 401, "you are not authtorized")
		return
	}
	token := strings.TrimSpace(parts[1])
	vars := mux.Vars(r)
	id := vars["roomId"]
	if id == "" {
		writeError(w, 400, "bad request, room id is required")
		return
	}
	date := r.URL.Query().Get("date")
	if date == "" {
		writeError(w, 400, "bad request, date is required")
		return
	}
	var input interfaces.GetAllSlotsDTO
	input.RoomID = id
	input.Date = date
	output := (h.slot).GetAllSlots(r.Context(), input, token)
	writeJSON(w, output.Status, output.Data)
}

func (h *BookingHandler) CreateBookingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeError(w, 401, "you are not authtorized")
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		writeError(w, 401, "you are not authtorized")
		return
	}
	token := strings.TrimSpace(parts[1])
	var input interfaces.CreateBookingDTO
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&input); err != nil {
		writeError(w, 400, "bad request")
		return
	}
	output := (h.booking).CreateBook(r.Context(), input, token)
	writeJSON(w, output.Status, output.Data)
}

func (h *BookingHandler) GetAllBookingsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeError(w, 401, "you are not authorized")
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		writeError(w, 401, "you are not authorized")
		return
	}
	token := strings.TrimSpace(parts[1])
	query := r.URL.Query()
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(query.Get("pageSize"))
	if err != nil {
		pageSize = 20
	}
	input := interfaces.GetAllBooksDTO{Page: page, PageSize: pageSize}
	output := (h.booking).GetAllBookings(r.Context(), input, token)
	writeJSON(w, output.Status, output.Data)
}

func (h *BookingHandler) GetMyBooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeError(w, 401, "you are not authorized")
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		writeError(w, 401, "you are not authorized")
		return
	}
	token := strings.TrimSpace(parts[1])
	output := (h.booking).GetMyBooks(r.Context(), token)
	writeJSON(w, output.Status, output.Data)
}

func (h *BookingHandler) CancelBookingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeError(w, 401, "you are not authorized")
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		writeError(w, 401, "you are not authorized")
		return
	}
	token := strings.TrimSpace(parts[1])
	vars := mux.Vars(r)
	bookingID := vars["bookingId"]
	if bookingID == "" {
		writeError(w, 400, "booking ID is required")
		return
	}
	input := interfaces.CancelBookingDTO{BookingID: bookingID}
	output := (h.booking).CancelBook(r.Context(), input, token)
	writeJSON(w, output.Status, output.Data)
}
