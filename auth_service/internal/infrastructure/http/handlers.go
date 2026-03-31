package http

import (
	"auth_service/internal/interfaces"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	uc interfaces.AuthService
}

func NewAuthHandler(uc interfaces.AuthService) *AuthHandler {
	return &AuthHandler{uc: uc}
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

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
		var input interfaces.LoginDTO
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&input); err != nil {
			writeError(w, 401, "bad request")
			return
		}

		output := (h.uc).Login(input)
		writeJSON(w, output.Status, output.Data)
}

func (h *AuthHandler) RegistrHandler(w http.ResponseWriter, r *http.Request) {
		var input interfaces.RegistrDTO
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&input); err != nil {
			writeError(w, 401, "bad request")
			return
		}

		output := (h.uc).Registr(input)
		writeJSON(w, output.Status, output.Data)
}

func (h *AuthHandler) DummyLoginHandler(w http.ResponseWriter, r *http.Request) {
		var input interfaces.DummyLoginDTO
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&input); err != nil {
			writeError(w, 401, "bad request")
			return
		}

		output := (h.uc).DummyLogin(input)
		writeJSON(w, output.Status, output.Data)
}