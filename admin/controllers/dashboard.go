package controllers

import (
	"go-web-admin/config"
	"go-web-admin/models"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type GenderStats struct {
	Male   int
	Female int
	Other  int
}

type AgeStats struct {
	Under18   int
	Age18to25 int
	Age26to35 int
	Age36to45 int
	Over45    int
}

type DashboardData struct {
	TotalUsers    int
	NewUsersToday int
	GenderStats   GenderStats
	AgeStats      AgeStats
	GrowthLabels  []string
	GrowthData    []int
	RecentLogs    []models.AdminLog
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Lấy tổng số người dùng
	var totalUsers int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalUsers)
	if err != nil {
		http.Error(w, "Lỗi lấy dữ liệu người dùng", http.StatusInternalServerError)
		return
	}

	// Lấy số người dùng mới hôm nay
	var newUsersToday int
	today := time.Now().Format("2006-01-02")
	err = config.DB.QueryRow("SELECT COUNT(*) FROM users WHERE DATE(created_at) = ?", today).Scan(&newUsersToday)
	if err != nil {
		http.Error(w, "Lỗi lấy dữ liệu người dùng mới", http.StatusInternalServerError)
		return
	}

	// Lấy thống kê giới tính
	var genderStats GenderStats
	rows, err := config.DB.Query("SELECT gender, COUNT(*) as count FROM users GROUP BY gender")
	if err != nil {
		http.Error(w, "Lỗi lấy thống kê giới tính", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var gender string
		var count int
		err = rows.Scan(&gender, &count)
		if err != nil {
			continue
		}
		switch gender {
		case "Nam":
			genderStats.Male = count
		case "Nữ":
			genderStats.Female = count
		default:
			genderStats.Other += count
		}
	}

	// Lấy thống kê độ tuổi
	var ageStats AgeStats
	query := `
			SELECT 
				SUM(CASE WHEN TIMESTAMPDIFF(YEAR, birthday, CURDATE()) < 18 THEN 1 ELSE 0 END) as under18,
				SUM(CASE WHEN TIMESTAMPDIFF(YEAR, birthday, CURDATE()) BETWEEN 18 AND 25 THEN 1 ELSE 0 END) as age18to25,
				SUM(CASE WHEN TIMESTAMPDIFF(YEAR, birthday, CURDATE()) BETWEEN 26 AND 35 THEN 1 ELSE 0 END) as age26to35,
				SUM(CASE WHEN TIMESTAMPDIFF(YEAR, birthday, CURDATE()) BETWEEN 36 AND 45 THEN 1 ELSE 0 END) as age36to45,
				SUM(CASE WHEN TIMESTAMPDIFF(YEAR, birthday, CURDATE()) > 45 THEN 1 ELSE 0 END) as over45
			FROM users
		`
	err = config.DB.QueryRow(query).Scan(
		&ageStats.Under18,
		&ageStats.Age18to25,
		&ageStats.Age26to35,
		&ageStats.Age36to45,
		&ageStats.Over45,
	)
	if err != nil {
		http.Error(w, "Lỗi lấy thống kê độ tuổi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Lấy dữ liệu tăng trưởng 7 ngày gần nhất
	var growthLabels []string
	var growthData []int
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		dateStr := date.Format("02/01")
		growthLabels = append(growthLabels, dateStr)

		var count int
		err = config.DB.QueryRow("SELECT COUNT(*) FROM users WHERE DATE(created_at) = ?", date.Format("2006-01-02")).Scan(&count)
		if err != nil {
			http.Error(w, "Lỗi lấy dữ liệu tăng trưởng", http.StatusInternalServerError)
			return
		}
		growthData = append(growthData, count)
	}

	// Lấy danh sách log gần đây
	recentLogs, err := models.GetAdminLogs(1, 10) // Lấy 10 log gần nhất
	if err != nil {
		http.Error(w, "Lỗi lấy dữ liệu log", http.StatusInternalServerError)
		return
	}

	// Tạo dữ liệu cho template
	data := DashboardData{
		TotalUsers:    totalUsers,
		NewUsersToday: newUsersToday,
		GenderStats:   genderStats,
		AgeStats:      ageStats,
		GrowthLabels:  growthLabels,
		GrowthData:    growthData,
		RecentLogs:    recentLogs,
	}

	// Render template
	tmpl, err := template.ParseFiles("views/dashboard.html")
	if err != nil {
		http.Error(w, "Lỗi tải template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Lỗi render template", http.StatusInternalServerError)
		return
	}
}

// Hàm lấy AdminID từ session
func getAdminIDFromSession(r *http.Request) int {
	cookie, err := r.Cookie("token")
	if err != nil {
		log.Printf("Không tìm thấy cookie token: %v", err)
		return 0
	}

	tokenString := cookie.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		log.Printf("Lỗi parse token: %v", err)
		return 0
	}

	if !token.Valid {
		log.Printf("Token không hợp lệ")
		return 0
	}

	// Lấy AdminID từ username trong claims
	admin, err := models.GetAdminByUsername(claims.Username)
	if err != nil {
		log.Printf("Không tìm thấy admin với username %s: %v", claims.Username, err)
		return 0
	}

	log.Printf("Tìm thấy admin với ID: %d", admin.ID)
	return admin.ID
}
