package models

import (
	"database/sql"
	"go-web-admin/config"
	"go-web-admin/utils"
	"time"
)

type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	Birthday  string
	Gender    string
	CCCD      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetAllUsers lấy danh sách tất cả người dùng
func GetAllUsers() ([]User, error) {
	query := `
		SELECT id, username, email, birthday, gender, cccd, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Birthday,
			&user.Gender,
			&user.CCCD,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// GetUserByID lấy thông tin người dùng theo ID
func GetUserByID(id int) (*User, error) {
	query := `
		SELECT id, username, email, birthday, gender, cccd, created_at, updated_at
		FROM users
		WHERE id = ?
	`
	row := config.DB.QueryRow(query, id)

	var user User
	var createdAt, updatedAt string
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Birthday,
		&user.Gender,
		&user.CCCD,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Parse thời gian
	user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	user.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	return &user, nil
}

// Create thêm người dùng mới
func (u *User) Create() error {
	// Mã hóa mật khẩu
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	// Thêm người dùng vào database
	query := `
		INSERT INTO users (username, email, password, birthday, gender, cccd, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	result, err := config.DB.Exec(query, u.Username, u.Email, hashedPassword, u.Birthday, u.Gender, u.CCCD)
	if err != nil {
		return err
	}

	// Lấy ID của người dùng vừa tạo
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = int(id)

	return nil
}

// Update cập nhật thông tin người dùng
func (u *User) Update() error {
	query := `
		UPDATE users 
		SET username = ?, email = ?, birthday = ?, gender = ?, cccd = ?, updated_at = NOW()
	`
	args := []interface{}{u.Username, u.Email, u.Birthday, u.Gender, u.CCCD}

	// Nếu có mật khẩu mới, cập nhật mật khẩu
	if u.Password != "" {
		hashedPassword, err := utils.HashPassword(u.Password)
		if err != nil {
			return err
		}
		query += ", password = ?"
		args = append(args, hashedPassword)
	}

	query += " WHERE id = ?"
	args = append(args, u.ID)

	_, err := config.DB.Exec(query, args...)
	return err
}

// DeleteUser xóa người dùng theo ID
func DeleteUser(id int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := config.DB.Exec(query, id)
	return err
}

func SearchUsers(query string) ([]User, error) {
	searchQuery := "%" + query + "%"
	sqlQuery := `
		SELECT id, username, email, birthday, gender, cccd, created_at, updated_at
		FROM users
		WHERE username LIKE ? OR email LIKE ? OR cccd LIKE ?
		ORDER BY created_at DESC
	`
	rows, err := config.DB.Query(sqlQuery, searchQuery, searchQuery, searchQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var createdAt, updatedAt string
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Birthday,
			&user.Gender,
			&user.CCCD,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse thời gian
		user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		user.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

		users = append(users, user)
	}

	return users, nil
}
 