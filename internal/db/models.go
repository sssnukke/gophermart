package db

import "time"

type User struct {
	ID           uint   `json:"-" gorm:"primaryKey"`
	Login        string `json:"login" gorm:"uniqueIndex;not null"`
	PasswordHash string `json:"-" gorm:"not null"`
	CreatedAt    time.Time
	Order        []Order `json:"orders"`
	Balance      float64 `json:"balance" gorm:"default:0"`
	Withdrawn    float64 `json:"withdrawn" gorm:"default:0"`
}

type Order struct {
	ID         uint      `json:"-" gorm:"primaryKey"`
	Number     string    `json:"number" gorm:"uniqueIndex;not null"`
	UserID     uint      `json:"-" gorm:"index;not null"`
	Status     string    `json:"status" gorm:"not null"`
	Accrual    *float64  `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at" gorm:"autoCreateTime"`
}

type Withdrawal struct {
	ID          uint      `json:"-" gorm:"primaryKey"`
	UserID      uint      `json:"-" gorm:"index"`
	OrderNumber string    `json:"order_number" gorm:"not null"`
	Sum         float64   `json:"sum" gorm:"not null"`
	ProcessedAt time.Time `json:"processed_at" gorm:"autoCreateTime"`
}
