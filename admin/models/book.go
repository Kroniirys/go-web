package models

import (
	"database/sql"
	"errors"
	"go-web-admin/config"
	"time"
)

type Book struct {
	ID          int
	Title       string
	Category    string
	Author      string
	Publisher   string
	Description string
	Price       int
	Quantity    int
	ImageURL    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GetAllBooks lấy tất cả sách
func GetAllBooks() ([]Book, error) {
	rows, err := config.DB.Query("SELECT id, title, category, author, publisher, description, price, quantity, image_url, created_at, updated_at FROM books ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		err := rows.Scan(&b.ID, &b.Title, &b.Category, &b.Author, &b.Publisher, &b.Description, &b.Price, &b.Quantity, &b.ImageURL, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

// GetBooksByCategory lấy sách theo thể loại
func GetBooksByCategory(category string) ([]Book, error) {
	rows, err := config.DB.Query("SELECT id, title, category, author, publisher, description, price, quantity, image_url, created_at, updated_at FROM books WHERE category = ? ORDER BY created_at DESC", category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		err := rows.Scan(&b.ID, &b.Title, &b.Category, &b.Author, &b.Publisher, &b.Description, &b.Price, &b.Quantity, &b.ImageURL, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

// GetBookByID lấy sách theo ID
func GetBookByID(id int) (*Book, error) {
	var b Book
	err := config.DB.QueryRow("SELECT id, title, category, author, publisher, description, price, quantity, image_url, created_at, updated_at FROM books WHERE id = ?", id).
		Scan(&b.ID, &b.Title, &b.Category, &b.Author, &b.Publisher, &b.Description, &b.Price, &b.Quantity, &b.ImageURL, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("book not found")
		}
		return nil, err
	}
	return &b, nil
}

// Create thêm sách mới
func (b *Book) Create() error {
	result, err := config.DB.Exec("INSERT INTO books (title, category, author, publisher, description, price, quantity, image_url) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		b.Title, b.Category, b.Author, b.Publisher, b.Description, b.Price, b.Quantity, b.ImageURL)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	b.ID = int(id)
	return nil
}

// Update cập nhật thông tin sách
func (b *Book) Update() error {
	_, err := config.DB.Exec("UPDATE books SET title = ?, category = ?, author = ?, publisher = ?, description = ?, price = ?, quantity = ?, image_url = ? WHERE id = ?",
		b.Title, b.Category, b.Author, b.Publisher, b.Description, b.Price, b.Quantity, b.ImageURL, b.ID)
	return err
}

// Delete xóa sách
func DeleteBook(id int) error {
	_, err := config.DB.Exec("DELETE FROM books WHERE id = ?", id)
	return err
}

// SearchBooks tìm kiếm sách
func SearchBooks(query string) ([]Book, error) {
	rows, err := config.DB.Query("SELECT id, title, category, author, publisher, description, price, quantity, image_url, created_at, updated_at FROM books WHERE title LIKE ? OR author LIKE ? OR publisher LIKE ? OR category LIKE ? ORDER BY created_at DESC",
		"%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		err := rows.Scan(&b.ID, &b.Title, &b.Category, &b.Author, &b.Publisher, &b.Description, &b.Price, &b.Quantity, &b.ImageURL, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}
