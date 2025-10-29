package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"gophermart/internal/auth"
	"gophermart/internal/config"
	"gophermart/internal/orders"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type Server struct {
	cfg    *config.Config
	db     *gorm.DB
	router *mux.Router
}

func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	s := &Server{
		cfg:    cfg,
		db:     db,
		router: mux.NewRouter(),
	}
	s.routers()
	return s
}

func (s *Server) routers() {
	authTokenManager := auth.NewTokenManager(s.cfg.SecretToken, 24*time.Hour)

	worker := orders.NewWorker(s.db, s.cfg.AccrualSystemAddress, 10*time.Second)
	worker.Start()

	protected := s.router.PathPrefix("/api/user").Subrouter()
	protected.Use(auth.Middleware(authTokenManager))

	//Auth
	authRepo := auth.NewRepository(s.db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService, authTokenManager)

	s.router.HandleFunc("/api/user/register", authHandler.Register).Methods("POST")
	s.router.HandleFunc("/api/user/login", authHandler.Login).Methods("POST")

	//Orders
	ordersRepo := orders.NewRepository(s.db)
	ordersService := orders.NewService(ordersRepo)
	ordersHandler := orders.NewHandler(ordersService)

	protected.HandleFunc("/orders", ordersHandler.List).Methods("GET")
	protected.HandleFunc("/orders", ordersHandler.Create).Methods("POST")
}

func (s *Server) Start() {
	addr := s.cfg.RunAddress
	fmt.Println("Server running at", addr)
	if err := http.ListenAndServe(addr, s.router); err != nil {
		log.Fatal(err)
	}
}
