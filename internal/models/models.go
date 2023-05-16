package models

import "time"

type Reservation struct {
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
}

// Users model
type Users struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Rooms model
type Rooms struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restrictions model
type Restrictions struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// RoomRestrictions model
type RoomRestrictions struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	RestrictionID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Rooms        // There is no necessity that I have to put the same amt of vars in these models as I do in the table
	Reservation   Reservations // There is no necessity that I have to put the same amt of vars in these models as I do in the table
	Restriction   Restrictions // There is no necessity that I have to put the same amt of vars in these models as I do in the table
}

// Reservations model
type Reservations struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Rooms // There is no necessity that I have to put the same amt of vars in these models as I do in the table
}
