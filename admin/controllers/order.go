package controllers

import (
	"fmt"
	"go-web-admin/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// OrdersHandler xử lý trang danh sách đơn hàng
func OrdersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Lấy danh sách đơn hàng
		orders, err := models.GetAllOrders()
		if err != nil {
			log.Printf("Lỗi khi lấy danh sách đơn hàng: %v", err)
			http.Error(w, "Lỗi khi lấy danh sách đơn hàng", http.StatusInternalServerError)
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

		tmpl := template.Must(template.New("list.html").Funcs(funcMap).ParseFiles("views/orders/list.html"))
		if err := tmpl.Execute(w, orders); err != nil {
			log.Printf("Lỗi khi render template: %v", err)
			http.Error(w, "Lỗi khi render template", http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		// Cập nhật trạng thái đơn hàng
		orderID, _ := strconv.Atoi(r.FormValue("id"))
		status := r.FormValue("status")

		if err := models.UpdateOrderStatus(orderID, status); err != nil {
			log.Printf("Lỗi khi cập nhật trạng thái đơn hàng: %v", err)
			http.Error(w, "Lỗi khi cập nhật trạng thái đơn hàng", http.StatusInternalServerError)
			return
		}

		// Ghi log
		adminID, ok := r.Context().Value("adminID").(int)
		if !ok {
			log.Printf("Không thể lấy adminID từ context")
		} else {
			logEntry := models.AdminLog{
				AdminID:   adminID,
				Action:    "UPDATE_ORDER",
				Details:   fmt.Sprintf("Cập nhật trạng thái đơn hàng #%d thành %s", orderID, status),
				IPAddress: r.RemoteAddr,
			}
			if err := logEntry.Create(); err != nil {
				log.Printf("Lỗi khi ghi log: %v", err)
			}
		}

		http.Redirect(w, r, "/admin/orders", http.StatusSeeOther)

	default:
		http.Error(w, "Phương thức không được hỗ trợ", http.StatusMethodNotAllowed)
	}
}

// OrderDetailHandler xử lý trang chi tiết đơn hàng
func OrderDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Lấy ID từ URL
	id, err := strconv.Atoi(r.URL.Path[len("/admin/orders/detail/"):])
	if err != nil {
		log.Printf("Lỗi khi chuyển đổi ID: %v", err)
		http.Error(w, "ID không hợp lệ", http.StatusBadRequest)
		return
	}

	// Lấy thông tin đơn hàng
	order, err := models.GetOrderByID(id)
	if err != nil {
		log.Printf("Lỗi khi lấy thông tin đơn hàng: %v", err)
		http.Error(w, "Không tìm thấy đơn hàng", http.StatusNotFound)
		return
	}

	// Lấy chi tiết đơn hàng
	items, err := models.GetOrderItems(id)
	if err != nil {
		log.Printf("Lỗi khi lấy chi tiết đơn hàng: %v", err)
		http.Error(w, "Lỗi khi lấy chi tiết đơn hàng", http.StatusInternalServerError)
		return
	}

	// Lấy thông tin sách cho từng item
	type OrderItemWithBook struct {
		models.OrderItem
		BookTitle string
	}

	var itemsWithBooks []OrderItemWithBook
	for _, item := range items {
		book, err := models.GetBookByID(item.BookID)
		if err != nil {
			log.Printf("Lỗi khi lấy thông tin sách: %v", err)
			continue
		}

		itemsWithBooks = append(itemsWithBooks, OrderItemWithBook{
			OrderItem: item,
			BookTitle: book.Title,
		})
	}

	// Tạo template với các hàm tùy chỉnh
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"multiply": func(a int, b int) int {
			return a * b
		},
		"formatMoney": func(amount int) string {
			formattedAmount := fmt.Sprintf("%d", amount)
			n := len(formattedAmount)

			// Thêm dấu phẩy
			if n <= 3 {
				return formattedAmount
			}

			// Xử lý từng nhóm 3 số
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

	// Render template
	tmpl := template.Must(template.New("detail.html").Funcs(funcMap).ParseFiles("views/orders/detail.html"))
	data := struct {
		Order *models.Order
		Items []OrderItemWithBook
	}{
		Order: order,
		Items: itemsWithBooks,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Lỗi khi render template: %v", err)
		http.Error(w, "Lỗi khi render template", http.StatusInternalServerError)
		return
	}
}
