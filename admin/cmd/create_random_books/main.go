package main

import (
	"database/sql"
	"fmt"
	"go-web-admin/config"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Book struct {
	Title       string
	Category    string
	Author      string
	Publisher   string
	Description string
	Price       float64
	Quantity    int
	ImageURL    string
}

var (
	categories = []string{
		"Tiểu thuyết",
		"Khoa học",
		"Kinh tế",
		"Văn học",
		"Lịch sử",
		"Tâm lý học",
		"Kỹ năng sống",
		"Thiếu nhi",
	}

	authors = []string{
		"Nguyễn Nhật Ánh",
		"Nguyễn Ngọc Tư",
		"Nguyễn Du",
		"Nam Cao",
		"Vũ Trọng Phụng",
		"Nguyễn Huy Thiệp",
		"Bảo Ninh",
		"Nguyễn Bình Phương",
		"Nguyễn Việt Hà",
		"Đỗ Bích Thúy",
	}

	publishers = []string{
		"NXB Trẻ",
		"NXB Kim Đồng",
		"NXB Văn học",
		"NXB Hội Nhà văn",
		"NXB Tổng hợp TP.HCM",
		"NXB Phụ nữ",
		"NXB Thanh niên",
		"NXB Lao động",
	}

	descriptions = []string{
		"Một tác phẩm đặc sắc về cuộc sống và tình yêu",
		"Cuốn sách mang đến những bài học quý giá về cuộc sống",
		"Tác phẩm kinh điển được nhiều thế hệ yêu thích",
		"Cuốn sách giúp bạn khám phá bản thân và thế giới xung quanh",
		"Một câu chuyện cảm động về tình người",
		"Tác phẩm đạt nhiều giải thưởng văn học",
		"Cuốn sách best-seller được dịch ra nhiều thứ tiếng",
		"Tác phẩm để đời của một nhà văn nổi tiếng",
	}

	imageURLs = []string{
		"https://images.unsplash.com/photo-1544947950-fa07a98d237f",
		"https://images.unsplash.com/photo-1541963463532-d68292c34b19",
		"https://images.unsplash.com/photo-1543002588-bfa74002ed7e",
		"https://images.unsplash.com/photo-1544716278-ca5e3f4abd8c",
		"https://images.unsplash.com/photo-1512820790803-83ca734da794",
		"https://images.unsplash.com/photo-1519682337058-a94d519337bc",
		"https://images.unsplash.com/photo-1516979187457-637abb4f9353",
		"https://images.unsplash.com/photo-1513475382585-bd139b7e8857",
	}
)

func main() {
	// Kết nối database
	err := config.InitDB()
	if err != nil {
		fmt.Println("Lỗi kết nối database:", err)
		return
	}
	defer config.DB.Close()

	// Kiểm tra kết nối
	err = config.DB.Ping()
	if err != nil {
		fmt.Println("Lỗi ping database:", err)
		return
	}

	// Tạo 50 sách ngẫu nhiên
	for i := 0; i < 50; i++ {
		book := generateRandomBook()
		err := insertBook(config.DB, book)
		if err != nil {
			fmt.Printf("Lỗi khi thêm sách %d: %v\n", i+1, err)
			continue
		}
		fmt.Printf("Đã thêm sách %d: %s\n", i+1, book.Title)
	}

	fmt.Println("Hoàn thành tạo sách ngẫu nhiên!")
}

func generateRandomBook() Book {
	rand.Seed(time.Now().UnixNano())

	// Tạo tiêu đề ngẫu nhiên
	title := fmt.Sprintf("Sách %d", rand.Intn(1000)+1)

	// Chọn ngẫu nhiên các thông tin
	category := categories[rand.Intn(len(categories))]
	author := authors[rand.Intn(len(authors))]
	publisher := publishers[rand.Intn(len(publishers))]
	description := descriptions[rand.Intn(len(descriptions))]
	imageURL := imageURLs[rand.Intn(len(imageURLs))]

	// Tạo giá và số lượng ngẫu nhiên
	price := float64(rand.Intn(500)+50) * 1000 // Giá từ 50.000 đến 550.000
	quantity := rand.Intn(100) + 1             // Số lượng từ 1 đến 100

	return Book{
		Title:       title,
		Category:    category,
		Author:      author,
		Publisher:   publisher,
		Description: description,
		Price:       price,
		Quantity:    quantity,
		ImageURL:    imageURL,
	}
}

func insertBook(db *sql.DB, book Book) error {
	query := `INSERT INTO books (title, category, author, publisher, description, price, quantity, image_url) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := db.Exec(query,
		book.Title,
		book.Category,
		book.Author,
		book.Publisher,
		book.Description,
		book.Price,
		book.Quantity,
		book.ImageURL,
	)

	return err
}
