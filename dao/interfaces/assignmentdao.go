package interfaces

import (
	"github.com/EddieAlvarez01/administrator_courses/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AssignmentDao interface {
	Create(assignment *models.Assignment) error
	Update(assignment *models.Assignment) error
	GetByPersonAndDateRange(idPerson string, startDate, endDate time.Time) (*models.Assignment, error)
	GetById(assignmentId string) (*models.Assignment, error)
	GetAllBySectionIdInAPeriod(sectionID primitive.ObjectID, startDate, endDate time.Time) ([]*models.Assignment, error)
}
