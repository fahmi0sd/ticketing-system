package route

import (
	"context"
	"strings"

	"github.com/fahmi0sd/ticketing-system/service/route"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) GetAll(filter route.SearchFilter) ([]route.Route, error) {
	ctx := context.Background()
	query := r.db.WithContext(ctx).Model(&route.Route{})

	if filter.Origin != "" {
		query = query.Where("LOWER(origin) LIKE ?", "%"+strings.ToLower(filter.Origin)+"%")
	}
	if filter.Destination != "" {
		query = query.Where("LOWER(destination) LIKE ?", "%"+strings.ToLower(filter.Destination)+"%")
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Date != "" {
		query = query.Where("DATE(departure_at) = ?", filter.Date)
	}

	query = query.Where("departure_at > NOW()").Order("departure_at ASC")

	var routes []route.Route
	err := query.Find(&routes).Error
	return routes, err
}

func (r *GormRepository) UpdateSold(id, newSold int) error {
	ctx := context.Background()
	return r.db.WithContext(ctx).
		Model(&route.Route{}).
		Where("id = ?", id).
		Update("sold", newSold).Error
}
