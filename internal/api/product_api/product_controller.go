package product_api

import (
	"api-gin/package/internal/domain/product"
	"api-gin/package/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productService *product.ProductService
}

func NewProductController(service *product.ProductService) *ProductController {
	return &ProductController{
		productService: service,
	}
}

// @Summary Add product
// @Tags Product
// @Accept  json
// @Produce  json
// @Param productRequest body ProductRequest false "Product Info"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Success 200 {object} map[string]string
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router /product/add [post]
func (pc *ProductController) Add(g *gin.Context) {

	var req ProductRequest

	httpCode, errDetail := helper.BindAndValid(g, &req)

	if httpCode != http.StatusOK {
		g.JSON(httpCode, gin.H{"error_message": errDetail})
		g.Abort()
		return
	}

	hasCategory := pc.productService.GetCategoryById(req.CategoryId)

	if !hasCategory {
		allCategories := pc.productService.GetAllCategories()
		g.JSON(http.StatusBadRequest, gin.H{"error_message": " Category did not found, select category from the list ", "Categories": allCategories})
		g.Abort()
		return
	}

	hasProduct := pc.productService.ExistByName(req.Name)

	if hasProduct {
		g.JSON(http.StatusBadRequest, gin.H{"error_message": "Product Already Exists."})
		g.Abort()
		return
	}

	newProduct := product.NewProduct(req.Name, req.Brand, req.CategoryId)

	err := pc.productService.Create(newProduct)

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{
			"error_message": "An error while adding user.",
		})
		g.Abort()
		return
	}

	g.JSON(http.StatusCreated, gin.H{
		"status": "The product has been created",
	})
}

func (pc *ProductController) GetAll(g *gin.Context) {

	products := pc.productService.GetAll()

	var resProduct []ResponseRequest

	for _, p := range products {
		resProduct = append(resProduct, ResponseRequest{Name: p.Name, Brand: p.Brand, Id: int(p.Id), Sku: p.Sku, Category: p.Category.Name})
	}

	g.JSON(http.StatusOK, gin.H{
		"categories": resProduct,
	})
}

// @Summary Gets searched products
// @Tags Product
// @Accept  json
// @Produce  json
// @Param searchRequest body SearchRequest false "Search Value"
// @Failure 400 {object} map[string]string
// @Success 200 {object} map[string]string
// @Router /product/search [post]
func (pc *ProductController) Search(g *gin.Context) {

	var req SearchRequest

	httpCode, errDetail := helper.BindAndValid(g, &req)

	if httpCode != http.StatusOK {
		g.JSON(httpCode, gin.H{"error_message": errDetail})
		g.Abort()
		return
	}

	products := pc.productService.Search(req.Value)

	var resProduct []ResponseRequest

	for _, p := range products {
		resProduct = append(resProduct, ResponseRequest{Name: p.Name, Brand: p.Brand, Id: int(p.Id), Sku: p.Sku, Category: p.Category.Name})
	}

	if len(resProduct) == 0 {
		g.JSON(httpCode, gin.H{"result": "No records to show."})
		g.Abort()
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"categories": resProduct,
	})
}

// @Summary Delete product
// @Tags Product
// @Accept  json
// @Produce  json
// @Param deleteRequest body DeleteRequest false "Product Id"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Success 200 {object} map[string]string
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router /product/delete [post]
func (pc *ProductController) Delete(g *gin.Context) {

	var req DeleteRequest

	httpCode, errDetail := helper.BindAndValid(g, &req)

	if httpCode != http.StatusOK {
		g.JSON(httpCode, gin.H{"error_message": errDetail})
		g.Abort()
		return
	}

	hasProduct := pc.productService.ExistById(req.Id)

	if !hasProduct {
		g.JSON(http.StatusBadRequest, gin.H{"error_message": "The product has already been deleted."})
		g.Abort()
		return
	}

	err := pc.productService.DeleteById(req.Id)

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{
			"error_message": "An error while deleting product.",
		})
		g.Abort()
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"categories": "The product has been deleted",
	})
}

// @Summary Update product
// @Tags Product
// @Accept  json
// @Produce  json
// @Param updateRequest body UpdateRequest false "Products Update Values"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Success 200 {object} map[string]string
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router /product/update [post]
func (pc *ProductController) Update(g *gin.Context) {

	var req UpdateRequest

	httpCode, errDetail := helper.BindAndValid(g, &req)

	if httpCode != http.StatusOK {
		g.JSON(httpCode, gin.H{"error_message": errDetail})
		g.Abort()
		return
	}

	wUpdatedProduct := product.NewProduct(req.Name, req.Brand, req.CategoryId)
	wUpdatedProduct.Id = uint(req.Id)

	err := pc.productService.Update(wUpdatedProduct)

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{
			"error_message": err.Error(),
		})
		g.Abort()
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"categories": "The product has been updated",
	})
}

// @Summary Gets all products with paginated result
// @Tags Product
// @Accept  json
// @Produce  json
// @Param paginationRequest body ProductGetAllByPaginationRequest false "pagination"
// @Failure 400 {object} map[string]string
// @Success 200 {object} map[string]helper.Pages
// @Router /product/ [post]
func (pc *ProductController) GetAllByPagination(g *gin.Context) {

	var req ProductGetAllByPaginationRequest

	httpCode, errDetail := helper.BindAndValid(g, &req)

	if httpCode != http.StatusOK {
		g.JSON(httpCode, gin.H{"error_message": errDetail})
		g.Abort()
		return
	}
	pages := helper.NewPages(req.Page, req.Limit)

	products := pc.productService.GetAllByPagination(pages.Page, pages.Limit)

	var resProduct []ResponseRequest

	for _, p := range products {
		resProduct = append(resProduct, ResponseRequest{Name: p.Name, Brand: p.Brand, Id: int(p.Id), Sku: p.Sku, Category: p.Category.Name})
	}
	pages.Items = resProduct
	g.JSON(http.StatusOK, gin.H{
		"produts": pages,
	})
}
