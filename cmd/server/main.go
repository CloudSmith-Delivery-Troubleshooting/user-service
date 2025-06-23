package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-service/internal/handler"
	"user-service/internal/repository"
	"user-service/internal/service"
	"user-service/pkg/logger"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-memdb"
	"go.uber.org/zap"
)

func main() {
	// Setup logger
	zapLogger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("cannot initialize logger: %v", err)
	}
	defer zapLogger.Sync()

	// Setup in-memory DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"user": {
				Name: "user",
				Indexes: map[string]*memdb.IndexSchema{
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
		zapLogger.Fatal("failed to create memdb", zap.Error(err))
	}

	// Initialize repository, service, handler
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, zapLogger)

	// Setup router and routes
	r := mux.NewRouter()
	r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/users", userHandler.ListUsers).Methods("GET")
	r.HandleFunc("/users/{email}", userHandler.GetUser).Methods("GET")
	r.HandleFunc("/users/{email}", userHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{email}", userHandler.DeleteUser).Methods("DELETE")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		zapLogger.Info("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("Server failed", zap.Error(err))
		}
	}()

	// Graceful shutdown on SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zapLogger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Fatal("Server forced to shutdown", zap.Error(err))
	}
	zapLogger.Info("Server exited properly")
}
