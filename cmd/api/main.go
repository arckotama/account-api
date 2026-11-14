package main

import (
	"log"
	"net/http" // baru

	"account-api/internal/config"
	"account-api/internal/database"
	"account-api/internal/handlers" 
	"account-api/internal/models"// baru
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DSN())
	if err != nil {
		log.Fatalf("gagal konek database: %v", err)
	}
	log.Println("terhubung ke database") // baru (menggantikan blok cek sementara)
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("gagal migrasi database: %v", err)
	}

	log.Println("database migration berhasil")

	userHandler := handlers.NewUserHandler(db) // baru

	// Routing: pola "METHOD /path" didukung ServeMux sejak Go 1.22.
	mux := http.NewServeMux() // baru
	mux.HandleFunc("GET /users", userHandler.List)
	mux.HandleFunc("POST /users", userHandler.Create)
	mux.HandleFunc("GET /users/{id}", userHandler.GetByID)
	mux.HandleFunc("PUT /users/{id}", userHandler.Update)
	mux.HandleFunc("DELETE /users/{id}", userHandler.Delete)

	addr := ":" + cfg.AppPort // baru
	log.Printf("server berjalan di port %s", cfg.AppPort)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server berhenti: %v", err)
	}
}