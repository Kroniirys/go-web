package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword mã hóa mật khẩu
func HashPassword(password string) (string, error) {
	// GenerateFromPassword tự động thêm salt và mã hóa mật khẩu
	// 14 là cost factor (số vòng lặp), càng cao càng an toàn nhưng cũng tốn nhiều tài nguyên hơn
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash kiểm tra mật khẩu có khớp với hash không
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
