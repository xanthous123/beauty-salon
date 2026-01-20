package main

import (
	"beauty-salon/internal/handlers"
	"beauty-salon/internal/middleware"
	"beauty-salon/internal/models"
	"beauty-salon/internal/repository"
	"beauty-salon/internal/service"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load()

	// DB
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.User{}, &models.Service{}, &models.Staff{}, &models.Booking{})

	// Redis
	rdb := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_HOST")})

	// Dependency Injection
	repo := repository.NewPostgresRepository(db)
	svc := service.NewSalonService(repo)
	h := handlers.NewHandler(svc)

	// Router
	r := gin.Default()
	r.Use(middleware.RateLimiter(rdb, 100, time.Minute)) // Анти-спам: 100 req/min

	api := r.Group("/api/v1")
	{
		api.POST("/register", h.Register)
		api.POST("/login", h.Login)

		auth := api.Group("/")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.POST("/logout", h.Logout)
			auth.GET("/users/me", h.GetMe)
			auth.GET("/users", h.GetAllUsers)
			auth.DELETE("/users/:id", h.DeleteUser)

			auth.POST("/services", h.AddService)
			auth.GET("/services", h.GetServices)
			auth.GET("/services/:id", h.GetServiceByID)
			auth.DELETE("/services/:id", h.DeleteService)

			auth.POST("/staff", h.AddStaff)
			auth.GET("/staff", h.GetStaff)
			auth.GET("/staff/:id", h.GetStaffByID)
			auth.DELETE("/staff/:id", h.DeleteStaff)

			auth.POST("/bookings", h.CreateBooking)
			auth.GET("/bookings", h.GetBookings)
			auth.GET("/bookings/:id", h.GetBookingByID)
			auth.PATCH("/bookings/:id", h.PatchBooking)
			auth.DELETE("/bookings/:id", h.DeleteBooking)
		}
	}
	r.Run(":" + os.Getenv("PORT"))
}
