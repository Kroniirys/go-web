package models

import (
	"database/sql"
	"errors"
	"go-web-admin/config"
	"time"
)

type Order struct {
	ID          int
	UserID      int
	TotalAmount int
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OrderItem struct {
	ID        int
	OrderID   int
	BookID    int
	Quantity  int
	Price     int
	CreatedAt time.Time
}

// GetAllOrders lấy tất cả đơn hàng
func GetAllOrders() ([]Order, error) {
	rows, err := config.DB.Query("SELECT id, user_id, total_amount, status, created_at, updated_at FROM orders ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		err := rows.Scan(&o.ID, &o.UserID, &o.TotalAmount, &o.Status, &o.CreatedAt, &o.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

// GetOrderByID lấy đơn hàng theo ID
func GetOrderByID(id int) (*Order, error) {
	var o Order
	err := config.DB.QueryRow("SELECT id, user_id, total_amount, status, created_at, updated_at FROM orders WHERE id = ?", id).
		Scan(&o.ID, &o.UserID, &o.TotalAmount, &o.Status, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &o, nil
}

// GetOrderItems lấy chi tiết đơn hàng
func GetOrderItems(orderID int) ([]OrderItem, error) {
	rows, err := config.DB.Query("SELECT id, order_id, book_id, quantity, price, created_at FROM order_items WHERE order_id = ?", orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []OrderItem
	for rows.Next() {
		var item OrderItem
		err := rows.Scan(&item.ID, &item.OrderID, &item.BookID, &item.Quantity, &item.Price, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// UpdateOrderStatus cập nhật trạng thái đơn hàng
func UpdateOrderStatus(id int, status string) error {
	_, err := config.DB.Exec("UPDATE orders SET status = ?, updated_at = NOW() WHERE id = ?", status, id)
	return err
}
