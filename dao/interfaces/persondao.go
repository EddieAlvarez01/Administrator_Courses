package interfaces

import (
	"github.com/EddieAlvarez01/administrator_courses/models"
)

//PersonDao .
type PersonDao interface {
	Create(p *models.Person) error
}
