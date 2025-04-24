package main

import (
	"go-web-admin/config"
	"go-web-admin/controllers"
	"log"
	"net/http"
)

func main() {
	// Khởi tạo database
	config.InitDB()
	defer config.DB.Close()

	// Cấu hình các route
	http.HandleFunc("/admin/login", controllers.LoginHandler)
	http.HandleFunc("/admin/logout", controllers.LogoutHandler)

	// Routes cho dashboard
	http.Handle("/admin/dashboard", controllers.AuthMiddleware(http.HandlerFunc(controllers.DashboardHandler)))

	// Routes cho quản lý users
	http.Handle("/admin/users", controllers.AuthMiddleware(http.HandlerFunc(controllers.UsersHandler)))
	http.Handle("/admin/users/edit/", controllers.AuthMiddleware(http.HandlerFunc(controllers.EditUserHandler)))
	http.Handle("/admin/users/search", controllers.AuthMiddleware(http.HandlerFunc(controllers.SearchUsersHandler)))

	// Routes cho quản lý sách
	http.Handle("/admin/books", controllers.AuthMiddleware(http.HandlerFunc(controllers.BooksHandler)))
	http.Handle("/admin/books/edit/", controllers.AuthMiddleware(http.HandlerFunc(controllers.EditBookHandler)))
	http.Handle("/admin/books/search", controllers.AuthMiddleware(http.HandlerFunc(controllers.SearchBooksHandler)))

	// Routes cho quản lý đơn hàng
	http.Handle("/admin/orders", controllers.AuthMiddleware(http.HandlerFunc(controllers.OrdersHandler)))
	http.Handle("/admin/orders/detail/", controllers.AuthMiddleware(http.HandlerFunc(controllers.OrderDetailHandler)))

	// Routes cho quản lý logs
	http.Handle("/admin/logs", controllers.AuthMiddleware(http.HandlerFunc(controllers.LogsHandler)))

	http.HandleFunc("/admin/", controllers.AdminHandler)

	// Cấu hình static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Khởi động server
	log.Println("Server đang chạy tại http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
