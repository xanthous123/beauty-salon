package handlers

import (
	"beauty-salon/internal/models"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Register(u, p string) error {
	return m.Called(u, p).Error(0)
}

func (m *MockService) Login(u, p string) (string, error) {
	args := m.Called(u, p)
	return args.String(0), args.Error(1)
}

func (m *MockService) GetUserByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockService) GetAllUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}
func (m *MockService) DeleteUser(id string) error { return m.Called(id).Error(0) }

func (m *MockService) AddService(s *models.Service) error { return m.Called(s).Error(0) }

func (m *MockService) GetServices() ([]models.Service, error) {
	args := m.Called()
	return args.Get(0).([]models.Service), args.Error(1)
}

func (m *MockService) GetService(id string) (*models.Service, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Service), args.Error(1)
}

func (m *MockService) DeleteService(id string) error { return m.Called(id).Error(0) }

func (m *MockService) AddStaff(s *models.Staff) error { return m.Called(s).Error(0) }

func (m *MockService) GetStaffList() ([]models.Staff, error) {
	args := m.Called()
	return args.Get(0).([]models.Staff), args.Error(1)
}

func (m *MockService) GetStaff(id string) (*models.Staff, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Staff), args.Error(1)
}

func (m *MockService) DeleteStaff(id string) error { return m.Called(id).Error(0) }

func (m *MockService) CreateBooking(b *models.Booking) error { return m.Called(b).Error(0) }

func (m *MockService) GetBookings() ([]models.Booking, error) {
	args := m.Called()
	return args.Get(0).([]models.Booking), args.Error(1)
}

func (m *MockService) GetBooking(id string) (*models.Booking, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Booking), args.Error(1)
}

func (m *MockService) UpdateBooking(id string, u map[string]interface{}) (*models.Booking, error) {
	args := m.Called(id, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Booking), args.Error(1)
}

func (m *MockService) CancelBooking(id string) error { return m.Called(id).Error(0) }

func setup() (*gin.Engine, *MockService, *Handler) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockService)
	h := NewHandler(mockSvc)
	r := gin.New()
	return r, mockSvc, h
}

func TestRegister(t *testing.T) {
	r, mockSvc, h := setup()
	r.POST("/register", h.Register)

	t.Run("Success", func(t *testing.T) {
		mockSvc.On("Register", "user", "pass").Return(nil).Once()
		body, _ := json.Marshal(map[string]string{"username": "user", "password": "pass"})
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, 201, w.Code)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte("{invalid}")))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc.On("Register", "error_user", "pass").Return(errors.New("db fail")).Once()

		body, _ := json.Marshal(map[string]string{"username": "error_user", "password": "pass"})
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, 500, w.Code)
		assert.Contains(t, w.Body.String(), "Failed")
	})
}

func TestLogin(t *testing.T) {
	r, mockSvc, h := setup()
	r.POST("/login", h.Login)

	t.Run("Success", func(t *testing.T) {
		mockSvc.On("Login", "user", "pass").Return("token123", nil).Once()
		body, _ := json.Marshal(map[string]string{"username": "user", "password": "pass"})
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "token123")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockSvc.On("Login", "wrong", "wrong").Return("", errors.New("unauthorized")).Once()
		body, _ := json.Marshal(map[string]string{"username": "wrong", "password": "wrong"})
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, 401, w.Code)
	})

	t.Run("Invalid JSON Syntax", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(`{"username": "user", "password": `)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid input")
	})

	t.Run("Empty Body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(``)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid input")
	})
}

