package service

import (
	"beauty-salon/internal/models"
	"beauty-salon/internal/repository"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

//

type Service interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
	GetUserByID(id uint) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	DeleteUser(id string) error

	AddService(s *models.Service) error
	GetServices() ([]models.Service, error)
	GetService(id string) (*models.Service, error)
	DeleteService(id string) error

	AddStaff(s *models.Staff) error
	GetStaffList() ([]models.Staff, error)
	GetStaff(id string) (*models.Staff, error)
	DeleteStaff(id string) error

	CreateBooking(b *models.Booking) error
	GetBookings() ([]models.Booking, error)
	GetBooking(id string) (*models.Booking, error)
	UpdateBooking(id string, updates map[string]interface{}) (*models.Booking, error)
	CancelBooking(id string) error
}

type SalonService struct {
	repo repository.Repository
}

func NewSalonService(repo repository.Repository) *SalonService {
	return &SalonService{repo: repo}
}

func (s *SalonService) Register(username, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}
	return s.repo.CreateUser(&models.User{Username: username, Password: string(hashed)})
}

func (s *SalonService) Login(username, password string) (string, error) {
	u, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": u.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (s *SalonService) GetUserByID(id uint) (*models.User, error) { return s.repo.GetUserByID(id) }
func (s *SalonService) GetAllUsers() ([]models.User, error)       { return s.repo.GetAllUsers() }
func (s *SalonService) DeleteUser(id string) error                { return s.repo.DeleteUser(id) }

func (s *SalonService) AddService(srv *models.Service) error   { return s.repo.CreateService(srv) }
func (s *SalonService) GetServices() ([]models.Service, error) { return s.repo.GetAllServices() }
func (s *SalonService) GetService(id string) (*models.Service, error) {
	return s.repo.GetServiceByID(id)
}
func (s *SalonService) DeleteService(id string) error { return s.repo.DeleteService(id) }

func (s *SalonService) AddStaff(st *models.Staff) error           { return s.repo.CreateStaff(st) }
func (s *SalonService) GetStaffList() ([]models.Staff, error)     { return s.repo.GetAllStaff() }
func (s *SalonService) GetStaff(id string) (*models.Staff, error) { return s.repo.GetStaffByID(id) }
func (s *SalonService) DeleteStaff(id string) error               { return s.repo.DeleteStaff(id) }

func (s *SalonService) CreateBooking(b *models.Booking) error  { return s.repo.CreateBooking(b) }
func (s *SalonService) GetBookings() ([]models.Booking, error) { return s.repo.GetAllBookings() }
func (s *SalonService) GetBooking(id string) (*models.Booking, error) {
	return s.repo.GetBookingByID(id)
}
func (s *SalonService) UpdateBooking(id string, updates map[string]interface{}) (*models.Booking, error) {
	b, err := s.repo.GetBookingByID(id)
	if err != nil {
		return nil, err
	}
	if err := s.repo.UpdateBooking(b, updates); err != nil {
		return nil, err
	}
	return s.repo.GetBookingByID(id)
}
func (s *SalonService) CancelBooking(id string) error { return s.repo.DeleteBooking(id) }
