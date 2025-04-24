package main

import (
	"go-web/config"
	"go-web/controllers"
	"log"
	"net/http"
)

func main() {
	// Khởi tạo kết nối database
	config.InitDB()

	// Phục vụ các file tĩnh
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Cấu hình route
	http.HandleFunc("/", controllers.HomeController)
	http.HandleFunc("/about", controllers.AboutController)
	http.HandleFunc("/register", controllers.RegisterController)
	http.HandleFunc("/login", controllers.LoginController)

	// Khởi động server
	log.Println("Server đang chạy tại http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
