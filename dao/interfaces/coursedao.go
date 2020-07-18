package interfaces

import "github.com/EddieAlvarez01/administrator_courses/models"

type CourseDao interface {
	Create(course *models.Course) error
	GetById(id string) (*models.Course, error)
	Update(course *models.Course) error
}
