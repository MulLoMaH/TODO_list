package domain

type User struct {
	ID      int // id
	Version int // версия

	FullName    string  // полное имя
	PhoneNumber *string // номер телефона
}
