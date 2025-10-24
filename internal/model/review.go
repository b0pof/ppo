package model

import "time"

type Review struct {
	ID            int64
	UserID        int64
	ItemID        int64
	Username      string
	Rating        int
	Advantages    string
	Disadvantages string
	Comment       string
	CreatedAt     time.Time
}
