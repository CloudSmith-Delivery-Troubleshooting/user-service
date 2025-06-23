// In internal/repository/user_repo_test.go
package repository

import (
	"context"
	"sync"
	"testing"
	"user-service/internal/model"

	"github.com/hashicorp/go-memdb"
)

func TestConcurrentUserCreation(t *testing.T) {
	// Setup in-memory DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"user": {
				Name: "user",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Email"},
					},
					"email": {
						Name:    "email",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Email"},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		t.Fatalf("failed to create memdb: %v", err)
	}

	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create a user that will be attempted to be created concurrently
	user := &model.User{
		Email: "concurrent@example.com",
		Name:  "Concurrent User",
		Age:   30,
	}

	// Number of concurrent goroutines
	concurrency := 10
	var wg sync.WaitGroup
	wg.Add(concurrency)

	// Channel to collect results
	results := make(chan error, concurrency)

	// Launch concurrent creation attempts
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			// Each goroutine tries to create the same user
			err := repo.Create(ctx, user)
			results <- err
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(results)

	// Count successful and failed creations
	successCount := 0
	for err := range results {
		if err == nil {
			successCount++
		}
	}

	// Only one creation should succeed
	if successCount != 1 {
		t.Errorf("Expected exactly 1 successful creation, got %d", successCount)
	}

	// Verify the user was actually created
	storedUser, err := repo.GetByEmail(ctx, user.Email)
	if err != nil {
		t.Errorf("Failed to get created user: %v", err)
	}
	if storedUser == nil {
		t.Error("User was not created")
	}
}
