package models

import "time"

type Book struct {
	Title       string     `json:"title"`
	CategoryId  int        `json:"categoryId"`
	Description string     `json:"description"`
	ImageUrl    string     `json:"imageUrl"`
	ReleaseYear int        `json:"releaseYear"`
	Price       int        `json:"price"`
	TotalPage   int        `json:"totalPage"`
	Thickness   string     `json:"thickness"`
	CreatedAt   time.Time  `json:"createdAt"`
	CreatedBy   string     `json:"createdBy"`
	ModifiedAt  *time.Time `json:"modifiedAt"`
	ModifiedBy  *string    `json:"modifiedBy"`
}
