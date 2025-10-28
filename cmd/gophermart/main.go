package main

import (
	"fmt"
	"gophermart/internal/app"
	"gophermart/internal/config"
	"gophermart/internal/db"
)

func main() {
	cfg := config.Load()
	fmt.Printf("Server on %s\nDB: %s\nAccrual: %s\n", cfg.RunAddress, cfg.DatabaseURI, cfg.AccrualSystemAddress)

	conn := db.Init(cfg.DatabaseURI)
	defer db.Close(conn)
	fmt.Println("Connected to database:", conn != nil)

	server := app.NewServer(cfg, conn)
	server.Start()
}
