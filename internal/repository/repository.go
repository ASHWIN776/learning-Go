package repository

import "github.com/ASHWIN776/learning-Go/internal/models"

type DatabaseRepo interface {
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestriction) error
	SearchAvailabilityByDates(startDate, endDate string, roomId int) (bool, error)
}
