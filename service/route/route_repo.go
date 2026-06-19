package route

type Repository interface {
	GetAll(filter SearchFilter) ([]Route, error)
	UpdateSold(id, newSold int) error
	GetByID(id int) (Route, error)
}
