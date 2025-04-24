package main

import (
	"bufio"
	"fmt"
	"go-web-admin/config"
	"go-web-admin/models"
	"go-web-admin/utils"
	"os"
	"strings"
)

func main() {
	// Khởi tạo kết nối database
	config.InitDB()
	defer config.DB.Close()

	reader := bufio.NewReader(os.Stdin)

	// Nhập thông tin admin
	fmt.Print("Nhập tên đăng nhập admin: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Nhập email admin: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Nhập mật khẩu admin: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	// Hash mật khẩu
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		fmt.Println("Lỗi hash mật khẩu:", err)
		return
	}

	// Tạo admin mới
	admin := &models.Admin{
		Username: username,
		Password: hashedPassword,
		Email:    email,
		Role:     "admin",
	}

	err = admin.Create()
	if err != nil {
		fmt.Println("Lỗi tạo admin:", err)
		return
	}

	fmt.Println("Tạo tài khoản admin thành công!")
	fmt.Println("Tên đăng nhập:", username)
	fmt.Println("Email:", email)
}
