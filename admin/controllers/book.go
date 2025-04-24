package controllers

import (
	"encoding/json"
	"fmt"
	"go-web-admin/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// BooksHandler xử lý trang danh sách sách
func BooksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Lấy danh sách sách
		books, err := models.GetAllBooks()
		if err != nil {
			log.Printf("Lỗi khi lấy danh sách sách: %v", err)
			http.Error(w, "Lỗi khi lấy danh sách sách", http.StatusInternalServerError)
			return
		}

		// Tạo template với các hàm tùy chỉnh
		funcMap := template.FuncMap{
			"formatMoney": func(amount int) string {
				formattedAmount := fmt.Sprintf("%d", amount)
				n := len(formattedAmount)

				if n <= 3 {
					return formattedAmount
				}

				var result []byte
				for i := 0; i < n; i++ {
					if i > 0 && (n-i)%3 == 0 {
						result = append(result, ',')
					}
					result = append(result, formattedAmount[i])
				}

				return string(result)
			},
		}

		tmpl := template.Must(template.New("list.html").Funcs(funcMap).ParseFiles("views/books/list.html"))
		if err := tmpl.Execute(w, books); err != nil {
			log.Printf("Lỗi khi render template: %v", err)
			http.Error(w, "Lỗi khi render template", http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		// Tạo sách mới
		price, _ := strconv.Atoi(r.FormValue("price"))
		quantity, _ := strconv.Atoi(r.FormValue("quantity"))

		book := &models.Book{
			Title:       r.FormValue("title"),
			Author:      r.FormValue("author"),
			Publisher:   r.FormValue("publisher"),
			Category:    r.FormValue("category"),
			Description: r.FormValue("description"),
			Price:       price,
			Quantity:    quantity,
			ImageURL:    r.FormValue("image_url"),
		}

		if err := book.Create(); err != nil {
			log.Printf("Lỗi khi tạo sách: %v", err)
			http.Error(w, "Lỗi khi tạo sách", http.StatusInternalServerError)
			return
		}

		// Ghi log
		adminID, ok := r.Context().Value("adminID").(int)
		if !ok {
			log.Printf("Không thể lấy adminID từ context")
		} else {
			logEntry := models.AdminLog{
				AdminID:   adminID,
				Action:    "CREATE_BOOK",
				Details:   fmt.Sprintf("Tạo sách mới: %s", book.Title),
				IPAddress: r.RemoteAddr,
			}
			if err := logEntry.Create(); err != nil {
				log.Printf("Lỗi khi ghi log: %v", err)
			}
		}

		http.Redirect(w, r, "/admin/books", http.StatusSeeOther)

	case http.MethodDelete:
		// Xóa sách
		id, _ := strconv.Atoi(r.FormValue("id"))
		book, err := models.GetBookByID(id)
		if err != nil {
			log.Printf("Lỗi khi tìm sách: %v", err)
			http.Error(w, "Lỗi khi tìm sách", http.StatusInternalServerError)
			return
		}

		if err := models.DeleteBook(id); err != nil {
			log.Printf("Lỗi khi xóa sách: %v", err)
			http.Error(w, "Lỗi khi xóa sách", http.StatusInternalServerError)
			return
		}

		// Ghi log
		adminID := getAdminIDFromSession(r)
		log.Printf("AdminID từ session: %d", adminID)
		if adminID > 0 {
			logEntry := &models.AdminLog{
				AdminID:   adminID,
				Action:    "DELETE_BOOK",
				Details:   fmt.Sprintf("Xóa sách: %s", book.Title),
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

	default:
		http.Error(w, "Phương thức không được hỗ trợ", http.StatusMethodNotAllowed)
	}
}

// EditBookHandler xử lý trang sửa sách
func EditBookHandler(w http.ResponseWriter, r *http.Request) {
	// Lấy ID từ URL
	id, err := strconv.Atoi(r.URL.Path[len("/admin/books/edit/"):])
	if err != nil {
		log.Printf("Lỗi khi chuyển đổi ID: %v", err)
		http.Error(w, "ID không hợp lệ", http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {
		// Lấy thông tin sách
		book, err := models.GetBookByID(id)
		if err != nil {
			log.Printf("Lỗi khi lấy thông tin sách: %v", err)
			http.Error(w, "Không tìm thấy sách", http.StatusNotFound)
			return
		}

		// Render form sửa
		tmpl, err := template.ParseFiles("views/books/edit.html")
		if err != nil {
			log.Printf("Lỗi khi tải template: %v", err)
			http.Error(w, "Lỗi tải template", http.StatusInternalServerError)
			return
		}

		data := struct {
			Book *models.Book
		}{
			Book: book,
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("Lỗi khi render template: %v", err)
			http.Error(w, "Lỗi render template", http.StatusInternalServerError)
			return
		}
		return
	}

	// Xử lý cập nhật thông tin sách
	price, err := strconv.Atoi(r.FormValue("price"))
	if err != nil {
		log.Printf("Lỗi khi chuyển đổi giá: %v", err)
		http.Error(w, "Giá không hợp lệ", http.StatusBadRequest)
		return
	}

	quantity, err := strconv.Atoi(r.FormValue("quantity"))
	if err != nil {
		log.Printf("Lỗi khi chuyển đổi số lượng: %v", err)
		http.Error(w, "Số lượng không hợp lệ", http.StatusBadRequest)
		return
	}

	book := &models.Book{
		ID:          id,
		Title:       r.FormValue("title"),
		Category:    r.FormValue("category"),
		Author:      r.FormValue("author"),
		Publisher:   r.FormValue("publisher"),
		Description: r.FormValue("description"),
		Price:       price,
		Quantity:    quantity,
		ImageURL:    r.FormValue("image_url"),
	}

	// Lấy thông tin sách cũ để so sánh
	oldBook, err := models.GetBookByID(id)
	if err != nil {
		log.Printf("Lỗi khi lấy thông tin sách cũ: %v", err)
		http.Error(w, "Không tìm thấy sách", http.StatusNotFound)
		return
	}

	err = book.Update()
	if err != nil {
		log.Printf("Lỗi khi cập nhật sách: %v", err)
		http.Error(w, "Lỗi cập nhật sách", http.StatusInternalServerError)
		return
	}

	// Tạo chi tiết log
	var details string
	if oldBook.Title != book.Title {
		details += "Tên sách: " + oldBook.Title + " -> " + book.Title + "\n"
	}
	if oldBook.Category != book.Category {
		details += "Thể loại: " + oldBook.Category + " -> " + book.Category + "\n"
	}
	if oldBook.Author != book.Author {
		details += "Tác giả: " + oldBook.Author + " -> " + book.Author + "\n"
	}
	if oldBook.Publisher != book.Publisher {
		details += "Nhà xuất bản: " + oldBook.Publisher + " -> " + book.Publisher + "\n"
	}
	if oldBook.Price != book.Price {
		details += "Giá: " + strconv.Itoa(oldBook.Price) + " -> " + strconv.Itoa(book.Price) + "\n"
	}
	if oldBook.Quantity != book.Quantity {
		details += "Số lượng: " + strconv.Itoa(oldBook.Quantity) + " -> " + strconv.Itoa(book.Quantity) + "\n"
	}

	// Ghi log
	adminID := getAdminIDFromSession(r)
	log.Printf("AdminID từ session: %d", adminID)
	if adminID > 0 {
		logEntry := &models.AdminLog{
			AdminID:   adminID,
			Action:    "UPDATE_BOOK",
			Details:   "Cập nhật sách: " + book.Title + "\n" + details,
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

	http.Redirect(w, r, "/admin/books", http.StatusSeeOther)
}

// SearchBooksHandler xử lý tìm kiếm sách
func SearchBooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query().Get("q")
	var books []models.Book
	var err error

	if query == "" {
		// Nếu không có từ khóa, trả về tất cả sách
		books, err = models.GetAllBooks()
	} else {
		// Nếu có từ khóa, tìm kiếm sách
		books, err = models.SearchBooks(query)
	}

	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Lỗi khi tìm kiếm sách"})
		return
	}

	json.NewEncoder(w).Encode(books)
}
