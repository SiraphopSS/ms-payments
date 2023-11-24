package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/promptlabth/ms-payments/controllers"
	"github.com/promptlabth/ms-payments/middlewares"
	"github.com/promptlabth/ms-payments/repository"
	"github.com/promptlabth/ms-payments/services"
	"github.com/promptlabth/ms-payments/usecases"
	"gorm.io/gorm"
)

func SubscriptionRoute(r *gin.Engine, DB *gorm.DB) {

	// initial a repo, usecase, control
	userRepo := repository.NewUserRepository(DB)
	userUseCases := usecases.NewUserUseCase(userRepo)

	planRepo := repository.NewPlanRepository(DB)
	planUsecase := usecases.NewPlanUsecase(planRepo)

	paymentSubscriptionRepo := repository.NewPaymentScriptionRepository(DB)
	paymentSubscriptionUsecase := usecases.NewPaymentSubscriptionUsecase(paymentSubscriptionRepo)

	subscriptionReqUrlController := controllers.NewSubscriptionReqUrlController(
		userUseCases,
		planUsecase,
		paymentSubscriptionUsecase,
	)

	// use a middleware to route subscription
	subScription := r.Group("/subscription")
	protect := subScription.Use(middlewares.AuthorizeFirebase())

	protect.POST("/get-url", subscriptionReqUrlController.GetSubscriptionUrl)

	protect.POST("/success", subscriptionReqUrlController.SaveSubscription)

	protect.POST("/cancle", func(c *gin.Context) {
		data, _ := services.CancelSubscriptionBySubID(
			"sub_1OD9dUAom1IgIvKKHzZyTo22",
		)
		c.JSON(200, gin.H{
			"data": data,
		})
	})
}
