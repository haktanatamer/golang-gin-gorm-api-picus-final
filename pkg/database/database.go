package database

import (
	"api-gin/package/internal/domain/category"
	"api-gin/package/internal/domain/product"
	"api-gin/package/internal/domain/shopping"
	"api-gin/package/internal/domain/users"
	"fmt"
	"time"

	"github.com/spf13/viper"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DB    *gorm.DB
	err   error
	DBErr error
)

type Database struct {
	*gorm.DB
}

// Setup initialize db
func Setup() error {
	var db = DB

	db, err := gorm.Open(mysql.Open(getConnectionString()), &gorm.Config{
		PrepareStmt: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   true,
		},
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return err
	}

	db.AutoMigrate(users.User{}, users.User_Role{}, users.Role{}, users.User_Token{}, category.Category{}, product.Product{}, shopping.Cart{}, shopping.Cart_Item{}, shopping.Order{}, shopping.Order_Item{})
	DB = db

	return nil
}

// getConnectionString get db connection string data
func getConnectionString() string {

	dbname := viper.GetString("database.dbname")
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")

	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, host, dbname)
}

func GetDB() *gorm.DB {
	return DB
}

// TestSetup initialize testdb
func TestSetup() (*gorm.DB, error) {

	db, err := gorm.Open(mysql.Open("root:1575@tcp(127.0.0.1:3306)/basket_api_test?parseTime=true&loc=Local"), &gorm.Config{
		PrepareStmt: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   true,
		},
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	return db, err

}
