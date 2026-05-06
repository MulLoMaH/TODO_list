package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/MulLoMaH/TODO_list.git/internal/core/errors"
)

type User struct {
	ID      int // id
	Version int // версия

	FullName    string  // полное имя
	PhoneNumber *string // номер телефона
}

func NewUser(
	id int,
	version int,
	fullName string,
	phoneNumber *string,
) User {
	return User{
		ID:      id,
		Version: version,

		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func NewUserUninitialized(
	fullName string,
	phoneNumber *string,
) User {
	return NewUser(
		UninitializedID,
		UninitializedVersion,

		fullName,
		phoneNumber,
	)
}

func (u *User) Validate() error {
	//проверка на количество символов в строке (русские буквы = 2 байта = 1 руна)
	fullNameLenght := len([]rune(u.FullName))

	if fullNameLenght < 3 || fullNameLenght > 100 {
		return fmt.Errorf(
			"invalid 'FullName' len: %d: %w",
			fullNameLenght,
			core_errors.ErrInvalidArgument,
		)
	}

	if u.PhoneNumber != nil {
		phoneNumberLen := len([]rune(*u.PhoneNumber))
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf(
				"invalid 'PhoneNumber' len: %d: %w",
				phoneNumberLen,
				core_errors.ErrInvalidArgument,
			)
		}

		re := regexp.MustCompile(`^\+[0-9]+$`)

		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf(
				"invalid 'PhoneNumber' format: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}
