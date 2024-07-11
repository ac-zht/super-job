package domain

type User struct {
	Id      int64
	Name    string
	Email   string
	IsAdmin uint8
	Status  uint8
}
