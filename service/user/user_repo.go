package user

type Repository interface {
	Create(u User) (User, error)
	GetByEmail(email string) (User, error)
	GetByID(id int) (User, error)
	ExistsByEmail(email string) bool
}
