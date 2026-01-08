package repo_user

import (
	"auth-workflow/internal/models"
	"database/sql"
	"errors"
)

type SaveRepo struct {
	DB *sql.DB
}

func (r *SaveRepo) SaveToMySQL(user *models.User) (*models.User, error) {

	_, err := r.DB.Exec(
		`INSERT INTO users (id, email, password_hash)
		 VALUES (?, ?, ?)`,
		user.ID,
		user.Email,
		user.PasswordHash,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *SaveRepo) GetPasswordHash(email string) (string, error) {
	var foundPass string

	query := `SELECT password_hash FROM users WHERE email = ?`

	err := r.DB.QueryRow(query, email).Scan(&foundPass)
	if err != nil {
		if err == sql.ErrNoRows {
			// email not found
			return "", errors.New("user not found")
		}
		// other DB error
		return "", err
	}


	return foundPass, nil
}

func (r *SaveRepo) GetuserID(email string) (string, error) {
	var foundID string

	query := `SELECT ID FROM users WHERE email = ?`

	err := r.DB.QueryRow(query, email).Scan(&foundID)
	if err != nil {
		if err == sql.ErrNoRows {
			// email not found
			return "", errors.New("ID not found")
		}
		// other DB error
		return "", err
	}


	return foundID, nil
}

//refresh token
func (r *SaveRepo) SaveRefToken(user *models.User) (*models.User, error) {
	

}







