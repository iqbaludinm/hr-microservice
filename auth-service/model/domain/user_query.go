package domain

type UserQueryFilter struct {
	Name  string 
	Email string 
	Phone string 

	// ShowDeleted is used for showing the soft-deleted user or not.
	ShowDeleted bool

	// Pagination is used for fetching the data by page. The default value is 1.
	Page  string 
	Limit string 
}

func (q *UserQueryFilter) BuildUserQueries() (filter string) {
	// filter user by name
	if q.Name != "" {
		filter += " AND u.name LIKE '%" + q.Name + "%'"
	}

	// filter user by email
	if q.Email != "" {
		filter += " AND u.email = '" + q.Email + "'"
	}

	// filter user by phone
	if q.Phone != "" {
		filter += " AND u.phone = " + q.Phone
	}

	// filter user by deleted_at is null
	if !q.ShowDeleted {
		filter += " AND u.deleted_at is null"
	}

	// remove the first ' AND ' from the filter
	if len(filter) > 0 {
		filter = filter[5:]
	}

	return filter
}
