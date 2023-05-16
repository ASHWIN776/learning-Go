package repository

import "github.com/ASHWIN776/learning-Go/internal/models"

type DatabaseRepo interface {
	InsertReservation(res models.Reservation) error
}
