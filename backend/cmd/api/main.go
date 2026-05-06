package main

import (
	"fmt"
	"log"
	"net/http"

	"weather-outfit-backend/internal/db"
	"weather-outfit-backend/internal/handler"
)

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	conn, err := db.NewDB()
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	defer conn.Close()

	handler.SetDB(conn)

	fmt.Println("database connected")

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Go backend is running")
	})

	mux.HandleFunc("/signup", handler.SignupHandler)
	mux.HandleFunc("/login", handler.LoginHandler)
	mux.HandleFunc("/me", handler.MeHandler)

	fmt.Println("server is running on :8080")
	http.ListenAndServe(":8080", enableCors(mux))
}