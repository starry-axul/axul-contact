package bootstrap

import (
	"fmt"
	"github.com/ncostamagna/axul_domain/domain"
	"github.com/ncostamagna/go-logger-hub/loghub"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

func SetupLogger() loghub.Logger {
	return loghub.New(
		loghub.NewNativeLogger(nil, loghub.FormatStringToNumber(os.Getenv("NATIVE_LOGGER_TRACE"))),
	)
}

func DBConnection() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"))
	fmt.Println("connect: ", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if os.Getenv("DATABASE_DEBUG") == "true" {
		db = db.Debug()
	}

	if os.Getenv("DATABASE_MIGRATE") == "true" {
		if err := db.AutoMigrate(domain.Contact{}); err != nil {
			panic(err)
		}
	}

	return db
}
