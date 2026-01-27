package repository

import (
	"beauty-salon/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	// Users
	CreateUser(u *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	DeleteUser(id string) error

	// Services
	CreateService(s *models.Service) error
	GetAllServices() ([]models.Service, error)
	GetServiceByID(id string) (*models.Service, error)
	DeleteService(id string) error

	// Staff
	CreateStaff(s *models.Staff) error
	GetAllStaff() ([]models.Staff, error)
	GetStaffByID(id string) (*models.Staff, error)
	DeleteStaff(id string) error

	// Bookings
	CreateBooking(b *models.Booking) error
	GetAllBookings() ([]models.Booking, error)
	GetBookingByID(id string) (*models.Booking, error)
	UpdateBooking(b *models.Booking, updates map[string]interface{}) error
	DeleteBooking(id string) error
}

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Users
func (r *PostgresRepository) CreateUser(u *models.User) error { return r.db.Create(u).Error }
func (r *PostgresRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}
func (r *PostgresRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return &user, err
}
func (r *PostgresRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}
func (r *PostgresRepository) DeleteUser(id string) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

// Services
func (r *PostgresRepository) CreateService(s *models.Service) error { return r.db.Create(s).Error }
func (r *PostgresRepository) GetAllServices() ([]models.Service, error) {
	var services []models.Service
	err := r.db.Find(&services).Error
	return services, err
}
func (r *PostgresRepository) GetServiceByID(id string) (*models.Service, error) {
	var service models.Service
	err := r.db.First(&service, "id = ?", id).Error
	return &service, err
}
func (r *PostgresRepository) DeleteService(id string) error {
	return r.db.Delete(&models.Service{}, "id = ?", id).Error
}

// Staffs
func (r *PostgresRepository) CreateStaff(s *models.Staff) error { return r.db.Create(s).Error }
func (r *PostgresRepository) GetAllStaff() ([]models.Staff, error) {
	var staff []models.Staff
	err := r.db.Find(&staff).Error
	return staff, err
}
func (r *PostgresRepository) GetStaffByID(id string) (*models.Staff, error) {
	var staff models.Staff
	err := r.db.First(&staff, "id = ?", id).Error
	return &staff, err
}
func (r *PostgresRepository) DeleteStaff(id string) error {
	return r.db.Delete(&models.Staff{}, "id = ?", id).Error
}

// Bookings
func (r *PostgresRepository) CreateBooking(b *models.Booking) error { return r.db.Create(b).Error }
func (r *PostgresRepository) GetAllBookings() ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Preload("User").Preload("Service").Preload("Staff").Find(&bookings).Error
	return bookings, err
}
func (r *PostgresRepository) GetBookingByID(id string) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.Preload("User").Preload("Service").Preload("Staff").First(&booking, "id = ?", id).Error
	return &booking, err
}
func (r *PostgresRepository) UpdateBooking(b *models.Booking, updates map[string]interface{}) error {
	return r.db.Model(b).Updates(updates).Error
}
func (r *PostgresRepository) DeleteBooking(id string) error {
	return r.db.Delete(&models.Booking{}, "id = ?", id).Error
}
