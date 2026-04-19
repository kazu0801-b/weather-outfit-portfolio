package main

import (
	"fmt"
	"log"
	"net/http"

	"weather-outfit-backend/internal/db"
	"weather-outfit-backend/internal/handler"
)

func main() {
	conn, err := db.NewDB()
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	defer conn.Close()

	handler.SetDB(conn)

	fmt.Println("database connected")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Go backend is running")
	})

	http.HandleFunc("/signup", handler.SignupHandler)
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/me", handler.MeHandler)

	fmt.Println("server is running on :8080")
	http.ListenAndServe(":8080", nil)
}