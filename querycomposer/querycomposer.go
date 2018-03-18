package querycomposer

type QueryComposer interface {
	Columns(input interface{}) QueryComposer
	PersistenceNames(names []string) QueryComposer
	Compose() string
}

type UpdateQueryComposer interface {
	QueryComposer
	Where(whereClause string) SelectQueryComposer
}

type SelectQueryComposer interface {
	UpdateQueryComposer
	OrderBy(orderClause string) SelectQueryComposer
	Paginate(offset, limit int) SelectQueryComposer
}
