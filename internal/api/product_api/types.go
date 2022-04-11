package product_api

type ProductRequest struct {
	Name       string `json:"name" valid:"Required"`
	Brand      string `json:"brand" valid:"Required"`
	CategoryId int    `json:"categoryId" valid:"Match(^[0-9]*$); Required"`
}

type ResponseRequest struct {
	Name     string `json:"name"`
	Brand    string `json:"brand"`
	Id       int    `json:"id"`
	Sku      string `json:"sku"`
	Category string `json:"category"`
}

type SearchRequest struct {
	Value string `json:"value" valid:"Required"`
}

type DeleteRequest struct {
	Id int `json:"id" valid:"Match(^[0-9]*$); Required"`
}

type UpdateRequest struct {
	Id         int    `json:"id" valid:"Match(^[0-9]*$); Required"`
	Name       string `json:"name"`
	Brand      string `json:"brand"`
	CategoryId int    `json:"categoryId" valid:"Match(^[0-9]*$)"`
}

type ProductGetAllByPaginationRequest struct {
	Page  int `json:"page" valid:"Match(^[0-9]*$)"`
	Limit int `json:"limit" valid:"Match(^[0-9]*$)"`
}
