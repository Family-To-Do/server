package main

import (
	"./models"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

const PathWeb = "./web/dist/"

func Init() *gorm.DB {
	/* - Import Environment - */
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	/* - Connect to Database - */
	db, err := gorm.Open(
		"postgres",
		"host="+os.Getenv("DATABASE_HOST")+" "+
			"port="+os.Getenv("DATABASE_PORT")+" "+
			"user="+os.Getenv("DATABASE_USER")+" "+
			"dbname="+os.Getenv("DATABASE_NAME")+" "+
			"password="+os.Getenv("DATABASE_PASS")+" "+
			"sslmode="+os.Getenv("DATABASE_SSL"))

	if err != nil {
		panic(err)
	}

	return db
}

func main() {
	db := Init()
	defer db.Close()

	/* - Initialization Iris - */
	app := iris.New()
	app.Logger().SetLevel(os.Getenv("LOGGER"))

	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())

	// SPA (git submodule, dist folder)
	app.StaticWeb("/", PathWeb)

	// FIXME Example
	app.Get("/api/users", func(ctx context.Context) {
		var user models.User
		db.First(&user, 1)

		ctx.JSON(user)
	})

	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
