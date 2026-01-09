package service

import (
	"auth-workflow/internal/models"
	"auth-workflow/internal/repo_user.go"
	"auth-workflow/internal/service"
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func NewUserService(db *sql.DB)*UserService  {
	return &UserService{
		service: repo_user.NewrepoUserInstance(db) ,
	}
}

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

	_, err = s.generateRefToken(loginReq)
	if err != nil {
		return "", err
	}

	return token, nil
}


var jwtSecret = []byte("SUPER_SECRET_KEY") 
func (s *UserService)GenerateJWT(loginreq *models.LoginRequest)(string,error)  {
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

func (s *UserService) generateRefToken(
	loqinDetail *models.LoginRequest,
) (*models.RefreshToken, error) {

	id, err := s.service.GetuserID(loqinDetail.Email)
	if err != nil {
		return nil, err
	}

	exists, err := s.service.IsRefTokenExists(id)
	if err != nil {
		return nil, err
	}

	if exists {
		// FIX: nil,nil mat bhejo — clear signal do
		return &models.RefreshToken{}, nil
	}

	tkn := &models.RefreshToken{
		ID:      uuid.NewString(),
		UserID:  id,
		Token:   uuid.NewString(),
		Revoked: false,
	}

	reftoken, err := s.service.SaveRefToken(tkn)
	if err != nil {
		return nil, err
	}

	return reftoken, nil
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *UserService) VerifyJWT(tokenString string) (*Claims, error) {

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}


