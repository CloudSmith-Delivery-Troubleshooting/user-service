package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-service/internal/model"

	"github.com/gorilla/mux"
	"go.uber.org/zap/zaptest"
)

// TestStatusCodes ensures that handlers return the expected HTTP status codes
func TestStatusCodes(t *testing.T) {
	logger := zaptest.NewLogger(t)
	svc := &mockUserService{users: make(map[string]*model.User)}
	handler := NewUserHandler(svc, logger)

	t.Run("CreateUser returns 201 Created", func(t *testing.T) {
		user := model.User{Email: "new@example.com", Name: "New User", Age: 25}
		body, _ := json.Marshal(user)
		req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.CreateUser(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
		}
	})

	t.Run("GetUser returns 404 for non-existent user", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/nonexistent@example.com", nil)
		req = mux.SetURLVars(req, map[string]string{"email": "nonexistent@example.com"})
		w := httptest.NewRecorder()

		handler.GetUser(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("DeleteUser returns 204 No Content", func(t *testing.T) {
		// First create a user
		user := &model.User{Email: "delete@example.com", Name: "Delete Me", Age: 30}
		svc.users[user.Email] = user

		req := httptest.NewRequest("DELETE", "/users/delete@example.com", nil)
		req = mux.SetURLVars(req, map[string]string{"email": "delete@example.com"})
		w := httptest.NewRecorder()

		handler.DeleteUser(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
		}
	})
}
