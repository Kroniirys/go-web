package controllers

import (
	"go-web/models"
	"go-web/utils"
	"html/template"
	"net/http"
)

func RegisterController(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, _ := template.ParseFiles("Views/register.html")
		tmpl.Execute(w, nil)
	} else if r.Method == "POST" {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		birthday := r.FormValue("birthday")
		gender := r.FormValue("gender")
		cccd := r.FormValue("cccd")

		// Mã hóa mật khẩu
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			http.Error(w, "Lỗi mã hóa mật khẩu", http.StatusInternalServerError)
			return
		}

		user := &models.User{
			Username: username,
			Email:    email,
			Password: hashedPassword,
			Birthday: birthday,
			Gender:   gender,
			CCCD:     cccd,
		}

		err = user.Create()
		if err != nil {
			http.Error(w, "Lỗi đăng ký", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func LoginController(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, _ := template.ParseFiles("Views/login.html")
		tmpl.Execute(w, nil)
	} else if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := models.GetUserByEmail(email)
		if err != nil {
			http.Error(w, "Email hoặc mật khẩu không đúng", http.StatusUnauthorized)
			return
		}

		// Kiểm tra mật khẩu đã mã hóa
		if !utils.CheckPasswordHash(password, user.Password) {
			http.Error(w, "Email hoặc mật khẩu không đúng", http.StatusUnauthorized)
			return
		}

		// Tạo session hoặc JWT token ở đây
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
