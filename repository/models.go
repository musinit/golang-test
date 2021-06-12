package repository

import "github.com/google/uuid"

type User struct {
	ID         uuid.UUID `db:"id"`
	Name       string    `db:"name"`
	Surname    string    `db:"surname"`
	Patronymic string    `db:"patronymic"`
}
