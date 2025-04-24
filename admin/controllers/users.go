package controllers

import (
	"encoding/json"
	"fmt"
	"go-web-admin/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type UserListData struct {
	Users       []models.User
	SearchQuery string
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Lấy từ khóa tìm kiếm
		searchQuery := r.URL.Query().Get("search")

		// Lấy danh sách người dùng
		var users []models.User
		var err error

		if searchQuery != "" {
			// Tìm kiếm người dùng
			users, err = models.SearchUsers(searchQuery)
		} else {
			// Lấy tất cả người dùng
			users, err = models.GetAllUsers()
		}

		if err != nil {
			http.Error(w, "Lỗi lấy danh sách người dùng", http.StatusInternalServerError)
			return
		}

		// Tạo dữ liệu cho template
		data := UserListData{
			Users:       users,
			SearchQuery: searchQuery,
		}

		// Render template
		tmpl, err := template.New("list.html").Funcs(template.FuncMap{
			"formatDate": func(date string) string {
				t, err := time.Parse("2006-01-02", date)
				if err != nil {
					return date
				}
				return t.Format("02/01/2006")
			},
		}).ParseFiles("views/users/list.html")
		if err != nil {
			http.Error(w, "Lỗi tải template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Lỗi render template", http.StatusInternalServerError)
			return
		}
		return
	}

	// Xử lý các thao tác với người dùng
	switch r.Method {
	case "POST":
		// Thêm người dùng mới
		user := &models.User{
			Username: r.FormValue("username"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
			Birthday: r.FormValue("birthday"),
			Gender:   r.FormValue("gender"),
			CCCD:     r.FormValue("cccd"),
		}
		err := user.Create()
		if err != nil {
			http.Error(w, "Lỗi tạo người dùng", http.StatusInternalServerError)
			return
		}

		// Ghi log thêm người dùng
		adminID := getAdminIDFromSession(r)
		log.Printf("AdminID từ session: %d", adminID)
		if adminID > 0 {
			logEntry := &models.AdminLog{
				AdminID:   adminID,
				Action:    "create_user",
				Details:   fmt.Sprintf("Thêm người dùng mới: %s", user.Username),
				IPAddress: r.RemoteAddr,
			}
			log.Printf("Chuẩn bị ghi log: %+v", logEntry)
			if err := logEntry.Create(); err != nil {
				log.Printf("Lỗi ghi log: %v", err)
			} else {
				log.Printf("Ghi log thành công")
			}
		} else {
			log.Printf("Không tìm thấy adminID")
		}

	case "DELETE":
		// Xóa người dùng
		id, _ := strconv.Atoi(r.FormValue("id"))
		user, err := models.GetUserByID(id)
		if err != nil {
			http.Error(w, "Lỗi tìm người dùng", http.StatusInternalServerError)
			return
		}

		err = models.DeleteUser(id)
		if err != nil {
			http.Error(w, "Lỗi xóa người dùng", http.StatusInternalServerError)
			return
		}

		// Ghi log xóa người dùng
		adminID := getAdminIDFromSession(r)
		log.Printf("AdminID từ session: %d", adminID)
		if adminID > 0 {
			logEntry := &models.AdminLog{
				AdminID:   adminID,
				Action:    "delete_user",
				Details:   fmt.Sprintf("Xóa người dùng: %s", user.Username),
				IPAddress: r.RemoteAddr,
			}
			log.Printf("Chuẩn bị ghi log: %+v", logEntry)
			if err := logEntry.Create(); err != nil {
				log.Printf("Lỗi ghi log: %v", err)
			} else {
				log.Printf("Ghi log thành công")
			}
		} else {
			log.Printf("Không tìm thấy adminID")
		}
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func EditUserHandler(w http.ResponseWriter, r *http.Request) {
	// Lấy ID từ URL
	id, err := strconv.Atoi(r.URL.Path[len("/admin/users/edit/"):])
	if err != nil {
		http.Error(w, "ID không hợp lệ", http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {
		// Lấy thông tin người dùng
		user, err := models.GetUserByID(id)
		if err != nil {
			http.Error(w, "Không tìm thấy người dùng", http.StatusNotFound)
			return
		}

		// Render form sửa
		tmpl, err := template.ParseFiles("views/users/edit.html")
		if err != nil {
			http.Error(w, "Lỗi tải template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, user)
		if err != nil {
			http.Error(w, "Lỗi render template", http.StatusInternalServerError)
			return
		}
		return
	}

	// Xử lý cập nhật thông tin người dùng
	user := &models.User{
		ID:       id,
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Birthday: r.FormValue("birthday"),
		Gender:   r.FormValue("gender"),
		CCCD:     r.FormValue("cccd"),
	}

	// Lấy thông tin người dùng cũ để so sánh
	oldUser, err := models.GetUserByID(id)
	if err != nil {
		http.Error(w, "Không tìm thấy người dùng", http.StatusNotFound)
		return
	}

	// Chỉ cập nhật mật khẩu nếu có giá trị
	if password := r.FormValue("password"); password != "" {
		user.Password = password
	}

	err = user.Update()
	if err != nil {
		http.Error(w, "Lỗi cập nhật người dùng", http.StatusInternalServerError)
		return
	}

	// Tạo chi tiết log
	var details strings.Builder
	details.WriteString(fmt.Sprintf("Cập nhật thông tin người dùng: %s\n", user.Username))

	if oldUser.Username != user.Username {
		details.WriteString(fmt.Sprintf("- Tên đăng nhập: %s -> %s\n", oldUser.Username, user.Username))
	}
	if oldUser.Email != user.Email {
		details.WriteString(fmt.Sprintf("- Email: %s -> %s\n", oldUser.Email, user.Email))
	}
	if oldUser.Birthday != user.Birthday {
		details.WriteString(fmt.Sprintf("- Ngày sinh: %s -> %s\n", oldUser.Birthday, user.Birthday))
	}
	if oldUser.Gender != user.Gender {
		details.WriteString(fmt.Sprintf("- Giới tính: %s -> %s\n", oldUser.Gender, user.Gender))
	}
	if oldUser.CCCD != user.CCCD {
		details.WriteString(fmt.Sprintf("- CCCD: %s -> %s\n", oldUser.CCCD, user.CCCD))
	}
	if r.FormValue("password") != "" {
		details.WriteString("- Mật khẩu đã được thay đổi\n")
	}

	// Ghi log cập nhật người dùng
	adminID := getAdminIDFromSession(r)
	log.Printf("AdminID từ session: %d", adminID)
	if adminID > 0 {
		logEntry := &models.AdminLog{
			AdminID:   adminID,
			Action:    "update_user",
			Details:   details.String(),
			IPAddress: r.RemoteAddr,
		}
		log.Printf("Chuẩn bị ghi log: %+v", logEntry)
		if err := logEntry.Create(); err != nil {
			log.Printf("Lỗi ghi log: %v", err)
		} else {
			log.Printf("Ghi log thành công")
		}
	} else {
		log.Printf("Không tìm thấy adminID")
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

// SearchUsersHandler xử lý tìm kiếm người dùng
func SearchUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Đặt header Content-Type là application/json
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query().Get("q")
	if query == "" {
		json.NewEncoder(w).Encode(map[string]string{"error": "Query parameter is required"})
		return
	}

	users, err := models.SearchUsers(query)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Error searching users"})
		return
	}

	// Trả về kết quả dưới dạng JSON
	json.NewEncoder(w).Encode(users)
}
