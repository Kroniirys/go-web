package main

import (
	"go-web-admin/config"
	"go-web-admin/models"
	"math/rand"
	"time"
)

func main() {
	// Khởi tạo kết nối database
	err := config.InitDB()
	if err != nil {
		println("Lỗi kết nối database:", err.Error())
		return
	}
	defer config.DB.Close()

	// Khởi tạo random seed
	rand.Seed(time.Now().UnixNano())

	// Danh sách tên ngẫu nhiên
	firstNames := []string{"Nguyễn", "Trần", "Lê", "Phạm", "Hoàng", "Huỳnh", "Phan", "Vũ", "Võ", "Đặng"}
	lastNames := []string{"Văn", "Thị", "Hữu", "Đức", "Minh", "Hồng", "Thanh", "Quốc", "Bảo", "Anh"}
	middleNames := []string{"Hải", "Thành", "Minh", "Hồng", "Thanh", "Quốc", "Bảo", "Anh", "Duy", "Tuấn"}

	// Danh sách email domain
	emailDomains := []string{"gmail.com", "yahoo.com", "hotmail.com", "outlook.com"}

	// Tạo i người dùng
	for i := 0; i < 10; i++ {
		// Tạo tên ngẫu nhiên
		firstName := firstNames[rand.Intn(len(firstNames))]
		middleName := middleNames[rand.Intn(len(middleNames))]
		lastName := lastNames[rand.Intn(len(lastNames))]
		username := firstName + " " + middleName + " " + lastName

		// Tạo email ngẫu nhiên
		email := username + "@" + emailDomains[rand.Intn(len(emailDomains))]

		// Tạo CCCD ngẫu nhiên (12 số)
		cccd := ""
		for j := 0; j < 12; j++ {
			cccd += string(rune(rand.Intn(10) + '0'))
		}

		// Tạo ngày sinh ngẫu nhiên (từ 18-60 tuổi)
		now := time.Now()
		years := rand.Intn(42) + 18 // 18-60 tuổi
		days := rand.Intn(365)
		birthday := now.AddDate(-years, 0, -days)

		// Tạo giới tính ngẫu nhiên
		genders := []string{"Nam", "Nữ", "Khác"}
		gender := genders[rand.Intn(len(genders))]

		// Tạo mật khẩu ngẫu nhiên
		password := "123456" // Mật khẩu mặc định

		// Tạo người dùng
		user := &models.User{
			Username: username,
			Email:    email,
			Password: password,
			Birthday: birthday.Format("2006-01-02"),
			Gender:   gender,
			CCCD:     cccd,
		}

		// Thêm vào database
		err := user.Create()
		if err != nil {
			println("Lỗi tạo người dùng", username, ":", err.Error())
		} else {
			println("Đã tạo người dùng:", username)
		}
	}

	println("Đã tạo xong 50 người dùng ngẫu nhiên")
}
