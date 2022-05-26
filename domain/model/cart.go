package model

type Cart struct {
	ID        int64 `gorm:"primary_key;not_null;auto_increment" json:"id"`
	ProductID int64 `gorm:"not_null" json:"product_id"`
	Num       int64 `gorm:"not_null" json:"num"`
	SizeID    int64 `gorm:"not_null" json:"size_ud"`
	UserID    int64 `gorm:"not_null" json:"user_ud"`
}
