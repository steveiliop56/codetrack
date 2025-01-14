package services

import (
	"codetrack/internal/queries"

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
	newUserErr := services.Queries.NewUser(email, password)

	if newUserErr != nil {
		return newUserErr
	}

	return nil
}

func (services *Services) UserExists(email string) (bool, error) {
	return services.Queries.UserExists(email)
}