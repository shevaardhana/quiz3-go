package models

import "time"

type Category struct {
	Name       string     `json:"title"`
	CreatedAt  time.Time  `json:"createdAt"`
	CreatedBy  string     `json:"createdBy"`
	ModifiedAt *time.Time `json:"modifiedAt"`
	ModifiedBy *string    `json:"modifiedBy"`
}
