package booking

import "time"

type Repository interface {
	Create(b Booking) (Booking, error)
	GetByID(id int) (Booking, error)
	GetByUserID(userID int) ([]Booking, error)
	GetByExternalID(externalID string) (Booking, error)
	UpdateStatus(id int, status Status, paidAt *time.Time) error
	FindWithRelations(id int) (Booking, error)
}
