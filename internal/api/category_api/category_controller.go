package category_api

import (
	"api-gin/package/internal/domain/category"
	"api-gin/package/pkg/helper"
	"encoding/csv"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	categoryService *category.CategoryService
}

func NewCategoryController(service *category.CategoryService) *CategoryController {
	return &CategoryController{
		categoryService: service,
	}
}

// @Summary Add category by CSV file
// @Tags Category
// @Accept  json
// @Produce  json
// @Param categoryCsv formData file true "Category CSV File"
// @Failure 409 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Success 201 {object} map[string]string
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router /category/add [post]
func (cc *CategoryController) Add(g *gin.Context) {

	file, _, err := g.Request.FormFile("categoryCsv")

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error_message": err.Error()})
		g.Abort()
		return
	}

	csvLines, readErr := csv.NewReader(file).ReadAll()
	if readErr != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error_message": readErr.Error()})
		g.Abort()
		return
	}

	var result []category.Category

	for _, line := range csvLines[1:] {

		data := category.Category{
			Name: line[0], Active: true,
		}

		result = append(result, data)
	}

	var notAddedCategories []string

	for _, row := range result {
		hasCategory := cc.categoryService.ExistByName(row)

		if hasCategory {
			notAddedCategories = append(notAddedCategories, row.Name)
		} else {
			cc.categoryService.Create(&row)
		}

	}

	if len(notAddedCategories) > 0 {
		g.JSON(http.StatusConflict, gin.H{
			"status": "This list has not been added because already exists.", "data": strings.Join(notAddedCategories, ","),
		})
		g.Abort()
		return
	}

	g.JSON(http.StatusCreated, gin.H{
		"status": "Categories added.",
	})
}

func (cc *CategoryController) GetAll(g *gin.Context) {

	categories := cc.categoryService.GetAll()

	g.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}

// @Summary Gets all categories with paginated result
// @Tags Category
// @Accept  json
// @Produce  json
// @Param paginationRequest body CategoryGetAllByPaginationRequest false "pagination"
// @Failure 400 {object} map[string]string
// @Success 200 {object} map[string]helper.Pages
// @Router /category/ [post]
func (cc *CategoryController) GetAllByPagination(g *gin.Context) {

	var req CategoryGetAllByPaginationRequest

	httpCode, errDetail := helper.BindAndValid(g, &req)

	if httpCode != http.StatusOK {
		g.JSON(httpCode, gin.H{"error_message": errDetail})
		g.Abort()
		return
	}

	pages := helper.NewPages(req.Page, req.Limit)
	categories := cc.categoryService.GetAllByPagination(pages.Page, pages.Limit)
	pages.Items = categories
	g.JSON(http.StatusOK, gin.H{
		"categories": pages,
	})
}
