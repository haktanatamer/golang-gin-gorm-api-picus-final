package category

import "time"

type Category struct {
	Id                uint `gorm:"primaryKey"`
	Name              string
	IsDeleted, Active bool      `json:"-"`
	CreatedAt         time.Time `gorm:"<-:create" json:"-"`
	UpdatedAt         time.Time `json:"-"`
}
