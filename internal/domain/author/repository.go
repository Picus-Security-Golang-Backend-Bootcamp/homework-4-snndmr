package author

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Migration() {
	err := r.db.AutoMigrate(&Author{})
	if err != nil {
		return
	}
}

func (r *Repository) InitializeWithSampleData(authors chan *Author) {
	for author := range authors {
		r.db.Where(Author{Name: author.Name}).FirstOrCreate(&author)
	}
}
