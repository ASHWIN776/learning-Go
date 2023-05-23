package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ASHWIN776/learning-Go/internal/models"
	"golang.org/x/crypto/bcrypt"
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
func (p *postgresDBRepo) SearchAvailabilityByRoomId(startDate, endDate time.Time, roomId int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var resCount int // stores the count of existing restricted rooms

	stmt := `
	select 
		count(id) 
	from 
		room_restrictions 
	where
		room_id = $1 AND 
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

// Returns a room type variable after searching one by Id
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

// Returns a user type variable after searching one by Id
func (p *postgresDBRepo) GetUserById(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	stmt := `
		select * 
		from 
			users
		where 
			id = $1
	`

	row := p.DB.QueryRowContext(ctx, stmt, id)

	var foundUser models.User
	err := row.Scan(
		&foundUser.ID,
		&foundUser.FirstName,
		&foundUser.LastName,
		&foundUser.Email,
		&foundUser.Password,
		&foundUser.CreatedAt,
		&foundUser.UpdatedAt,
		&foundUser.AccessLevel,
	)

	if err != nil {
		return foundUser, err
	}

	return foundUser, nil
}

// Updates a user in the database
func (p *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	stmt := `
		Update users
		set 
			first_name = $1
			last_name = $2
			email = $3
			access_level = $4
			updated_at = $5
		where 
			id = $6
	`

	_, err := p.DB.ExecContext(ctx, stmt, u.FirstName, u.LastName, u.Email, u.AccessLevel, time.Now(), u.ID)

	if err != nil {
		return err
	}

	return nil
}

// Will check the creds entered by the user and returns the user id, hashed password if the creds are correct
func (p *postgresDBRepo) Authenticate(enteredEmail, enteredPass string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	// Get the user from the database using the email(unique)
	var id int
	var hashedPass string

	stmt := `
		select id, password
		from users
		where 
			email = $1
	`

	row := p.DB.QueryRowContext(ctx, stmt, enteredEmail)
	err := row.Scan(&id, &hashedPass)
	if err != nil {
		return -1, "", err
	}

	// Compare the password hash in the db and the one given by the user
	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(enteredPass))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return -1, "", errors.New("incorrect password")
	} else if err != nil {
		return -1, "", err
	}

	// if it gets here, then password entered is correct
	return id, hashedPass, nil
}

// Returns a slice of all reservations
func (p *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	stmt := `
		select r.*, rm.* 
		from 
			reservations r 
			left join rooms rm
		on
			(r.room_id = rm.id)
		order by
			r.start_date asc
	`

	rows, err := p.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var reservations []models.Reservation

	for rows.Next() {
		var res models.Reservation
		err := rows.Scan(
			&res.ID, // reservation info
			&res.FirstName,
			&res.LastName,
			&res.Email,
			&res.Phone,
			&res.StartDate,
			&res.EndDate,
			&res.RoomID,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.Room.ID, // start getting the room info
			&res.Room.RoomName,
			&res.Room.CreatedAt,
			&res.Room.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		reservations = append(reservations, res)
	}

	// Catching any errors during iteration
	if rows.Err() != nil {
		return nil, err
	}

	return reservations, nil
}

// Returns the list of reservations which are not processed(processed = 0)
func (p *postgresDBRepo) NewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	stmt := `
		select r.*, rm.* 
		from 
			reservations r 
			left join rooms rm
		on
			(r.room_id = rm.id)
		where 
			r.processed = 0
		order by
			r.start_date asc
	`

	rows, err := p.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var reservations []models.Reservation

	for rows.Next() {
		var res models.Reservation
		err := rows.Scan(
			&res.ID, // reservation info
			&res.FirstName,
			&res.LastName,
			&res.Email,
			&res.Phone,
			&res.StartDate,
			&res.EndDate,
			&res.RoomID,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.Processed,
			&res.Room.ID, // start getting the room info
			&res.Room.RoomName,
			&res.Room.CreatedAt,
			&res.Room.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		reservations = append(reservations, res)
	}

	// Catching any errors during iteration
	if rows.Err() != nil {
		return nil, err
	}

	return reservations, nil
}
