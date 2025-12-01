package models

import "time"

type Room struct {
	ID        string    `bson:"id" json:"id"`
	Pin       string    `bson:"pin" json:"pin"`
	DJID      string    `bson:"dj_id" json:"dj_id"`
	DJEmail   string    `bson:"dj_email" json:"dj_email"`
	Active    bool      `bson:"active" json:"active"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

