package handler

import (
	"auth-workflow/internal/models"
	"auth-workflow/internal/service"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

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