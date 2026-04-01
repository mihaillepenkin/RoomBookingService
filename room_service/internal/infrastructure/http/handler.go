package http

import (
	"encoding/json"
	"net/http"
	"room_service/internal/interfaces"
	"strings"

	"github.com/gorilla/mux"
)

type RoomHandler struct {
	serv interfaces.RoomService
}

func NewHandler(serv interfaces.RoomService) *RoomHandler {
	return &RoomHandler{serv: serv}
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

func (h *RoomHandler) GetListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		writeError(w, 401, "you are not authtorized")
		return
	}
	_ = strings.TrimSpace(parts[1])
	output := (h.serv).ListRooms()
	writeJSON(w, output.Status, output.Data)
}

func (h *RoomHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		writeError(w, 401, "you are not authtorized")
		return
	}
	_ = strings.TrimSpace(parts[1])
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		writeError(w, 401, "bad request")
		return
	}
	var input interfaces.GetDTO
	input.Name = name
	output := (h.serv).GetRoom(input)
	writeJSON(w, output.Status, output.Data)
}

func (h *RoomHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		writeError(w, 401, "you are not authtorized")
		return
	}
	token := strings.TrimSpace(parts[1])
	var input interfaces.CreateDTO
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&input); err != nil {
		writeError(w, 400, "bad request")
		return
	}
	output := (h.serv).CreateRoom(input, token)
	writeJSON(w, output.Status, output.Data)
}