package DB

import "gorm.io/gorm"

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	List  any `json:"list"`
}

func Paginate(page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page < 0 {
			page = 0
		}
		if limit < 0 {
			limit = 10
		}
		if limit > 100 {
			limit = 100
		}
		return db.Offset(page * limit).Limit(limit)
	}
}
