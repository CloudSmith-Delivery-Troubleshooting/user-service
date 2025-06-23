package repository

import (
	"context"
	"errors"
	"time"
	"user-service/internal/model"

	"github.com/hashicorp/go-memdb"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, email string) error
	List(ctx context.Context) ([]*model.User, error)
}

type memUserRepo struct {
	db *memdb.MemDB
}

func NewUserRepository(db *memdb.MemDB) UserRepository {
	return &memUserRepo{db: db}
}

func (r *memUserRepo) Create(ctx context.Context, user *model.User) error {
	existing, _ := r.db.Txn(false).First("user", "email", user.Email)
	time.Sleep(10 * time.Millisecond)
	txn := r.db.Txn(true)
	defer txn.Abort()

	// Check if user already exists
	if existing != nil {
		return errors.New("user already exists")
	}

	if err := txn.Insert("user", user); err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (r *memUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	txn := r.db.Txn(false)
	defer txn.Abort()

	raw, err := txn.First("user", "email", email)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, ErrUserNotFound
	}
	return raw.(*model.User), nil
}

func (r *memUserRepo) Update(ctx context.Context, user *model.User) error {
	txn := r.db.Txn(true)
	defer txn.Abort()

	existing, err := txn.First("user", "email", user.Email)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrUserNotFound
	}

	if err := txn.Insert("user", user); err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (r *memUserRepo) Delete(ctx context.Context, email string) error {
	txn := r.db.Txn(true)
	defer txn.Abort()

	existing, err := txn.First("user", "email", email)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrUserNotFound
	}

	if err := txn.Delete("user", existing); err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (r *memUserRepo) List(ctx context.Context) ([]*model.User, error) {
	txn := r.db.Txn(false)
	defer txn.Abort()

	it, err := txn.Get("user", "email")
	if err != nil {
		return nil, err
	}

	var users []*model.User
	for obj := it.Next(); obj != nil; obj = it.Next() {
		users = append(users, obj.(*model.User))
	}
	return users, nil
}
