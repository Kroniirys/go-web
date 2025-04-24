package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"go-web-admin/config"
	"go-web-admin/models"
)

func main() {
	// Kết nối database
	err := config.InitDB()
	if err != nil {
		println("Lỗi kết nối database:", err.Error())
		return
	}
	defer config.DB.Close()

	// Khởi tạo random seed
	rand.Seed(time.Now().UnixNano())

	// Số lượng đơn hàng cần tạo
	numOrders := 10

	// Lấy danh sách user
	users, err := models.GetAllUsers()
	if err != nil {
		log.Fatal(err)
	}

	// Lấy danh sách sách
	books, err := models.GetAllBooks()
	if err != nil {
		log.Fatal(err)
	}

	// Tạo các đơn hàng
	for i := 0; i < numOrders; i++ {
		// Chọn ngẫu nhiên user
		user := users[rand.Intn(len(users))]

		// Tạo đơn hàng
		order := models.Order{
			UserID:      user.ID,
			TotalAmount: 0,
			Status:      getRandomStatus(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Lưu đơn hàng
		result, err := config.DB.Exec("INSERT INTO orders (user_id, total_amount, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
			order.UserID, order.TotalAmount, order.Status, order.CreatedAt, order.UpdatedAt)
		if err != nil {
			log.Fatal(err)
		}

		// Lấy ID của đơn hàng vừa tạo
		orderID, err := result.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}

		// Số lượng sản phẩm trong đơn hàng (1-5 sản phẩm)
		numItems := rand.Intn(5) + 1
		totalAmount := 0.0

		// Thêm các sản phẩm vào đơn hàng
		for j := 0; j < numItems; j++ {
			// Chọn ngẫu nhiên sách
			book := books[rand.Intn(len(books))]

			// Số lượng ngẫu nhiên (1-3)
			quantity := rand.Intn(3) + 1

			// Tính thành tiền
			subtotal := book.Price * int(quantity)
			totalAmount += float64(subtotal)

			// Thêm vào chi tiết đơn hàng
			_, err := config.DB.Exec("INSERT INTO order_items (order_id, book_id, quantity, price, created_at) VALUES (?, ?, ?, ?, ?)",
				orderID, book.ID, quantity, book.Price, time.Now())
			if err != nil {
				log.Fatal(err)
			}
		}

		// Cập nhật tổng tiền đơn hàng
		_, err = config.DB.Exec("UPDATE orders SET total_amount = ? WHERE id = ?", totalAmount, orderID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Đã tạo đơn hàng #%d cho user %s với tổng tiền %.2f VNĐ\n",
			orderID, user.Username, totalAmount)
	}

	fmt.Println("Đã tạo xong các đơn hàng ngẫu nhiên!")
}

// Hàm trả về trạng thái ngẫu nhiên
func getRandomStatus() string {
	statuses := []string{"pending", "processing", "completed", "cancelled"}
	return statuses[rand.Intn(len(statuses))]
}
