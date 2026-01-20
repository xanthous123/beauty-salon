package repository

import (
	"beauty-salon/internal/models"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type RepositorySuite struct {
	suite.Suite
	mock  sqlmock.Sqlmock
	db    *gorm.DB
	repo  Repository
	sqlDB *sql.DB
}

func (s *RepositorySuite) SetupTest() {
	var err error
	s.sqlDB, s.mock, err = sqlmock.New()
	assert.NoError(s.T(), err)

	dialector := postgres.New(postgres.Config{
		Conn: s.sqlDB,
	})

	s.db, err = gorm.Open(dialector, &gorm.Config{})
	assert.NoError(s.T(), err)

	s.repo = NewPostgresRepository(s.db)
}

func (s *RepositorySuite) TearDownTest() {
	s.sqlDB.Close()
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

func (s *RepositorySuite) TestCreateUser() {
	user := &models.User{Username: "test"}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	err := s.repo.CreateUser(user)
	assert.NoError(s.T(), err)
}

func (s *RepositorySuite) TestGetUserByUsername() {
	username := "admin"
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 AND "users"."deleted_at" IS NULL`)).
		WithArgs(username, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, username))

	res, err := s.repo.GetUserByUsername(username)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), username, res.Username)
}

func (s *RepositorySuite) TestGetUserByID() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL`)).
		WithArgs(uint(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	res, err := s.repo.GetUserByID(1)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), uint(1), res.ID)
}

func (s *RepositorySuite) TestGetAllUsers() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "u1").AddRow(2, "u2"))

	res, err := s.repo.GetAllUsers()
	assert.NoError(s.T(), err)
	assert.Len(s.T(), res, 2)
}

func (s *RepositorySuite) TestDeleteUser() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "deleted_at"=`)).
		WithArgs(sqlmock.AnyArg(), "1").
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	err := s.repo.DeleteUser("1")
	assert.NoError(s.T(), err)
}

func (s *RepositorySuite) TestCreateService() {
	srv := &models.Service{Title: "Haircut"}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "services"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	err := s.repo.CreateService(srv)
	assert.NoError(s.T(), err)
}

func (s *RepositorySuite) TestGetAllServices() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "services" WHERE "services"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "S1"))

	res, err := s.repo.GetAllServices()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "S1", res[0].Title)
}

func (s *RepositorySuite) TestGetServiceByID() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "services" WHERE id = $1 AND "services"."deleted_at" IS NULL`)).
		WithArgs("1", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "S1"))

	res, err := s.repo.GetServiceByID("1")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "S1", res.Title)
}

func (s *RepositorySuite) TestDeleteService() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "services" SET "deleted_at"=`)).
		WithArgs(sqlmock.AnyArg(), "1").
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	err := s.repo.DeleteService("1")
	assert.NoError(s.T(), err)
}

func (s *RepositorySuite) TestCreateStaff() {
	st := &models.Staff{FullName: "Anna"}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "staffs"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	err := s.repo.CreateStaff(st)
	assert.NoError(s.T(), err)
}

func (s *RepositorySuite) TestGetAllStaff() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "staffs" WHERE "staffs"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "full_name"}).AddRow(1, "Anna"))

	res, err := s.repo.GetAllStaff()
	assert.NoError(s.T(), err)
	assert.Len(s.T(), res, 1)
}

func (s *RepositorySuite) TestGetStaffByID() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "staffs" WHERE id = $1 AND "staffs"."deleted_at" IS NULL`)).
		WithArgs("1", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "full_name"}).AddRow(1, "Anna"))

	res, err := s.repo.GetStaffByID("1")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Anna", res.FullName)
}

func (s *RepositorySuite) TestDeleteStaff() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "staffs" SET "deleted_at"=`)).
		WithArgs(sqlmock.AnyArg(), "1").
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	err := s.repo.DeleteStaff("1")
	assert.NoError(s.T(), err)
}

func (s *RepositorySuite) TestGetAllBookings() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bookings" WHERE "bookings"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "service_id", "staff_id"}).
			AddRow(1, 1, 1, 1))

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "services" WHERE "services"."id" = $1 AND "services"."deleted_at" IS NULL`)).
		WithArgs(uint(1)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "staffs" WHERE "staffs"."id" = $1 AND "staffs"."deleted_at" IS NULL`)).
		WithArgs(uint(1)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL`)).
		WithArgs(uint(1)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	res, err := s.repo.GetAllBookings()
	assert.NoError(s.T(), err)
	assert.Len(s.T(), res, 1)
}

func (s *RepositorySuite) TestGetBookingByID() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bookings" WHERE id = $1 AND "bookings"."deleted_at" IS NULL`)).
		WithArgs("1", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "service_id", "staff_id"}).AddRow(1, 1, 1, 1))

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "services" WHERE "services"."id" = $1 AND "services"."deleted_at" IS NULL`)).
		WithArgs(uint(1)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "staffs" WHERE "staffs"."id" = $1 AND "staffs"."deleted_at" IS NULL`)).
		WithArgs(uint(1)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL`)).
		WithArgs(uint(1)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	res, err := s.repo.GetBookingByID("1")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), res)
}

func (s *RepositorySuite) TestUpdateBooking() {
	b := &models.Booking{Model: gorm.Model{ID: 1}}
	updates := map[string]interface{}{"status": "confirmed"}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "bookings" SET "status"=$1,"updated_at"=$2 WHERE "bookings"."deleted_at" IS NULL AND "id" = $3`)).
		WithArgs("confirmed", sqlmock.AnyArg(), uint(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	err := s.repo.UpdateBooking(b, updates)
	assert.NoError(s.T(), err)
}

func (s *RepositorySuite) TestCreateBooking() {
	booking := &models.Booking{
		UserID:    1,
		ServiceID: 1,
		StaffID:   1,
		Date:      "2026-01-20",
		Status:    "pending",
	}

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "bookings"`)).
		WithArgs(
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			nil,              // deleted_at
			booking.UserID,
			booking.ServiceID,
			booking.StaffID,
			booking.Date, // поле Date
			booking.Status,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectCommit()

	err := s.repo.CreateBooking(booking)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), uint(1), booking.ID)
}

func (s *RepositorySuite) TestCreateBooking_Error() {
	booking := &models.Booking{UserID: 1}

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "bookings"`)).
		WillReturnError(errors.New("db connection lost"))

	s.mock.ExpectRollback()

	err := s.repo.CreateBooking(booking)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), "db connection lost", err.Error())
}

func (s *RepositorySuite) TestDeleteBooking() {
	bookingID := "1"

	s.mock.ExpectBegin()

	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "bookings" SET "deleted_at"=`)).
		WithArgs(
			sqlmock.AnyArg(),
			bookingID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	s.mock.ExpectCommit()

	err := s.repo.DeleteBooking(bookingID)

	assert.NoError(s.T(), err)
}

func (s *RepositorySuite) TestDeleteBooking_Error() {
	bookingID := "99"

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "bookings" SET "deleted_at"=`)).
		WithArgs(sqlmock.AnyArg(), bookingID).
		WillReturnError(errors.New("db error on delete"))

	s.mock.ExpectRollback()

	err := s.repo.DeleteBooking(bookingID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), "db error on delete", err.Error())
}
