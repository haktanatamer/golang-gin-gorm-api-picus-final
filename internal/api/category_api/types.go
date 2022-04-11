package category_api

type CategoryRequest struct {
	Name string `json:"name" valid:"Required"`
}

type CategoryResponse struct {
	Name string `json:"name"`
}

type CategoryGetAllByPaginationRequest struct {
	Page  int `json:"page" valid:"Match(^[0-9]*$)"`
	Limit int `json:"limit" valid:"Match(^[0-9]*$)"`
}

type CategoryFileRequest struct {
	Csv []byte `json:"categoryCsv" valid:"Required"`
}
