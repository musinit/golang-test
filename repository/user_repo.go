package repository

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	usersGetAllQuery = "SELECT * FROM users"
	userCreateQuery  = "INSERT INTO USERS VALUES ($1, $2, $3, $4)"
)

type UserRepository struct {
	Db *pgxpool.Pool
}

func (this *UserRepository) GetAllUsers() ([]*User, error) {
	var users []*User
	ctx := context.Background()
	var err error
	err = pgxscan.Get(ctx, this.Db, &users, usersGetAllQuery)
	if err != nil && err.Error() != "no rows in result set" {
		return nil, err
	}
	return users, nil
}

func (this *UserRepository) Create(user *User) error {
	ctx := context.Background()
	var err error
	_, err = this.Db.Exec(ctx, userCreateQuery,
		user.ID, user.Name, user.Surname, user.Patronymic,
	)
	if err != nil {
		return err
	}

	return nil
}

func (this *UserRepository) Close(user *User) {
	this.Db.Close()
}
