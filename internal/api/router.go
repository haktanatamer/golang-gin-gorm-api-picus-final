package api

import (
	"api-gin/package/internal/api/category_api"
	"api-gin/package/internal/api/product_api"
	"api-gin/package/internal/api/shopping_api"
	"api-gin/package/internal/api/user_api"
	"api-gin/package/internal/domain/category"
	"api-gin/package/internal/domain/product"
	"api-gin/package/internal/domain/shopping"
	"api-gin/package/internal/domain/users"
	"api-gin/package/pkg/database"
	"api-gin/package/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	_ "api-gin/package/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterHandlers create endpoints
func RegisterHandlers(r *gin.Engine) {

	swaggerGroup := r.Group("/swagger")
	swaggerGroup.GET("/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	db := database.DB

	userRepository := users.NewUserRepository(db)
	userService := users.NewUserService(*userRepository)
	userController := user_api.NewUserController(userService)

	userGroup := r.Group("/user")
	userGroup.POST("/add", middleware.DBTransactionMiddleware(db), userController.Add)
	userGroup.POST("/login", userController.Login)

	categoryRepository := category.NewCategoryRepository(db)
	categoryService := category.NewCategoryService(*categoryRepository)
	categoryController := category_api.NewCategoryController(categoryService)

	categoryGroup := r.Group("/category")
	categoryGroup.POST("/add", middleware.AuthMiddleware(viper.GetString("server.secret"), "admin"), categoryController.Add)
	categoryGroup.POST("/", categoryController.GetAllByPagination)

	productRepository := product.NewProductRepository(db)
	productService := product.NewProductService(*productRepository)
	productController := product_api.NewProductController(productService)

	productGroup := r.Group("/product")
	productGroup.POST("/add", middleware.AuthMiddleware(viper.GetString("server.secret"), "admin"), middleware.DBTransactionMiddleware(db), productController.Add)
	productGroup.POST("/", productController.GetAllByPagination)
	productGroup.POST("/search", productController.Search)
	productGroup.POST("/delete", middleware.AuthMiddleware(viper.GetString("server.secret"), "admin"), productController.Delete)
	productGroup.POST("/update", middleware.AuthMiddleware(viper.GetString("server.secret"), "admin"), productController.Update)

	shoppingRepository := shopping.NewShoppingRepository(db)
	shoppingService := shopping.NewShoppingService(*shoppingRepository)
	shoppingController := shopping_api.NewShoppingController(shoppingService)

	shoppingGroup := r.Group("/shopping")
	shoppingGroup.POST("/add", middleware.AuthMiddleware(viper.GetString("server.secret"), "customer"), middleware.DBTransactionMiddleware(db), shoppingController.AddToCart)
	shoppingGroup.GET("/list", middleware.AuthMiddleware(viper.GetString("server.secret"), "customer"), shoppingController.GetCart)
	shoppingGroup.POST("/delete", middleware.AuthMiddleware(viper.GetString("server.secret"), "customer"), middleware.DBTransactionMiddleware(db), shoppingController.DeleteToCart)

	orderGroup := r.Group("/order")
	orderGroup.POST("/add", middleware.AuthMiddleware(viper.GetString("server.secret"), "customer"), middleware.DBTransactionMiddleware(db), shoppingController.AddOrder)
	orderGroup.GET("/list", middleware.AuthMiddleware(viper.GetString("server.secret"), "customer"), shoppingController.GetOrders)
	orderGroup.POST("/cancel", middleware.AuthMiddleware(viper.GetString("server.secret"), "customer"), middleware.DBTransactionMiddleware(db), shoppingController.CancelOrders)
}

// RegisterHandlersTest create endpoints for test
func RegisterHandlersTest(r *gin.Engine) {
	userGroup := r.Group("/user")
	userGroup.POST("/add", testMethod())

	orderGroup := r.Group("/order")
	orderGroup.GET("/list", testMethod())

}

func testMethod() gin.HandlerFunc {
	return func(g *gin.Context) {
		g.JSON(http.StatusOK, gin.H{
			"result": "func",
		})
	}
}
