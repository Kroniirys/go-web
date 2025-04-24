package models

import (
	"go-web/config"
)

type User struct {
	ID       int
	Username string
	Email    string
	Password string
	Birthday string
	Gender   string
	CCCD     string
}

func (u *User) Create() error {
	query := "INSERT INTO users (username, email, password, birthday, gender, cccd) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := config.DB.Exec(query, u.Username, u.Email, u.Password, u.Birthday, u.Gender, u.CCCD)
	return err
}

func GetUserByEmail(email string) (*User, error) {
	user := &User{}
	query := "SELECT id, username, email, password, birthday, gender, cccd FROM users WHERE email = ?"
	err := config.DB.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Birthday, &user.Gender, &user.CCCD)
	if err != nil {
		return nil, err
	}
	return user, nil
}
