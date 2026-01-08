package service

import (
	"auth-workflow/internal/models"
	"auth-workflow/internal/repo_user.go"
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	service *repo_user.SaveRepo
}

func (s *UserService) ProcessUser(user *models.User) error {

	var existingEmail string

	err := s.service.DB.QueryRow(
		"SELECT email FROM users WHERE email = ?",
		user.Email,
	).Scan(&existingEmail)

	// ✅ User does NOT exist → create new
	if err == sql.ErrNoRows {

		// hash password
		hash, err := bcrypt.GenerateFromPassword(
			[]byte(user.PasswordHash),
			bcrypt.DefaultCost,
		)
		if err != nil {
			return err
		}

		user.PasswordHash = string(hash)

		_, err = s.service.SaveToMySQL(user)
		return err
	}

	//  Real DB error
	if err != nil {
		return err
	}

	//  User already exists
	return errors.New("user already exists")
}

func (s *UserService) Authentication(loginReq *models.LoginRequest) (string, error) {
	// get password hash using email
	passHash, err := s.service.GetPasswordHash(loginReq.Email)
	if err != nil {
		return "", err
	}

	// compare password
	err = bcrypt.CompareHashAndPassword(
		[]byte(passHash),
		[]byte(loginReq.Password),
	)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// generate JWT
	token, err := s.GenerateJWT(loginReq)
	if err != nil {
		return "", err
	}

	//befor returning accessed token make the refresh token
	//if refresh token not found then make new one
	//else found then simple return without doing anything

	

	return token, nil
}


func (s *UserService)GenerateJWT(loginreq *models.LoginRequest)(string,error)  {
	var jwtSecret = []byte("SUPER_SECRET_KEY") 
	userid,err := s.service.GetuserID(loginreq.Email)
	if err!=nil {
		return "",nil
	}

	//claim->token->signedtoken
	claim := jwt.MapClaims{
		ID : userid,
		"email":   loginreq.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256,claim)

	signedToken, err := token.SignedString(jwtSecret)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}
