package service

import (
	"beauty-salon/internal/models"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) CreateUser(u *models.User) error { return m.Called(u).Error(0) }
func (m *MockRepo) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockRepo) GetUserByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockRepo) GetAllUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}
func (m *MockRepo) DeleteUser(id string) error            { return m.Called(id).Error(0) }
func (m *MockRepo) CreateService(s *models.Service) error { return m.Called(s).Error(0) }
func (m *MockRepo) GetAllServices() ([]models.Service, error) {
	args := m.Called()
	return args.Get(0).([]models.Service), args.Error(1)
}
func (m *MockRepo) GetServiceByID(id string) (*models.Service, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Service), args.Error(1)
}
func (m *MockRepo) DeleteService(id string) error     { return m.Called(id).Error(0) }
func (m *MockRepo) CreateStaff(s *models.Staff) error { return m.Called(s).Error(0) }
func (m *MockRepo) GetAllStaff() ([]models.Staff, error) {
	args := m.Called()
	return args.Get(0).([]models.Staff), args.Error(1)
}
func (m *MockRepo) GetStaffByID(id string) (*models.Staff, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Staff), args.Error(1)
}
func (m *MockRepo) DeleteStaff(id string) error           { return m.Called(id).Error(0) }
func (m *MockRepo) CreateBooking(b *models.Booking) error { return m.Called(b).Error(0) }
func (m *MockRepo) GetAllBookings() ([]models.Booking, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]models.Booking), args.Error(1)
}
func (m *MockRepo) GetBookingByID(id string) (*models.Booking, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Booking), args.Error(1)
}
func (m *MockRepo) UpdateBooking(b *models.Booking, updates map[string]interface{}) error {
	return m.Called(b, updates).Error(0)
}
func (m *MockRepo) DeleteBooking(id string) error { return m.Called(id).Error(0) }

// --- ТЕСТЫ ---

func TestRegister(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockRepo)
		svc := NewSalonService(mockRepo)

		mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil).Once()

		err := svc.Register("testuser", "password123")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo := new(MockRepo)
		svc := NewSalonService(mockRepo)

		mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).
			Return(errors.New("user already exists")).Once()

		err := svc.Register("testuser", "password123")

		assert.Error(t, err)
		assert.Equal(t, "user already exists", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Bcrypt Error", func(t *testing.T) {
		mockRepo := new(MockRepo)
		svc := NewSalonService(mockRepo)

		longPass := make([]byte, 80)
		_ = svc.Register("testuser", string(longPass))
	})
}

func TestLogin(t *testing.T) {
	os.Setenv("JWT_SECRET", "secret")
	mockRepo := new(MockRepo)
	svc := NewSalonService(mockRepo)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("pass"), 10)
	user := &models.User{Username: "admin", Password: string(hashed)}
	user.ID = 1

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetUserByUsername", "admin").Return(user, nil).Once()
		token, err := svc.Login("admin", "pass")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockRepo.On("GetUserByUsername", "nonexistent").Return(nil, errors.New("not found")).Once()
		_, err := svc.Login("nonexistent", "pass")
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		mockRepo.On("GetUserByUsername", "admin").Return(user, nil).Once()
		_, err := svc.Login("admin", "wrong")
		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
	})
}

// Тесты пользователей
func TestUserMethods(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewSalonService(mockRepo)

	t.Run("GetUserByID", func(t *testing.T) {
		mockRepo.On("GetUserByID", uint(1)).Return(&models.User{}, nil)
		_, err := svc.GetUserByID(1)
		assert.NoError(t, err)
	})

	t.Run("GetAllUsers", func(t *testing.T) {
		mockRepo.On("GetAllUsers").Return([]models.User{{}, {}}, nil)
		users, _ := svc.GetAllUsers()
		assert.Len(t, users, 2)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		mockRepo.On("DeleteUser", "1").Return(nil)
		err := svc.DeleteUser("1")
		assert.NoError(t, err)
	})
}

// Тесты услуг (Services)
func TestServiceMethods(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewSalonService(mockRepo)

	t.Run("AddService", func(t *testing.T) {
		mockRepo.On("CreateService", mock.Anything).Return(nil)
		err := svc.AddService(&models.Service{})
		assert.NoError(t, err)
	})

	t.Run("GetServices", func(t *testing.T) {
		mockRepo.On("GetAllServices").Return([]models.Service{{Title: "S1"}}, nil)
		res, _ := svc.GetServices()
		assert.Equal(t, "S1", res[0].Title)
	})

	t.Run("GetServiceByID", func(t *testing.T) {
		mockRepo.On("GetServiceByID", "1").Return(&models.Service{Title: "S1"}, nil)
		res, _ := svc.GetService("1")
		assert.Equal(t, "S1", res.Title)
	})
}

