package order

type Repository struct {
	db database
}

func New(d database) *Repository {
	return &Repository{
		db: d,
	}
}
