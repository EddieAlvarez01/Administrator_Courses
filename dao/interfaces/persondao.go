package interfaces

import (
	"github.com/EddieAlvarez01/administrator_courses/models"

	"go.mongodb.org/mongo-driver/mongo/options"
)

//PersonDao .
type PersonDao interface {
	Create(p *models.Person) error
	Update(person *models.Person) error
	GetOne(id string) (*models.Person, error)
	GetAllByFilter(filter interface{}, opt *options.FindOptions) ([]*models.Person, error)
	GetByEmail(email string) *models.Person
	CreateProfessor(p *models.Person) error
}
