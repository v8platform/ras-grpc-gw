package domain

import "time"

type AccessToken struct {
	UserId   int32
	ClientId int32
	Expiries time.Time
}
