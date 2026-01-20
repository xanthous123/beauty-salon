package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Password string `json:"-"`
	Role     string `gorm:"default:client" json:"role"` // client, admin
}

type Service struct {
	gorm.Model
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	DurationMin int     `json:"duration_min"` // Длительность процедуры в минутах
	Description string  `json:"description"`
}

type Staff struct {
	gorm.Model
	FullName   string `json:"full_name"`
	Speciality string `json:"speciality"` // Например: "Топ-стилист", "Нейл-мастер"
}

type Booking struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	ServiceID uint   `json:"service_id"`
	StaffID   uint   `json:"staff_id"`
	Date      string `json:"date"`                          // YYYY-MM-DD HH:MM
	Status    string `gorm:"default:pending" json:"status"` // pending, confirmed, cancelled

	User    User    `gorm:"foreignKey:UserID" json:"user"`
	Service Service `gorm:"foreignKey:ServiceID" json:"service"`
	Staff   Staff   `gorm:"foreignKey:StaffID" json:"staff"`
}
