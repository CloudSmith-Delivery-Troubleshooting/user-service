package service

import (
	"context"
	"user-service/internal/model"
	"user-service/internal/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, email string) error
	ListUsers(ctx context.Context) ([]*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, user *model.User) error {
	return s.repo.Create(ctx, user)
}

func (s *userService) GetUser(ctx context.Context, email string) (*model.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *userService) UpdateUser(ctx context.Context, user *model.User) error {
	return s.repo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, email string) error {
	return s.repo.Delete(ctx, email)
}

func (s *userService) ListUsers(ctx context.Context) ([]*model.User, error) {
	return s.repo.List(ctx)
}
