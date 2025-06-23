package service

import (
	"context"
	"errors"
	"testing"
	"user-service/internal/model"
	"user-service/internal/repository"
)

// mockUserRepo implements UserRepository for testing
type mockUserRepo struct {
	users map[string]*model.User
}

func (m *mockUserRepo) Create(ctx context.Context, user *model.User) error {
	if _, exists := m.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user, ok := m.users[email]
	if !ok {
		return nil, repository.ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserRepo) Update(ctx context.Context, user *model.User) error {
	if _, ok := m.users[user.Email]; !ok {
		return repository.ErrUserNotFound
	}
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepo) Delete(ctx context.Context, email string) error {
	if _, ok := m.users[email]; !ok {
		return repository.ErrUserNotFound
	}
	delete(m.users, email)
	return nil
}

func (m *mockUserRepo) List(ctx context.Context) ([]*model.User, error) {
	var list []*model.User
	for _, u := range m.users {
		list = append(list, u)
	}
	return list, nil
}

func TestUserService(t *testing.T) {
	repo := &mockUserRepo{users: make(map[string]*model.User)}
	svc := NewUserService(repo)

	ctx := context.Background()
	user := &model.User{Email: "test@example.com", Name: "Test User", Age: 25}

	// Test CreateUser
	if err := svc.CreateUser(ctx, user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// Test CreateUser duplicate
	if err := svc.CreateUser(ctx, user); err == nil {
		t.Fatalf("Expected error for duplicate user, got nil")
	}

	// Test GetUser
	got, err := svc.GetUser(ctx, user.Email)
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if got.Name != user.Name {
		t.Errorf("GetUser returned wrong name: got %s, want %s", got.Name, user.Name)
	}

	// Test UpdateUser
	user.Name = "Updated Name"
	if err := svc.UpdateUser(ctx, user); err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}
	got, _ = svc.GetUser(ctx, user.Email)
	if got.Name != "Updated Name" {
		t.Errorf("UpdateUser did not update name: got %s", got.Name)
	}

	// Test DeleteUser
	if err := svc.DeleteUser(ctx, user.Email); err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}
	if _, err := svc.GetUser(ctx, user.Email); err == nil {
		t.Errorf("Expected error for deleted user, got nil")
	}

	// Test ListUsers
	_ = svc.CreateUser(ctx, &model.User{Email: "a@example.com", Name: "A", Age: 20})
	_ = svc.CreateUser(ctx, &model.User{Email: "b@example.com", Name: "B", Age: 30})
	users, err := svc.ListUsers(ctx)
	if err != nil {
		t.Fatalf("ListUsers failed: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("ListUsers returned wrong count: got %d, want 2", len(users))
	}
}
