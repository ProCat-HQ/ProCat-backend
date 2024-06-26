package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/procat-hq/procat-backend/internal/app/service"
	v3 "github.com/swaggest/swgui/v3"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.StaticFS("/assets", gin.Dir("./assets", false))
	router.StaticFS("/docs", gin.Dir("./api", false))

	swaggerHandler := v3.NewHandler("ProCat API", "/docs/api.json", "/swagger")

	router.GET("/swagger/*any", gin.WrapH(swaggerHandler))

	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	m.SetSlowTime(5)
	m.Use(router)

	users := router.Group("/users")
	{
		usersAuthenticatedGroup := users.Group("", h.UserIdentify)
		{
			usersAuthenticatedGroup.GET("", h.CheckRole("admin"), h.GetAllUsers)
			usersAuthenticatedGroup.GET("/:id", h.MustBelongsToUser, h.GetUser)
			usersAuthenticatedGroup.DELETE("/:id", h.CheckRole("admin"), h.DeleteUser)

			usersAuthenticatedGroup.POST("/logout", h.Logout)
		}
		users.POST("/sign-in", h.SignIn)
		users.POST("/sign-up", h.SignUp)
		users.POST("/refresh", h.RefreshToken)

		verification := users.Group("/verification", h.UserIdentify)
		{
			verification.PUT("/send", h.SendCode)     // TODO
			verification.POST("/check", h.CheckCode)  // TODO
			verification.POST("/iin-bin", h.CheckIIN) // TODO
		}

		change := users.Group("/change", h.UserIdentify)
		{
			change.POST("/iin-bin", h.ChangeIIN)
			change.POST("/fullname", h.ChangeFullName)
			change.POST("/password", h.ChangePassword)
			change.POST("/phone", h.ChangePhone)
			change.POST("/email", h.ChangeEmail)
			change.PATCH("/role/:id", h.CheckRole("admin"), h.ChangeRole)
		}

		deliverymen := users.Group("/deliverymen", h.UserIdentify, h.CheckRole("deliveryman"))
		{
			deliverymen.GET("", h.GetAllDeliverymen)
			deliverymen.GET("/:id", h.MustBelongsToUser, h.GetDeliveryman)
			deliverymen.POST("/:id", h.CheckRole("admin"), h.CreateDeliveryman)
			deliverymen.PATCH("/:id", h.CheckRole("admin"), h.ChangeDeliverymanData)
			deliverymen.DELETE("/:id", h.CheckRole("admin"), h.DeleteDeliveryman)

			deliveries := deliverymen.Group("/deliveries")
			{
				deliveries.GET("", h.CheckRole("admin"), h.GetAllDeliveries)
				deliveries.GET("/:id", h.GetAllDeliveriesForOneDeliveryman)
				deliveries.GET("/delivery/:id", h.GetDelivery)
				deliveries.PATCH("/:id", h.ChangeDeliveryStatus)
				deliveries.POST("/create-route", h.CreateRoute)
			}
		}

		admin := users.Group("/admin", h.UserIdentify, h.CheckRole("admin"))
		{
			admin.POST("/cluster", h.Cluster)
			admin.GET("/deliveries-to-sort", h.GetAllDeliveriesToSort)
			admin.PATCH("/change-delivery", h.ChangeDeliveryData)
		}

		cart := users.Group("/cart", h.UserIdentify)
		{
			cart.GET("", h.GetCartItems)
			cart.POST("", h.AddItemsToCart)
			cart.DELETE("/:id", h.DeleteItemsFromCart)
		}

		orders := users.Group("/orders", h.UserIdentify)
		{
			orders.GET("", h.GetAllOrders)
			orders.GET("/:id", h.GetOrder)
			orders.POST("", h.CreateOrder)
			orders.PATCH("/cancel/:id", h.CancelOrder)
			orders.PATCH("/extend/:id", h.ExtendOrder)
			orders.POST("/return/:id", h.ReturnOrder)
			orders.POST("/confirm-extension/:id", h.CheckRole("admin"), h.ConfirmOrderExtension)
			orders.PATCH("/status/:id", h.CheckRole("admin"), h.ChangeOrderStatus)
			orders.POST("/need-repair/:id", h.CheckRole("admin"), h.NeedRepairForOrder)

			payment := orders.Group("/payment")
			{
				payment.GET("/:id", h.GetPaymentData)
				payment.PATCH("/:id", h.CheckRole("admin"), h.ChangePaymentStatus)
			}
		}

		subscriptions := users.Group("/subscriptions", h.UserIdentify)
		{
			subscriptions.GET("", h.GetAllSubscriptions)
			subscriptions.POST("", h.SubscribeToItem)
			subscriptions.DELETE("/:id", h.DeleteItemFromSubscriptions)
		}

		notifications := users.Group("/notifications", h.UserIdentify)
		{
			notifications.GET("", h.GetAllNotifications)
			notifications.POST("/:id", h.CheckRole("admin"), h.SendNotification)
			notifications.PATCH("/:id", h.ViewNotification)
			notifications.DELETE("/:id", h.CheckRole("admin"), h.DeleteNotification)
		}

	}

	categories := router.Group("/categories")
	{
		categories.GET("/:id", h.GetCategoriesForParent)
		categories.GET("/route/:id", h.GetCategoryRoute)
		categories.POST("/:id", h.UserIdentify, h.CheckRole("admin"), h.CreateCategory)
		categories.PATCH("/:id", h.UserIdentify, h.CheckRole("admin"), h.ChangeCategory)
		categories.DELETE("/:id", h.UserIdentify, h.CheckRole("admin"), h.DeleteCategory)
	}

	items := router.Group("/items")
	{
		items.GET("", h.GetAllItems)
		items.GET("/:id", h.GetItem)
		items.POST("", h.UserIdentify, h.CheckRole("admin"), h.CreateItem)
		items.PATCH("/:id", h.UserIdentify, h.CheckRole("admin"), h.ChangeItem)
		items.DELETE("/:id", h.UserIdentify, h.CheckRole("admin"), h.DeleteItem)

		stock := items.Group("/stock", h.UserIdentify, h.CheckRole("admin"))
		{
			stock.PUT("/:id", h.ChangeStock)
		}

		infos := items.Group("/infos", h.UserIdentify, h.CheckRole("admin"))
		{
			infos.POST("/:id", h.AddInfo)
			infos.PATCH("/:id", h.ChangeInfo)
			infos.DELETE("/:id", h.DeleteInfo)
		}

		images := items.Group("/images", h.UserIdentify, h.CheckRole("admin"))
		{
			images.POST("/:id", h.AddImages)
			images.DELETE("/:id", h.DeleteImages)
		}
	}

	stores := router.Group("/stores")
	{
		stores.GET("", h.GetAllStores)
		stores.POST("", h.UserIdentify, h.CheckRole("admin"), h.CreateStore)
		stores.PATCH("/:id", h.UserIdentify, h.CheckRole("admin"), h.ChangeStore)
		stores.DELETE("/:id", h.UserIdentify, h.CheckRole("admin"), h.DeleteStore)
	}

	return router
}
