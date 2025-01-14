package services

import (
	"codetrack/internal/queries"
	"errors"

	"github.com/samber/do/v2"
)

func NewServices(i do.Injector) (*Services, error) {
	return &Services{
		Queries: do.MustInvoke[*queries.Queries](i),
	}, nil
}

type Services struct {
	Queries *queries.Queries
}

func (services *Services) Register(email string, password string) (error) {
	userExists, userExistsErr := services.UserExists(email)

	if userExistsErr != nil {
		return userExistsErr
	}

	if userExists {
		return errors.New("user already exists")
	}

	newUserErr := services.Queries.NewUser(email, password)

	if newUserErr != nil {
		return newUserErr
	}

	return nil
}

func (services *Services) UserExists(email string) (bool, error) {
	return services.Queries.UserExists(email)
}

func (services *Services) Login(email string, password string) (bool, error) {
	user, userErr := services.Queries.GetUser(email)

	if userErr != nil {
		return false, userErr
	}

	if user == nil {
		return false, nil
	}

	if user.PasswordHash != password {
		return false, nil
	}

	return true, nil
}

func (services *Services) EmailLogin(email string) (bool, error) {
	userExists, userExistsErr := services.UserExists(email)

	if userExistsErr != nil {
		return false, userExistsErr
	}

	if !userExists {
		return false, nil
	}

	return true, nil
}

func (services *Services) DeleteUser(email string) (error) {
	userExists, userExistsErr := services.UserExists(email)

	if userExistsErr != nil {
		return userExistsErr
	}

	if !userExists {
		return errors.New("user does not exist")
	}

	deleteUserErr := services.Queries.DeleteUser(email)

	if deleteUserErr != nil {
		return deleteUserErr
	}

	return nil
}