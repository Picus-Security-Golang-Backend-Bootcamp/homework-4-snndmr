package book

import (
	"errors"
	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-snndmr/constants"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Migration() {
	err := r.db.AutoMigrate(&Book{})
	if err != nil {
		return
	}
}

func (r *Repository) InitializeWithSampleData(books chan *Book) {
	for book := range books {
		r.db.Where(Book{Title: book.Title}).FirstOrCreate(&book)
	}
}

func (r *Repository) GetById(id int) (error, Book) {
	var book Book
	result := r.db.First(&book, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return constants.ErrBookNotFound, Book{}
	}
	return nil, book
}

func (r *Repository) GetBooks() []Book {
	var books []Book
	r.db.Find(&books)
	return books
}

func (r *Repository) Search(substr string) []Book {
	var books []Book
	substr = strings.ToLower(substr)
	r.db.Preload(clause.Associations).Joins("JOIN Authors on Authors.id = Books.author_id").
		Where("lower(Books.title) LIKE ?", "%"+substr+"%").
		Or("lower(Books.stock_id) LIKE ?", "%"+substr+"%").
		Or("lower(Authors.name) LIKE ?", "%"+substr+"%").
		Find(&books)
	return books
}

func (r *Repository) Create(book *Book) error {
	result := r.db.Create(book)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *Repository) Update(book Book) error {
	result := r.db.Save(book)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
