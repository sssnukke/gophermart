package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"gophermart/internal/auth"
	"gophermart/internal/config"
	"gorm.io/gorm"
	"log"
	"net/http"
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
	authRepo := auth.NewRepository(s.db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	s.router.HandleFunc("/api/user/register", authHandler.Register).Methods("POST")
	s.router.HandleFunc("/api/user/login", authHandler.Login).Methods("POST")
}

func (s *Server) Start() {
	addr := s.cfg.RunAddress
	fmt.Println("Server running at", addr)
	if err := http.ListenAndServe(addr, s.router); err != nil {
		log.Fatal(err)
	}
}
