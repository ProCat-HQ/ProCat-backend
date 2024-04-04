package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	users := router.Group("/users")
	{
		usersAdminGroup := users.Group("", h.UserIdentify, h.CheckRole("admin"))
		{
			usersAdminGroup.GET("/", h.GetAllUsers)
			usersAdminGroup.GET("/:id", h.GetUser)
		}
		users.POST("/sign-in", h.SignIn)
		users.POST("/sign-up", h.SignUp)

		verification := users.Group("/verification", h.UserIdentify)
		{
			// неуверен тут насчёт прав доступа
			verification.PUT("/send", h.SendCode)
			verification.POST("/check", h.CheckCode)
			verification.POST("/iin-bin", h.CheckIIN)
		}

		change := users.Group("/change", h.UserIdentify)
		{
			change.POST("/iin-bin", h.ChangeIIN)
			change.POST("/fullname", h.ChangeFullName)
			change.POST("/password", h.ChangePassword)
			change.POST("/email", h.ChangeEmail)
			change.POST("/phone", h.ChangePhone)
			change.PATCH("/role/:id", h.CheckRole("admin"), h.ChangeRole)
		}

		deliverymen := users.Group("/deliverymen", h.UserIdentify, h.CheckRole("deliveryman"))
		{
			deliverymen.GET("/", h.GetAllDeliverymen)
			deliverymen.GET("/:id", h.GetDeliveryman)
			deliverymen.POST("/:id", h.CheckRole("admin"), h.CreateDeliveryman)
			deliverymen.PATCH("/:id", h.CheckRole("admin"), h.ChangeDeliverymanData)
			deliverymen.DELETE("/:id", h.CheckRole("admin"), h.DeleteDeliveryman)

			deliveries := deliverymen.Group("/deliveries")
			{
				deliveries.GET("/", h.GetAllDeliveries)
				deliveries.GET("/:id", h.GetDelivery)
				deliveries.PATCH("/:id", h.ChangeDeliveryStatus)
				deliveries.POST("/create-route", h.CreateRoute)
			}
		}

		admin := users.Group("/admin", h.UserIdentify, h.CheckRole("admin"))
		{
			admin.GET("/deliveries-to-sort", h.GetAllDeliveriesToSort)
			admin.PATCH("/change-delivery", h.ChangeDeliveryData)
			admin.POST("/cluster", h.Cluster)
		}

		cart := users.Group("/cart", h.UserIdentify)
		{
			cart.GET("/", h.GetCartItems)
			cart.POST("/", h.AddItemsToCart)
			cart.DELETE("/", h.DeleteItemsFromCart)
		}

		orders := users.Group("/orders", h.UserIdentify)
		{
			orders.GET("/", h.GetAllOrders)
			orders.GET("/:id", h.GetOrder)
			orders.POST("/", h.CreateOrder)
			orders.POST("/cancel/:id", h.CancelOrder)
			orders.PATCH("/status/:id", h.CheckRole("admin"), h.ChangeOrderStatus)

			payment := orders.Group("/payment")
			{
				payment.GET("/:id", h.GetPaymentData)
				payment.PATCH("/:id", h.CheckRole("admin"), h.ChangePaymentStatus)
			}
		}

		subscriptions := users.Group("/subscriptions", h.UserIdentify)
		{
			subscriptions.GET("/", h.GetAllSubscriptions)
			subscriptions.POST("/", h.SubscribeToItem)
			subscriptions.DELETE("/:id", h.DeleteItemFromSubscriptions)
		}

		notifications := users.Group("/notifications", h.UserIdentify)
		{
			notifications.GET("/", h.GetAllNotifications)
			notifications.POST("/:id", h.CheckRole("admin"), h.SendNotification)
			notifications.PATCH("/:id", h.ViewNotification)
			notifications.DELETE("/:id", h.CheckRole("admin"), h.DeleteNotification)
		}

	}

	categories := router.Group("/categories")
	{
		categories.GET("/:id", h.GetCategory)
		categories.GET("/route/:id", h.GetCategoryRoute)
		categories.POST("/:id", h.UserIdentify, h.CheckRole("admin"), h.CreateCategory)
		categories.PATCH("/:id", h.UserIdentify, h.CheckRole("admin"), h.ChangeCategory)
		categories.DELETE("/:id", h.UserIdentify, h.CheckRole("admin"), h.DeleteCategory)
	}

	items := router.Group("/items")
	{
		items.GET("/", h.GetAllItems)
		items.GET("/:id", h.GetItem)
		items.POST("/", h.UserIdentify, h.CheckRole("admin"), h.CreateItem)
		items.PATCH("/:id", h.UserIdentify, h.CheckRole("admin"), h.ChangeItem)
		items.DELETE("/:id", h.UserIdentify, h.CheckRole("admin"), h.DeleteItem)

		stock := items.Group("/stock", h.UserIdentify, h.CheckRole("admin"))
		{
			stock.POST("/:id", h.CreateStock)
			stock.PATCH("/:id", h.ChangeStock)
		}

		infos := items.Group("/infos", h.UserIdentify, h.CheckRole("admin"))
		{
			infos.POST("/:id", h.AddInfo)
			infos.PUT("/:id", h.ChangeInfo)
			infos.DELETE("/:id", h.DeleteInfo)
		}

		images := items.Group("/images", h.UserIdentify, h.CheckRole("admin"))
		{
			images.POST("/:id", h.AddImages)
			images.PATCH("/:id", h.ChangeImages)
			images.DELETE("/:id", h.DeleteImages)
		}
	}

	stores := router.Group("/stores")
	{
		stores.GET("/", h.GetAllStores)
		stores.POST("/", h.UserIdentify, h.CheckRole("admin"), h.CreateStore)
		stores.PATCH("/:id", h.UserIdentify, h.CheckRole("admin"), h.ChangeStore)
		stores.DELETE("/:id", h.UserIdentify, h.CheckRole("admin"), h.DeleteStore)
	}

	return router
}
