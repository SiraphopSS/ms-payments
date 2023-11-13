package main

import (
	"log"
	"os"
	"promptlabth/ms-payments/controllers"
	"promptlabth/ms-payments/database"
	"promptlabth/ms-payments/repository"
	"promptlabth/ms-payments/routes"
	"promptlabth/ms-payments/usecases"

	"gorm.io/driver/postgres"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

func CORSMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {

			c.AbortWithStatus(204)

			return

		}

		c.Next()

	}

}

var err error

func main() {
	database.DB, err = gorm.Open(postgres.Open(database.DbURL(database.BuildDBConfig())), &gorm.Config{})
	// defer database.DB.Close()
	if err != nil {
		log.Fatal("database connect error: ", err)
	}
	// auto migrate
	// database.DB.AutoMigrate(
	// 	&Coin{},
	// 	&Feature{},
	// 	&Payment{},
	// 	&PaymentMethod{},
	// 	&Feature{},
	// 	&User{},
	// 	&PaymentSubscription{},
	// )
	// database.DB.AutoMigrate()

	repo := &repository.PaymentRepository{}
	db, err := database.DB.DB()
	if err != nil {
		log.Fatal(err)
	}
	repo.DB = db
	usecase := usecases.NewPaymentUsecase(repo)
	controller := controllers.PaymentController{Usecase: usecase}
	if os.Getenv("BaseOn") != "DEV" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	// to set CORS
	r.Use(CORSMiddleware())
	// the clean arch
	routes.CoinRoute(r, database.DB)
	routes.PaymentSubscriptionRoute(r, database.DB)

	r.POST("/payment", controller.CreatePayment)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "world"})
	})

	r.Run(":8080")
}
