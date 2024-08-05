package domain

import "time"

const TokenDuration = time.Hour * 4

type User struct {
	Id       int64
	Name     string
	Email    string
	Password string
	IsAdmin  uint8
	Status   uint8
	Salt     string
	Token    string
}

type LoginLog struct {
	Id       int64
	Username string
	Ip       string
	Ctime    time.Time
}
