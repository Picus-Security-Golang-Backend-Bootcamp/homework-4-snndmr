package book

import (
	"fmt"
	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-snndmr/constants"
)

type Book struct {
	ID         uint32  `gorm:"primary_key" json:"id,omitempty"`
	Title      string  `gorm:"type:varchar(100)" json:"title,omitempty"`
	StockId    string  `gorm:"type:varchar(100)" json:"stock_id,omitempty"`
	ISBN       string  `gorm:"type:varchar(20)" json:"isbn,omitempty"`
	PageCount  int     `json:"page_count,omitempty"`
	StockCount int     `json:"stock_count,omitempty"`
	Price      float64 `json:"price,omitempty"`
	IsDeleted  bool    `json:"is_deleted,omitempty"`
	AuthorID   uint32  `gorm:"foreignKey:AuthorID" json:"author_id,omitempty"`
}

func NewBook(title, stockId, isbn string, pageCount, stockCount int, price float64, isDeleted bool, author uint32) *Book {
	return &Book{
		Title:      title,
		StockId:    stockId,
		ISBN:       isbn,
		PageCount:  pageCount,
		StockCount: stockCount,
		Price:      price,
		IsDeleted:  isDeleted,
		AuthorID:   author,
	}
}

func (Book) TableName() string {
	return "books"
}

type Deletable interface {
	Delete() error
}

func (book *Book) Delete() error {
	if book.IsDeleted {
		return constants.ErrBookAlreadyDeleted
	}

	book.IsDeleted = true
	return nil
}

func (book *Book) DecreaseAmount(amount int) error {
	if amount < 0 {
		return constants.ErrNegativeAmount
	}

	if amount > book.StockCount {
		return constants.ErrBookOutOfStock
	}

	book.StockCount -= amount
	return nil
}

func (book *Book) ToString() string {
	return fmt.Sprintf(
		"Title: %s Stock ID: %s ISBN: %s Page Count: %d Stock Count: %d Price: %f Author: %v", book.Title, book.StockId,
		book.ISBN, book.PageCount, book.StockCount, book.Price, book.AuthorID,
	)
}