// Тесты сотрудников (Staff)
func TestStaffMethods(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewSalonService(mockRepo)

	t.Run("AddStaff", func(t *testing.T) {
		mockRepo.On("CreateStaff", mock.Anything).Return(nil)
		err := svc.AddStaff(&models.Staff{})
		assert.NoError(t, err)
	})

	t.Run("GetStaffList", func(t *testing.T) {
		mockRepo.On("GetAllStaff").Return([]models.Staff{{FullName: "Anna"}}, nil)
		res, _ := svc.GetStaffList()
		assert.Equal(t, "Anna", res[0].FullName)
	})
}

// Тесты бронирований (Bookings)
func TestBookingMethods(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewSalonService(mockRepo)

	t.Run("CreateBooking", func(t *testing.T) {
		mockRepo.On("CreateBooking", mock.Anything).Return(nil)
		err := svc.CreateBooking(&models.Booking{})
		assert.NoError(t, err)
	})

	t.Run("UpdateBooking - Success", func(t *testing.T) {
		booking := &models.Booking{Status: "pending"}
		booking.ID = 1
		updates := map[string]interface{}{"status": "confirmed"}

		// Последовательность вызовов в методе UpdateBooking:
		// 1. GetBookingByID
		// 2. UpdateBooking
		// 3. GetBookingByID (снова, чтобы вернуть обновленный)
		mockRepo.On("GetBookingByID", "1").Return(booking, nil).Twice()
		mockRepo.On("UpdateBooking", booking, updates).Return(nil).Once()

		res, err := svc.UpdateBooking("1", updates)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("CancelBooking", func(t *testing.T) {
		mockRepo.On("DeleteBooking", "1").Return(nil)
		err := svc.CancelBooking("1")
		assert.NoError(t, err)
	})
}

func TestDeleteService(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewSalonService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("DeleteService", "1").Return(nil).Once()

		err := svc.DeleteService("1")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.On("DeleteService", "99").Return(errors.New("db error")).Once()

		err := svc.DeleteService("99")

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestGetStaff(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewSalonService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		expectedStaff := &models.Staff{FullName: "Anna"}
		mockRepo.On("GetStaffByID", "1").Return(expectedStaff, nil).Once()

		res, err := svc.GetStaff("1")

		assert.NoError(t, err)
		assert.Equal(t, "Anna", res.FullName)
		assert.Equal(t, expectedStaff, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.On("GetStaffByID", "99").Return(nil, errors.New("staff not found")).Once()

		res, err := svc.GetStaff("99")

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "staff not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteStaff(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewSalonService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("DeleteStaff", "1").Return(nil).Once()

		err := svc.DeleteStaff("1")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo.On("DeleteStaff", "99").Return(errors.New("db error")).Once()

		err := svc.DeleteStaff("99")

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestGetBooking(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewSalonService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		expectedBooking := &models.Booking{Status: "confirmed"}
		mockRepo.On("GetBookingByID", "1").Return(expectedBooking, nil).Once()

		res, err := svc.GetBooking("1")

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "confirmed", res.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.On("GetBookingByID", "99").Return(nil, errors.New("not found")).Once()

		res, err := svc.GetBooking("99")

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestGetBookings(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewSalonService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		expectedBookings := []models.Booking{
			{Status: "pending"},
			{Status: "confirmed"},
		}

		mockRepo.On("GetAllBookings").Return(expectedBookings, nil).Once()

		res, err := svc.GetBookings()

		assert.NoError(t, err)
		assert.Len(t, res, 2)
		assert.Equal(t, "pending", res[0].Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.On("GetAllBookings").Return(nil, errors.New("db error")).Once()

		res, err := svc.GetBookings()

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "db error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateBooking(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewSalonService(mockRepo)

	id := "123"
	updates := map[string]interface{}{"status": "confirmed"}

	t.Run("Error on initial Fetch", func(t *testing.T) {
		mockRepo.On("GetBookingByID", id).Return(nil, errors.New("booking not found")).Once()

		res, err := svc.UpdateBooking(id, updates)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "booking not found", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error on Update execution", func(t *testing.T) {
		existingBooking := &models.Booking{Status: "pending"}
		existingBooking.ID = 123

		mockRepo.On("GetBookingByID", id).Return(existingBooking, nil).Once()
		mockRepo.On("UpdateBooking", existingBooking, updates).Return(errors.New("db write error")).Once()

		res, err := svc.UpdateBooking(id, updates)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "db write error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
