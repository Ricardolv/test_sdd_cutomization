package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/sdd-cod3r/test-app/apps/backend/config"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/database"
	customerHandlers "github.com/sdd-cod3r/test-app/apps/backend/internal/modules/customer/handlers"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/customer/repositories"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/customer/services"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/handlers"
	appmiddleware "github.com/sdd-cod3r/test-app/apps/backend/internal/modules/middleware"
)

func main() {
	cfg := config.Load()

	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}
	defer db.Close()

	customerRepo := repositories.NewPostgresCustomerRepository(db)
	customerService := services.NewCustomerService(customerRepo)
	customerHandler := customerHandlers.NewCustomerHandler(customerService)

	router := chi.NewRouter()

	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.Timeout(60 * time.Second))
	router.Use(appmiddleware.ErrorHandler)

	router.Get("/health", handlers.HealthCheck)
	router.Mount("/customers", customerHandler.Routes())

	fmt.Printf("Server starting on :%s\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
