package domain

import "time"

type Category struct {
	ID   string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"ID"`
	Name string `gorm:"size:100;unique;not null" json:"Name"`
}

type Product struct {
	ID          string   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"ID"`
	ProductName string   `gorm:"size:150;not null" json:"ProductName"`
	Description string   `gorm:"type:text;null" json:"Description"`
	Price       float64  `gorm:"not null" json:"Price"`
	Quantity    int      `gorm:"not null" json:"Quantity"`
	CategoryID  string   `json:"CategoryID"`
	IsActive    bool     `gorm:"default:true" json:"IsActive"`
	Category    Category `gorm:"foreignKey:CategoryID" json:"Category"`
}

type Order struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"ID"`
	ProductID string    `json:"ProductID"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"Product"`
	OrderDate time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"OrderDate"`
	Quantity  int       `gorm:"not null" json:"Quantity"`
}
