package controllers

import (
	"go-web-admin/models"
	"go-web-admin/utils"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your-secret-key") // Thay đổi key này trong môi trường production

type Claims struct {
	Username string
	Role     string
	jwt.StandardClaims
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, _ := template.ParseFiles("views/login.html")
		tmpl.Execute(w, nil)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	log.Printf("Attempting login for username: %s", username)

	admin, err := models.GetAdminByUsername(username)
	if err != nil {
		log.Printf("Error getting admin: %v", err)
		http.Error(w, "Tên đăng nhập hoặc mật khẩu không đúng", http.StatusUnauthorized)
		return
	}

	log.Printf("Found admin: %+v", admin)

	if !utils.CheckPasswordHash(password, admin.Password) {
		log.Printf("Password mismatch for user: %s", username)
		http.Error(w, "Tên đăng nhập hoặc mật khẩu không đúng", http.StatusUnauthorized)
		return
	}

	// Tạo JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: admin.Username,
		Role:     admin.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Error creating token: %v", err)
		http.Error(w, "Lỗi tạo token", http.StatusInternalServerError)
		return
	}

	// Lưu token vào cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	// Cập nhật thời gian đăng nhập cuối
	admin.UpdateLastLogin()

	// Ghi log đăng nhập
	logEntry := &models.AdminLog{
		AdminID:   admin.ID,
		Action:    "login",
		Details:   "Đăng nhập thành công",
		IPAddress: r.RemoteAddr,
	}
	logEntry.Create()

	log.Printf("Login successful for user: %s", username)
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
				return
			}
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := c.Value
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	// Xử lý các request đến /admin/
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

// LogoutHandler xử lý đăng xuất
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Xóa cookie session
	cookie := &http.Cookie{
		Name:   "admin_session",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Xóa cookie
	}
	http.SetCookie(w, cookie)

	// Chuyển hướng về trang đăng nhập
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}
