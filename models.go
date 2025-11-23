package main

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex" json:"username"`
	Password  string    `json:"password"`
	Token     string    `gorm:"index" json:"token"`
	CartID    *uint     `json:"cart_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Item struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Cart struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	UserID    uint        `json:"user_id"`
	Name      string      `json:"name"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	Items     []CartItem  `json:"items"`
}

type CartItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CartID    uint      `json:"cart_id"`
	ItemID    uint      `json:"item_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Order struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CartID    uint      `json:"cart_id"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}