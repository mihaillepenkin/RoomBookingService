package http

import (
	"booking_service/internal/interfaces"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type BookingHandler struct {
	sch interfaces.ScheduleService
}

func NewHandler(sch interfaces.ScheduleService) *BookingHandler {
	return &BookingHandler{sch: sch}
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