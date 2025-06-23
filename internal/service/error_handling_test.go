package service

import (
	"context"
	"testing"
	"user-service/internal/model"
)

// TestErrorTypeConsistency ensures that the service returns the expected error types
func TestErrorTypeConsistency(t *testing.T) {
	repo := &mockUserRepo{users: make(map[string]*model.User)}
	svc := NewUserService(repo)
	ctx := context.Background()

	// Test GetUser with non-existent email
	_, err := svc.GetUser(ctx, "nonexistent@example.com")
	if err != ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}

	// Test UpdateUser with non-existent email
	user := &model.User{Email: "nonexistent@example.com", Name: "Test", Age: 25}
	err = svc.UpdateUser(ctx, user)
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}

	// Test DeleteUser with non-existent email
	err = svc.DeleteUser(ctx, "nonexistent@example.com")
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
}
