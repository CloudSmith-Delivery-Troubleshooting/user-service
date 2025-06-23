package handler

import (
	"encoding/json"
	"net/http"
	"user-service/internal/model"
	"user-service/internal/service"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService service.UserService
	logger      *zap.Logger
}

func NewUserHandler(userService service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{userService: userService, logger: logger}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.userService.CreateUser(r.Context(), &user); err != nil {
		h.logger.Error("Failed to create user", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]
	user, err := h.userService.GetUser(r.Context(), email)
	if err != nil {
		h.logger.Error("User not found", zap.Error(err))
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	user.Email = email

	if err := h.userService.UpdateUser(r.Context(), &user); err != nil {
		h.logger.Error("Failed to update user", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]

	if err := h.userService.DeleteUser(r.Context(), email); err != nil {
		h.logger.Error("Failed to delete user", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.ListUsers(r.Context())
	if err != nil {
		h.logger.Error("Failed to list users", zap.Error(err))
		http.Error(w, "Failed to list users", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(users); err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
