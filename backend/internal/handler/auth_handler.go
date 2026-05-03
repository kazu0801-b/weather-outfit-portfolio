package handler

import (
        "database/sql"
        "encoding/json"
        "net/http"
        "strings"

        "weather-outfit-backend/internal/model"
        "weather-outfit-backend/internal/repository"
        "weather-outfit-backend/internal/service"
)

var db *sql.DB

func SetDB(conn *sql.DB) {
        db = conn
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
                http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
                return
        }

        if db == nil {
                http.Error(w, "database not initialized", http.StatusInternalServerError)
                return
        }

        var req model.SignupRequest
        err := json.NewDecoder(r.Body).Decode(&req)
        if err != nil {
                http.Error(w, "invalid request body", http.StatusBadRequest)
                return
        }

        if req.Username == "" || req.Email == "" || req.Password == "" {
                http.Error(w, "username, email, and password are required", http.StatusBadRequest)
                return
        }

        hashedPassword, err := service.HashPassword(req.Password)
        if err != nil {
                http.Error(w, "failed to hash password", http.StatusInternalServerError)
                return
        }

        user := model.User{
                Username:     req.Username,
                Email:        req.Email,
                PasswordHash: hashedPassword,
        }

        userRepo := repository.NewUserRepository(db)
        err = userRepo.CreateUser(user)
        if err != nil {
                if strings.Contains(err.Error(), "Duplicate entry") {
                        http.Error(w, "email already exists", http.StatusConflict)
                        return
                }
                http.Error(w, "failed to create user", http.StatusInternalServerError)
                return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
                "message": "user created",
        })
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
                http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
                return
        }

        if db == nil {
                http.Error(w, "database not initialized", http.StatusInternalServerError)
                return
        }

        var req model.LoginRequest
        err := json.NewDecoder(r.Body).Decode(&req)
        if err != nil {
                http.Error(w, "invalid request body", http.StatusBadRequest)
                return
        }

        if req.Email == "" || req.Password == "" {
                http.Error(w, "email and password are required", http.StatusBadRequest)
                return
        }

        userRepo := repository.NewUserRepository(db)
        user, err := userRepo.FindUserByEmail(req.Email)
        if err != nil {
                if err == sql.ErrNoRows {
                        http.Error(w, "invalid email or password", http.StatusUnauthorized)
                        return
                }
                http.Error(w, "failed to find user", http.StatusInternalServerError)
                return
        }

        err = service.CheckPasswordHash(req.Password, user.PasswordHash)
        if err != nil {
                http.Error(w, "invalid email or password", http.StatusUnauthorized)
                return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
                "message": "login successful",
        })
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
                "message": "me endpoint reached",
        })
}