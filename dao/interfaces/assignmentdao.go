package interfaces

import (
	"github.com/EddieAlvarez01/administrator_courses/models"
	"time"
)

type AssignmentDao interface {
	Create(assignment *models.Assignment) error
	GetByPersonAndDateRange(idPerson string, startDate, endDate time.Time) (*models.Assignment, error)
}
