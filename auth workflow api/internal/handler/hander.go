package handler

import (
	"auth-workflow/internal/models"
	"auth-workflow/internal/service"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func NewHnadlerInstance(db *sql.DB)*HandlerUser  {
	return &HandlerUser{
		handler: service.NewUserService(db),
	}
}

type HandlerUser struct{
	handler *service.UserService
}

func (h *HandlerUser) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	user.ID = uuid.NewString()

	err := h.handler.ProcessUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *HandlerUser)LoginUser(w http.ResponseWriter,r *http.Request) {
	var logedUser models.LoginRequest

	token,err := h.handler.Authentication(&logedUser)
	//this return token when the user in the database is avlid and generate the new jwt token and return this token
	if token=="" {
		http.Error(w,"User not found",http.StatusBadGateway)
		return
	}

	if err!=nil {
		http.Error(w,"User not found",http.StatusInternalServerError)
		return
	}

	//at this stage token genearted after login
	// send token response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token":token,
	})
}

func (h *HandlerUser) Profile(w http.ResponseWriter, r *http.Request) {

	// 1️⃣ Authorization header lo
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "missing authorization header", http.StatusUnauthorized)
		return
	}

	// Expected: "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "invalid authorization format", http.StatusUnauthorized)
		return
	}

	tokenString := parts[1]

	//  JWT verify + parse
	claims, err := h.handler.VerifyJWT(tokenString)
	if err != nil {
		http.Error(w, "invalid or expired token", http.StatusUnauthorized)
		return
	}

	//  userID claims se nikalo
	userID := claims.UserID
	if userID == "" {
		http.Error(w, "invalid token claims", http.StatusUnauthorized)
		return
	}

	//  DB se user profile lao
	user, err := h.handler.GetUserByID(userID)
	if err != nil {
		http.Error(w, "failed to fetch profile", http.StatusInternalServerError)
		return
	}

	//  Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
