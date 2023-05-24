package repository

import (
	"time"

	"github.com/ASHWIN776/learning-Go/internal/models"
)

type DatabaseRepo interface {
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestriction) error
	SearchAvailabilityByRoomId(startDate, endDate time.Time, roomId int) (bool, error)
	SearchAvailabilityForAllRooms(startDate, endDate time.Time) ([]models.Room, error)
	GetRoomById(id int) (models.Room, error)
	GetUserById(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, enteredPass string) (int, string, error)
	AllReservations() ([]models.Reservation, error)
	NewReservations() ([]models.Reservation, error)
	GetReservationById(id int) (models.Reservation, error)
	UpdateReservation(res models.Reservation) error
	DeleteReservation(id int) error
	UpdateProcessedForReservation(id, processed int) error
	AllRooms() ([]models.Room, error)
	GetRestrictionsForRoomsByDate(roomId int, startDate, endDate time.Time) ([]models.RoomRestriction, error)
}
