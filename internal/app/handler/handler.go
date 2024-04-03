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
		users.GET("/", h.GetAllUsers)
		users.GET("/:id", h.GetUser)
		users.POST("/sign-in", h.SignIn)
		users.POST("/sign-up", h.SignUp)

		verification := users.Group("/verification")
		{
			verification.PUT("/send", h.SendCode)
			verification.POST("/check", h.CheckCode)
			verification.POST("/iin-bin", h.CheckIIN)
		}

		change := users.Group("/change")
		{
			change.POST("/iin-bin", h.ChangeIIN)
			change.POST("/fullname", h.ChangeFullName)
			change.POST("/password", h.ChangePassword)
			change.POST("/email", h.ChangeEmail)
			change.POST("/phone", h.ChangePhone)
			change.PATCH("/role/:id", h.ChangeRole)
		}

		deliverymen := users.Group("/deliverymen")
		{
			deliverymen.GET("/", h.GetAllDeliverymen)
			deliverymen.GET("/:id", h.GetDeliveryman)
			deliverymen.POST("/:id", h.CreateDeliveryman)
			deliverymen.PATCH("/:id", h.ChangeDeliverymanData)
			deliverymen.DELETE("/:id", h.DeleteDeliveryman)

			deliveries := deliverymen.Group("/deliveries")
			{
				deliveries.GET("/", h.GetAllDeliveries)
				deliveries.GET("/:id", h.GetDelivery)
				deliveries.PATCH("/:id", h.ChangeDeliveryStatus)
				deliveries.POST("/create-route", h.CreateRoute)
			}
		}

		admin := users.Group("/admin")
		{
			admin.GET("/deliveries-to-sort", h.GetAllDeliveriesToSort)
			admin.PATCH("/change-delivery", h.ChangeDeliveryData)
			admin.POST("/cluster", h.Cluster)
		}

		cart := users.Group("/cart")
		{
			cart.GET("/", h.GetCartItems)
			cart.POST("/", h.AddItemsToCart)
			cart.DELETE("/", h.DeleteItemsFromCart)
		}

		orders := users.Group("/orders")
		{
			orders.GET("/", h.GetAllOrders)
			orders.GET("/:id", h.GetOrder)
			orders.POST("/", h.CreateOrder)
			orders.POST("/cancel/:id", h.CancelOrder)
			orders.PATCH("/status/:id", h.ChangeOrderStatus)

			payment := orders.Group("/payment")
			{
				payment.GET("/:id", h.GetPaymentData)
				payment.PATCH("/:id", h.ChangePaymentStatus)
			}
		}

		subscriptions := users.Group("/subscriptions")
		{
			subscriptions.GET("/", h.GetAllSubscriptions)
			subscriptions.POST("/", h.SubscribeToItem)
			subscriptions.DELETE("/:id", h.DeleteItemFromSubscriptions)
		}

		notifications := users.Group("/notifications")
		{
			notifications.GET("/", h.GetAllNotifications)
			notifications.POST("/:id", h.SendNotification)
			notifications.PATCH("/:id", h.ViewNotification)
			notifications.DELETE("/:id", h.DeleteNotification)
		}

	}

	categories := router.Group("/categories")
	{
		categories.GET("/:id", h.GetCategory)
		categories.GET("/route/:id", h.GetCategoryRoute)
		categories.POST("/:id", h.CreateCategory)
		categories.PATCH("/:id", h.ChangeCategory)
		categories.DELETE("/:id", h.DeleteCategory)
	}

	items := router.Group("/items")
	{
		items.GET("/", h.GetAllItems)
		items.GET("/:id", h.GetItem)
		items.POST("/", h.CreateItem)
		items.PATCH("/:id", h.ChangeItem)
		items.DELETE("/:id", h.DeleteItem)

		stock := items.Group("/stock")
		{
			stock.POST("/:id", h.CreateStock)
			stock.PATCH("/:id", h.ChangeStock)
		}

		infos := items.Group("/infos")
		{
			infos.POST("/:id", h.AddInfo)
			infos.PUT("/:id", h.ChangeInfo)
			infos.DELETE("/:id", h.DeleteInfo)
		}

		images := items.Group("/images")
		{
			images.POST("/:id", h.AddImages)
			images.PATCH("/:id", h.ChangeImages)
			images.DELETE("/:id", h.DeleteImages)
		}
	}

	stores := router.Group("/stores")
	{
		stores.GET("/", h.GetAllStores)
		stores.POST("/", h.CreateStore)
		stores.PATCH("/:id", h.ChangeStore)
		stores.DELETE("/:id", h.DeleteStore)
	}

	return router
}
