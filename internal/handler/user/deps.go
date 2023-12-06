package user

import (
	"github.com/Nahbox/crud-users-api/internal/models"
)

// TODO: add mock tests
type DB interface {
	GetAllUsers() ([]models.User, error)
	AddUser(user *models.User) error
	GetUserById(userId int) (models.User, error)
	UpdateUser(userId int, user *models.User) (models.User, error)
	DeleteUser(userId int) error
}
