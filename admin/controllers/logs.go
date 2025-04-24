package controllers

import (
	"go-web-admin/models"
	"html/template"
	"net/http"
)

type LogsData struct {
	RecentLogs []models.AdminLog
}

func LogsHandler(w http.ResponseWriter, r *http.Request) {
	// Lấy danh sách log gần đây
	recentLogs, err := models.GetAdminLogs(1, 50)
	if err != nil {
		http.Error(w, "Lỗi lấy dữ liệu log", http.StatusInternalServerError)
		return
	}

	// Tạo dữ liệu cho template
	data := LogsData{
		RecentLogs: recentLogs,
	}

	// Render template
	tmpl, err := template.ParseFiles("views/logs.html")
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
