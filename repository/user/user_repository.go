package user

import (
	"context"

	"github.com/fahmi0sd/ticketing-system/service/user"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(u user.User) (user.User, error) {
	ctx := context.Background()
	result := r.db.WithContext(ctx).Create(&u)
	return u, result.Error
}

func (r *GormRepository) GetByEmail(email string) (user.User, error) {
	ctx := context.Background()
	var u user.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&u)
	return u, result.Error
}

func (r *GormRepository) GetByID(id int) (user.User, error) {
	ctx := context.Background()
	var u user.User
	result := r.db.WithContext(ctx).First(&u, id)
	return u, result.Error
}

func (r *GormRepository) ExistsByEmail(email string) bool {
	ctx := context.Background()
	var count int64
	r.db.WithContext(ctx).Model(&user.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}
