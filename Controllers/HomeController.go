package controllers

import (
	"go-web/models"
	"html/template"
	"net/http"
)

func HomeController(w http.ResponseWriter, r *http.Request) {
	// Tạo dữ liệu mẫu
	user := models.User{
		ID:       1,
		Username: "Người dùng mẫu",
		Email:    "example@email.com",
	}

	// Parse template
	tmpl, err := template.ParseFiles("Views/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render template với dữ liệu
	err = tmpl.Execute(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AboutController(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("Views/about.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
