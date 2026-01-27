package handlers

import (
	"beauty-salon/internal/models"
	"beauty-salon/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc service.Service
}

func NewHandler(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(c *gin.Context) {
	var i struct{ Username, Password string }
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Register(i.Username, i.Password); err != nil {
		c.JSON(500, gin.H{"error": "Failed"})
		return
	}
	c.JSON(201, gin.H{"message": "Registered"})
}

func (h *Handler) Login(c *gin.Context) {
	var i struct{ Username, Password string }
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	t, err := h.svc.Login(i.Username, i.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"token": t})
}

func (h *Handler) Logout(c *gin.Context) { c.JSON(200, gin.H{"message": "Logged out"}) }

// Users
func (h *Handler) GetMe(c *gin.Context) {
	u, err := h.svc.GetUserByID(c.MustGet("userID").(uint))
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	c.JSON(200, u)
}

func (h *Handler) GetAllUsers(c *gin.Context) {
	u, _ := h.svc.GetAllUsers()
	c.JSON(200, u)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	if err := h.svc.DeleteUser(c.Param("id")); err != nil {
		c.JSON(500, gin.H{"error": "Failed"})
		return
	}
	c.Status(204)
}

// Services
func (h *Handler) AddService(c *gin.Context) {
	var s models.Service
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.AddService(&s); err != nil {
		c.JSON(500, gin.H{"error": "Failed"})
		return
	}
	c.JSON(201, s)
}

func (h *Handler) GetServices(c *gin.Context) {
	s, _ := h.svc.GetServices()
	c.JSON(200, s)
}

func (h *Handler) GetServiceByID(c *gin.Context) {
	s, err := h.svc.GetService(c.Param("id"))
	if err != nil {
		c.JSON(404, gin.H{"error": "Service not found"})
		return
	}
	c.JSON(200, s)
}

func (h *Handler) DeleteService(c *gin.Context) {
	h.svc.DeleteService(c.Param("id"))
	c.Status(204)
}

// Staff
func (h *Handler) AddStaff(c *gin.Context) {
	var s models.Staff
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	h.svc.AddStaff(&s)
	c.JSON(201, s)
}

func (h *Handler) GetStaff(c *gin.Context) {
	s, _ := h.svc.GetStaffList()
	c.JSON(200, s)
}

func (h *Handler) GetStaffByID(c *gin.Context) {
	s, err := h.svc.GetStaff(c.Param("id"))
	if err != nil {
		c.JSON(404, gin.H{"error": "Staff not found"})
		return
	}
	c.JSON(200, s)
}

func (h *Handler) DeleteStaff(c *gin.Context) {
	h.svc.DeleteStaff(c.Param("id"))
	c.Status(204)
}

// Bookings
func (h *Handler) CreateBooking(c *gin.Context) {
	var b models.Booking
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	b.UserID = c.MustGet("userID").(uint)
	if err := h.svc.CreateBooking(&b); err != nil {
		c.JSON(500, gin.H{"error": "Failed"})
		return
	}
	c.JSON(201, b)
}

func (h *Handler) GetBookings(c *gin.Context) {
	b, _ := h.svc.GetBookings()
	c.JSON(200, b)
}

func (h *Handler) GetBookingByID(c *gin.Context) {
	b, err := h.svc.GetBooking(c.Param("id"))
	if err != nil {
		c.JSON(404, gin.H{"error": "Booking not found"})
		return
	}
	c.JSON(200, b)
}

func (h *Handler) PatchBooking(c *gin.Context) {
	var u map[string]interface{}
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	b, err := h.svc.UpdateBooking(c.Param("id"), u)
	if err != nil {
		c.JSON(500, gin.H{"error": "Update failed"})
		return
	}
	c.JSON(200, b)
}

func (h *Handler) DeleteBooking(c *gin.Context) {
	h.svc.CancelBooking(c.Param("id"))
	c.Status(204)
}
