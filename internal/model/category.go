package model

type Category struct {
	ID       int64
	Name     string
	ParentID int64
}

type CategoryExtended struct {
	ID       int64
	Name     string
	Parent   *Category
	Children []*Category
}
