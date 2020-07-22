package interfaces

import (
	"github.com/EddieAlvarez01/administrator_courses/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SectionDao interface {
	Create(section *models.Section) error
	GetByCourseAndSection(idCourse primitive.ObjectID, section string) *models.Section
	GetProfessor(idProfessor primitive.ObjectID) *models.Person
	Update(section *models.Section) error
	GetById(id string) (*models.Section, error)
	GetAllByCourseID(id primitive.ObjectID) ([]models.Section, error)
}
