package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-service/internal/model"
	"user-service/internal/service"

	"github.com/gorilla/mux"
	"go.uber.org/zap/zaptest"
)

// mockUserService implements service.UserService for testing
type mockUserService struct {
	users map[string]*model.User
}

func (m *mockUserService) CreateUser(_ context.Context, user *model.User) error {
	if _, exists := m.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	m.users[user.Email] = user
	return nil
}

func (m *mockUserService) GetUser(_ context.Context, email string) (*model.User, error) {
	user, ok := m.users[email]
	if !ok {
		return nil, service.ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserService) UpdateUser(_ context.Context, user *model.User) error {
	if _, ok := m.users[user.Email]; !ok {
		return service.ErrUserNotFound
	}
	m.users[user.Email] = user
	return nil
}

func (m *mockUserService) DeleteUser(_ context.Context, email string) error {
	if _, ok := m.users[email]; !ok {
		return service.ErrUserNotFound
	}
	delete(m.users, email)
	return nil
}

func (m *mockUserService) ListUsers(_ context.Context) ([]*model.User, error) {
	var list []*model.User
	for _, u := range m.users {
		list = append(list, u)
	}
	return list, nil
}

func setupHandler() (*UserHandler, *mockUserService) {
	svc := &mockUserService{users: make(map[string]*model.User)}
	logger := zaptest.NewLogger(nil)
	return NewUserHandler(svc, logger), svc
}

func TestCreateUserHandler(t *testing.T) {
	handler, _ := setupHandler()

	user := model.User{Email: "test@example.com", Name: "Test", Age: 25}
	body, _ := json.Marshal(user)

	req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateUser(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %d", w.Code)
	}
}

func TestGetUserHandler(t *testing.T) {
	handler, svc := setupHandler()

	user := &model.User{Email: "test@example.com", Name: "Test", Age: 25}
	svc.users[user.Email] = user

	req := httptest.NewRequest("GET", "/users/test@example.com", nil)
	w := httptest.NewRecorder()

	// We need to set the route variables for mux
	r := mux.NewRouter()
	r.HandleFunc("/users/{email}", handler.GetUser)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", w.Code)
	}

	var got model.User
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, got.Email)
	}
}

func TestUpdateUserHandler(t *testing.T) {
	handler, svc := setupHandler()

	user := &model.User{Email: "test@example.com", Name: "Test", Age: 25}
	svc.users[user.Email] = user

	updated := model.User{Name: "Updated", Age: 30}
	body, _ := json.Marshal(updated)

	req := httptest.NewRequest("PUT", "/users/test@example.com", bytes.NewReader(body))
	w := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/users/{email}", handler.UpdateUser)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", w.Code)
	}

	got, _ := svc.GetUser(nil, user.Email)
	if got.Name != "Updated" || got.Age != 30 {
		t.Errorf("user not updated correctly, got %+v", got)
	}
}

func TestDeleteUserHandler(t *testing.T) {
	handler, svc := setupHandler()

	user := &model.User{Email: "test@example.com", Name: "Test", Age: 25}
	svc.users[user.Email] = user

	req := httptest.NewRequest("DELETE", "/users/test@example.com", nil)
	w := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/users/{email}", handler.DeleteUser)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204 No Content, got %d", w.Code)
	}

	if _, err := svc.GetUser(context.TODO(), user.Email); err == nil {
		t.Errorf("user was not deleted")
	}
}

func TestListUsersHandler(t *testing.T) {
	handler, svc := setupHandler()

	svc.users["a@example.com"] = &model.User{Email: "a@example.com", Name: "A", Age: 20}
	svc.users["b@example.com"] = &model.User{Email: "b@example.com", Name: "B", Age: 30}

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	handler.ListUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", w.Code)
	}

	var users []model.User
	if err := json.NewDecoder(w.Body).Decode(&users); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}
