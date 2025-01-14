package queries

import (
	"codetrack/internal/db"

	"github.com/samber/do/v2"
)

func NewQueries(i do.Injector) (*Queries, error) {
	return &Queries{
		Database: do.MustInvoke[*db.Database](i),
	}, nil
}

type Queries struct {
	Database *db.Database
}

func (query *Queries) GetUser(email string) (*db.UsersTable, error) {
	user := db.UsersTable{
		Email: email,
	}

	result := query.Database.Database.Find(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &user, nil
}

func (query *Queries) UserExists(email string) (bool, error) {
	user, err := query.GetUser(email)

	if err != nil {
		return false, err
	}

	if user == nil {
		return false, nil
	}

	return true, nil
}

func (query *Queries) NewUser(email string, passwordHash string) (error) {
	user := db.UsersTable{
		Email: email,
		PasswordHash: passwordHash,
	}

	result := query.Database.Database.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}