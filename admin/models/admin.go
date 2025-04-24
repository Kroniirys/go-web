package models

import (
	"go-web-admin/config"
	"time"
)

type Admin struct {
	ID        int
	Username  string
	Password  string
	Email     string
	Role      string
	LastLogin *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AdminLog struct {
	ID        int
	AdminID   int
	Action    string
	Details   string
	IPAddress string
	CreatedAt time.Time
}

func (a *Admin) Create() error {
	query := `INSERT INTO admins (username, password, email, role) 
	          VALUES (?, ?, ?, ?)`
	_, err := config.DB.Exec(query, a.Username, a.Password, a.Email, a.Role)
	return err
}

func GetAdminByUsername(username string) (*Admin, error) {
	admin := &Admin{}
	query := `SELECT id, username, password, email, role, last_login, created_at, updated_at 
	          FROM admins WHERE username = ?`
	err := config.DB.QueryRow(query, username).Scan(
		&admin.ID, &admin.Username, &admin.Password, &admin.Email, &admin.Role,
		&admin.LastLogin, &admin.CreatedAt, &admin.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

func (a *Admin) UpdateLastLogin() error {
	query := `UPDATE admins SET last_login = ? WHERE id = ?`
	_, err := config.DB.Exec(query, time.Now(), a.ID)
	return err
}

func (l *AdminLog) Create() error {
	query := `INSERT INTO admin_logs (admin_id, action, details, ip_address) 
	          VALUES (?, ?, ?, ?)`
	_, err := config.DB.Exec(query, l.AdminID, l.Action, l.Details, l.IPAddress)
	return err
}

func GetAdminLogs(page, limit int) ([]AdminLog, error) {
	offset := (page - 1) * limit
	query := `SELECT id, admin_id, action, details, ip_address, created_at 
	          FROM admin_logs ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := config.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []AdminLog
	for rows.Next() {
		var log AdminLog
		err := rows.Scan(
			&log.ID, &log.AdminID, &log.Action, &log.Details,
			&log.IPAddress, &log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}
