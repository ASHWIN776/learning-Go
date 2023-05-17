package dbrepo

import (
	"context"
	"log"
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

// Returns a bool specifying if there is an availability(of the specified room) or not, and a potential error
func (p *postgresDBRepo) SearchAvailabilityByRoomId(startDate, endDate string, roomId int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var resCount int // stores the count of existing restricted rooms

	stmt := `
	select 
		count(id) 
	from 
		room_restrictions 
	where
		roomId = $1 AND 
		$2 < end_date AND $3 > start_date
	`

	row := p.DB.QueryRowContext(ctx, stmt, roomId, startDate, endDate)

	err := row.Scan(&resCount)
	if err != nil {
		return false, err
	}

	// If satisfied, means there is no booking overlapping the specified start and end dates
	if resCount == 0 {
		return true, nil
	}

	return false, nil
}

// Returns a slice of all rooms available to book for the startDate and endDate range, and potentially an error if any
func (p *postgresDBRepo) SearchAvailabilityForAllRooms(startDate, endDate time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	stmt := `
	select 
		*
	from 
		rooms r
	where 
		r.id not in 
		(select room_id from room_restrictions where $1 < end_date and $2 > start_date);
	`

	rows, err := p.DB.QueryContext(ctx, stmt, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var allRooms []models.Room
	var id int
	var room_name string
	var created_at, updated_at time.Time

	for rows.Next() {
		err := rows.Scan(&id, &room_name, &created_at, &updated_at)
		if err != nil {
			return nil, err
		}

		allRooms = append(allRooms, models.Room{
			ID:        id,
			RoomName:  room_name,
			CreatedAt: created_at,
			UpdatedAt: updated_at,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	log.Println(allRooms)
	return allRooms, nil
}

func (p *postgresDBRepo) GetRoomById(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	stmt := `
	select 
		*
	from rooms
	where
		id = $1 
	`

	row := p.DB.QueryRowContext(ctx, stmt, id)
	var room models.Room

	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}

	return room, nil
}
