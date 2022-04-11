package helper

var (
	DefaultPage  = 1
	DefaultLimit = 4
)

type Pages struct {
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Items interface{} `json:"items"`
}

func NewPages(page, limit int) *Pages {
	if page <= 0 {
		page = DefaultPage
	}

	if limit <= 0 {
		limit = DefaultLimit
	}

	return &Pages{
		Page:  page,
		Limit: limit,
	}
}
