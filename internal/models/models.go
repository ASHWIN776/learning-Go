package models

import (
	"html/template"
	"time"
)

// User model
type User struct {
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
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restriction model
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// RoomRestriction model
type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	RestrictionID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Room        // There is no necessity that I have to put the same amt of vars in these models as I do in the table
	Reservation   Reservation // There is no necessity that I have to put the same amt of vars in these models as I do in the table
	Restriction   Restriction // There is no necessity that I have to put the same amt of vars in these models as I do in the table
}

// Reservations model
type Reservation struct {
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
	Room      Room // There is no necessity that I have to put the same amt of vars in these models as I do in the table
}

// holds an email message
type MailData struct {
	To      string
	From    string
	Subject string
	Content template.HTML
}
