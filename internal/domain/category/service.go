package category

type CategoryService struct {
	r CategoryRepository
}

func NewCategoryService(r CategoryRepository) *CategoryService {
	return &CategoryService{
		r: r,
	}
}

func (cs *CategoryService) ExistByName(c Category) bool {
	hasUser := cs.r.ExistByName(c)

	return hasUser
}

func (cs *CategoryService) Create(newC *Category) error {
	err := cs.r.Create(newC)

	return err
}

func (cs *CategoryService) GetAll() []Category {
	result := cs.r.GetAll()

	return result
}

func (cs *CategoryService) GetAllByPagination(page, limit int) []Category {
	result := cs.r.GetAllByPagination(page, limit)

	return result
}