func TestLogout(t *testing.T) {
	r, _, h := setup()
	r.POST("/logout", h.Logout)
	req, _ := http.NewRequest("POST", "/logout", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestGetMe(t *testing.T) {
	r, mockSvc, h := setup()

	r.GET("/me", func(c *gin.Context) {
		userID := uint(1)
		if c.GetHeader("X-Mock-Fail") == "true" {
			userID = uint(99)
		}
		c.Set("userID", userID)
		h.GetMe(c)
	})

	t.Run("Success", func(t *testing.T) {
		mockSvc.On("GetUserByID", uint(1)).Return(&models.User{Username: "tester"}, nil).Once()

		req, _ := http.NewRequest("GET", "/me", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "tester")
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockSvc.On("GetUserByID", uint(99)).Return(nil, errors.New("user not found")).Once()

		req, _ := http.NewRequest("GET", "/me", nil)
		req.Header.Set("X-Mock-Fail", "true") // Триггерим смену ID на 99
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
		assert.Contains(t, w.Body.String(), "User not found")
	})
}

func TestGetAllUsers(t *testing.T) {
	r, mockSvc, h := setup()
	r.GET("/users", h.GetAllUsers)

	mockSvc.On("GetAllUsers").Return([]models.User{{Username: "u1"}, {Username: "u2"}}, nil)
	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestDeleteUser(t *testing.T) {
	r, mockSvc, h := setup()
	r.DELETE("/users/:id", h.DeleteUser)

	t.Run("Success", func(t *testing.T) {
		mockSvc.On("DeleteUser", "1").Return(nil).Once()

		req, _ := http.NewRequest("DELETE", "/users/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 204, w.Code)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockSvc.On("DeleteUser", "99").Return(errors.New("database connection failed")).Once()

		req, _ := http.NewRequest("DELETE", "/users/99", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 500, w.Code)
		assert.Contains(t, w.Body.String(), "Failed")
	})
}

func TestAddService(t *testing.T) {
	r, mockSvc, h := setup()
	r.POST("/services", h.AddService)

	t.Run("Success", func(t *testing.T) {
		mockSvc.On("AddService", mock.Anything).Return(nil).Once()
		body, _ := json.Marshal(models.Service{Title: "Haircut"})
		req, _ := http.NewRequest("POST", "/services", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		body := []byte(`{"title": "Haircut`)
		req, _ := http.NewRequest("POST", "/services", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("Service Failure", func(t *testing.T) {
		mockSvc.On("AddService", mock.Anything).Return(errors.New("db error")).Once()

		body, _ := json.Marshal(models.Service{Title: "Massage"})
		req, _ := http.NewRequest("POST", "/services", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 500, w.Code)
		assert.Contains(t, w.Body.String(), "Failed")
	})
}

func TestGetServices(t *testing.T) {
	r, mockSvc, h := setup()
	r.GET("/services", h.GetServices)

	mockSvc.On("GetServices").Return([]models.Service{{Title: "S1"}}, nil)
	req, _ := http.NewRequest("GET", "/services", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestGetServiceByID(t *testing.T) {
	r, mockSvc, h := setup()
	r.GET("/services/:id", h.GetServiceByID)

	t.Run("Found", func(t *testing.T) {
		mockSvc.On("GetService", "1").Return(&models.Service{Title: "S1"}, nil).Once()
		req, _ := http.NewRequest("GET", "/services/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockSvc.On("GetService", "99").Return(nil, errors.New("not found")).Once()
		req, _ := http.NewRequest("GET", "/services/99", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, 404, w.Code)
	})
}

func TestDeleteService(t *testing.T) {
	r, mockSvc, h := setup()
	r.DELETE("/services/:id", h.DeleteService)
	mockSvc.On("DeleteService", "1").Return(nil)
	req, _ := http.NewRequest("DELETE", "/services/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 204, w.Code)
}

func TestAddStaff(t *testing.T) {
	r, mockSvc, h := setup()
	r.POST("/staff", h.AddStaff)

	t.Run("Success", func(t *testing.T) {
		mockSvc.On("AddStaff", mock.Anything).Return(nil).Once()
		body, _ := json.Marshal(models.Staff{FullName: "Master"})
		req, _ := http.NewRequest("POST", "/staff", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)
		assert.Contains(t, w.Body.String(), "Master")
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		body := []byte(`{ "fullName": "Incomplete JSON ...`)
		req, _ := http.NewRequest("POST", "/staff", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})
}

func TestGetStaff(t *testing.T) {
	r, mockSvc, h := setup()
	r.GET("/staff", h.GetStaff)
	mockSvc.On("GetStaffList").Return([]models.Staff{}, nil)
	req, _ := http.NewRequest("GET", "/staff", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestGetStaffByID(t *testing.T) {
	r, mockSvc, h := setup()
	r.GET("/staff/:id", h.GetStaffByID)

	t.Run("Success", func(t *testing.T) {
		mockSvc.On("GetStaff", "1").Return(&models.Staff{FullName: "Anna"}, nil).Once()

		req, _ := http.NewRequest("GET", "/staff/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "Anna")
	})

	t.Run("Not Found", func(t *testing.T) {
		mockSvc.On("GetStaff", "99").Return(nil, errors.New("not found")).Once()

		req, _ := http.NewRequest("GET", "/staff/99", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
		assert.Contains(t, w.Body.String(), "Staff not found")
	})
}

func TestDeleteStaff(t *testing.T) {
	r, mockSvc, h := setup()
	r.DELETE("/staff/:id", h.DeleteStaff)
	mockSvc.On("DeleteStaff", "1").Return(nil)
	req, _ := http.NewRequest("DELETE", "/staff/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 204, w.Code)
}

func TestCreateBooking(t *testing.T) {
	r, mockSvc, h := setup()
	r.POST("/bookings", func(c *gin.Context) {
		c.Set("userID", uint(1)) // Имитация авторизации
		h.CreateBooking(c)
	})

	t.Run("Success", func(t *testing.T) {
		mockSvc.On("CreateBooking", mock.Anything).Return(nil).Once()

		body, _ := json.Marshal(models.Booking{ServiceID: 1, StaffID: 1})
		req, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		body := []byte(`{"service_id": 1, "staff_id": }`)
		req, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("Service Failure", func(t *testing.T) {
		mockSvc.On("CreateBooking", mock.Anything).Return(errors.New("db error")).Once()

		body, _ := json.Marshal(models.Booking{ServiceID: 1, StaffID: 1})
		req, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 500, w.Code)
		assert.Contains(t, w.Body.String(), "Failed")
	})
}

func TestGetBookings(t *testing.T) {
	r, mockSvc, h := setup()
	r.GET("/bookings", h.GetBookings)
	mockSvc.On("GetBookings").Return([]models.Booking{}, nil)
	req, _ := http.NewRequest("GET", "/bookings", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestGetBookingByID(t *testing.T) {
	r, mockSvc, h := setup()
	r.GET("/bookings/:id", h.GetBookingByID)

	t.Run("Success", func(t *testing.T) {
		mockSvc.On("GetBooking", "1").Return(&models.Booking{Status: "pending"}, nil).Once()

		req, _ := http.NewRequest("GET", "/bookings/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "pending")
	})

	t.Run("Not Found", func(t *testing.T) {
		mockSvc.On("GetBooking", "99").Return(nil, errors.New("booking not found")).Once()

		req, _ := http.NewRequest("GET", "/bookings/99", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
		assert.Contains(t, w.Body.String(), "Booking not found")
	})
}

func TestPatchBooking(t *testing.T) {
	r, mockSvc, h := setup()
	r.PATCH("/bookings/:id", h.PatchBooking)

	t.Run("Success", func(t *testing.T) {
		mockSvc.On("UpdateBooking", "1", mock.MatchedBy(func(u map[string]interface{}) bool {
			return u["status"] == "confirmed"
		})).Return(&models.Booking{Status: "confirmed"}, nil).Once()

		body, _ := json.Marshal(map[string]interface{}{"status": "confirmed"})
		req, _ := http.NewRequest("PATCH", "/bookings/1", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "confirmed")
	})

	t.Run("Invalid JSON (400)", func(t *testing.T) {
		req, _ := http.NewRequest("PATCH", "/bookings/1", bytes.NewBuffer([]byte(`{"status": "no-quotes-here`)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid input")
	})

	t.Run("Update Failed (500)", func(t *testing.T) {
		mockSvc.On("UpdateBooking", "99", mock.Anything).Return(nil, errors.New("db error")).Once()

		body, _ := json.Marshal(map[string]interface{}{"status": "cancelled"})
		req, _ := http.NewRequest("PATCH", "/bookings/99", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 500, w.Code)
		assert.Contains(t, w.Body.String(), "Update failed")
	})
}

func TestDeleteBooking(t *testing.T) {
	r, mockSvc, h := setup()
	r.DELETE("/bookings/:id", h.DeleteBooking)
	mockSvc.On("CancelBooking", "1").Return(nil)
	req, _ := http.NewRequest("DELETE", "/bookings/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 204, w.Code)
}
