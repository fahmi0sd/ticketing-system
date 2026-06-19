package booking

import (
	"context"
	"time"

	"github.com/fahmi0sd/ticketing-system/service/booking"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(b booking.Booking) (booking.Booking, error) {
	ctx := context.Background()
	result := r.db.WithContext(ctx).Create(&b)
	return b, result.Error
}

func (r *GormRepository) GetByID(id int) (booking.Booking, error) {
	ctx := context.Background()
	var b booking.Booking
	result := r.db.WithContext(ctx).First(&b, id)
	return b, result.Error
}

func (r *GormRepository) GetByUserID(userID int) ([]booking.Booking, error) {
	ctx := context.Background()
	var bookings []booking.Booking
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&bookings)
	return bookings, result.Error
}

func (r *GormRepository) GetByExternalID(externalID string) (booking.Booking, error) {
	ctx := context.Background()
	var b booking.Booking
	result := r.db.WithContext(ctx).Where("external_id = ?", externalID).First(&b)
	return b, result.Error
}

func (r *GormRepository) UpdateStatus(id int, status booking.Status, paidAt *time.Time) error {
	ctx := context.Background()
	updates := map[string]any{"status": status}
	if paidAt != nil {
		updates["paid_at"] = paidAt
	}
	return r.db.WithContext(ctx).
		Model(&booking.Booking{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *GormRepository) FindWithRelations(id int) (booking.Booking, error) {
	ctx := context.Background()
	var b booking.Booking
	result := r.db.WithContext(ctx).
		Preload("User").
		Preload("Route").
		First(&b, id)
	return b, result.Error
}
