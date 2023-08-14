package common

import "time"

type EventID int64

type Event struct {
	ID          EventID   `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DateStart   time.Time `json:"dateStart"`
	DateEnd     time.Time `json:"dateEnd"`
	NotifyIn    int64     `json:"notifyIn"`
	OwnerID     int64     `json:"ownerId"`
}

type Notice struct {
	EventID  EventID   `json:"eventId"`
	Title    string    `json:"title"`
	Datetime time.Time `json:"datetime"`
	OwnerID  int64     `json:"ownerId"`
}

type Exception struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseEventID struct {
	Code    int     `json:"code"`
	EventID EventID `json:"eventId"`
}

type ResponseEventList struct {
	Code   int     `json:"code"`
	Events []Event `json:"events"`
}
