package author

type Author struct {
	ID   uint32 `gorm:"primary_key" json:"id,omitempty"`
	Name string `gorm:"type:varchar(100)" json:"name,omitempty"`
}

func New(name string) *Author {
	return &Author{Name: name}
}

func (Author) TableName() string {
	return "authors"
}
