package dbrepo

import (
	"context"
	"time"

	"github.com/ASHWIN776/learning-Go/internal/models"
)

// To insert the form data from make-reservation form to the reservations table in the database
func (p *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	// This query will insert a reservation and also return the corresponding id
	stmt := `insert into reservations 
	(first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) 
	values 
	($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	var reservationId int

	row := p.DB.QueryRowContext(ctx, stmt, res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, res.RoomID, time.Now(), time.Now())

	err := row.Scan(&reservationId)

	if err != nil {
		return -1, err
	}

	return reservationId, nil
}

// Inserts a room restriction into the corresponding table in the database
func (p *postgresDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	stmt := `insert into room_restrictions 
	(start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at) 
	values 
	($1, $2, $3, $4, $5, $6, $7)`

	_, err := p.DB.ExecContext(ctx, stmt, res.StartDate, res.EndDate, res.RoomID, res.ReservationID, res.RestrictionID, time.Now(), time.Now())

	if err != nil {
		return err
	}

	return nil
}
