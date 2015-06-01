package services

// DefaultPerPage is the default number of items to return per page in a
// paginated result set.
const DefaultPerPage = 10

// ListOptions specifies general pagination options for fetching a list of
// results.
type ListOptions struct {
	PerPage int `url:"per_page,omitempty" json:"per_page,omitempty"`
	Page    int `url:"page,omitempty" json:"page,omitempty"`
}

func (o ListOptions) PageOrDefault() int {
	if o.Page <= 0 {
		return 1
	}

	return o.Page
}

func (o ListOptions) PerPageOrDefault() int {
	if o.PerPage <= 0 {
		return DefaultPerPage
	}

	return o.PerPage
}

// Offset returns the offset - in number of items - that this page represents.
// It is suitable for using in a SQL 'OFFSET' clause.
func (o ListOptions) Offset() int {
	return (o.PageOrDefault() - 1) * o.PerPageOrDefault()
}
