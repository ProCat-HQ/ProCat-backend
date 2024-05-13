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

	router.StaticFS("/assets", gin.Dir("./assets", false))

	users := router.Group("/users")
	{
		usersAuthenticatedGroup := users.Group("", h.UserIdentify)
		{
			usersAuthenticatedGroup.GET("", h.CheckRole("admin"), h.GetAllUsers)
			usersAuthenticatedGroup.GET("/:id", h.MustBelongsToUser, h.GetUser)
			usersAuthenticatedGroup.DELETE("/:id", h.CheckRole("admin"), h.DeleteUser) // TODO

			users.POST("/logout", h.Logout)
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
			change.POST("/iin-bin", h.ChangeIIN)                          // TODO
			change.POST("/fullname", h.ChangeFullName)                    // TODO
			change.POST("/password", h.ChangePassword)                    // TODO
			change.POST("/email", h.ChangeEmail)                          // TODO
			change.POST("/phone", h.ChangePhone)                          // TODO
			change.PATCH("/role/:id", h.CheckRole("admin"), h.ChangeRole) // TODO
		}

		deliverymen := users.Group("/deliverymen", h.UserIdentify, h.CheckRole("deliveryman"))
		{
			deliverymen.GET("", h.GetAllDeliverymen)                                 // TODO
			deliverymen.GET("/:id", h.GetDeliveryman)                                // TODO
			deliverymen.POST("/:id", h.CheckRole("admin"), h.CreateDeliveryman)      // TODO
			deliverymen.PATCH("/:id", h.CheckRole("admin"), h.ChangeDeliverymanData) // TODO
			deliverymen.DELETE("/:id", h.CheckRole("admin"), h.DeleteDeliveryman)    // TODO

			deliveries := deliverymen.Group("/deliveries")
			{
				deliveries.GET("", h.CheckRole("admin"), h.GetAllDeliveries) // TODO
				deliveries.GET("/:id", h.GetAllDeliveriesForOneDeliveryman)  // TODO
				deliveries.PATCH("/:id", h.ChangeDeliveryStatus)             // TODO
				deliveries.POST("/create-route", h.CreateRoute)              // TODO
			}
		}

		admin := users.Group("/admin", h.UserIdentify, h.CheckRole("admin"))
		{
			admin.POST("/cluster", h.Cluster)                          // TODO
			admin.GET("/deliveries-to-sort", h.GetAllDeliveriesToSort) // TODO
			admin.PATCH("/change-delivery", h.ChangeDeliveryData)      // TODO
		}

		cart := users.Group("/cart", h.UserIdentify)
		{
			cart.GET("", h.GetCartItems)           // TODO
			cart.POST("", h.AddItemsToCart)        // TODO
			cart.DELETE("", h.DeleteItemsFromCart) // TODO
		}

		orders := users.Group("/orders", h.UserIdentify)
		{
			orders.GET("", h.GetAllOrders)                                         // TODO
			orders.GET("/:id", h.GetOrder)                                         // TODO
			orders.POST("", h.CreateOrder)                                         // TODO
			orders.POST("/cancel/:id", h.CancelOrder)                              // TODO
			orders.PATCH("/status/:id", h.CheckRole("admin"), h.ChangeOrderStatus) // TODO

			payment := orders.Group("/payment")
			{
				payment.GET("/:id", h.GetPaymentData)                              // TODO
				payment.PATCH("/:id", h.CheckRole("admin"), h.ChangePaymentStatus) // TODO
			}
		}

		subscriptions := users.Group("/subscriptions", h.UserIdentify)
		{
			subscriptions.GET("", h.GetAllSubscriptions)                // TODO
			subscriptions.POST("", h.SubscribeToItem)                   // TODO
			subscriptions.DELETE("/:id", h.DeleteItemFromSubscriptions) // TODO
		}

		notifications := users.Group("/notifications", h.UserIdentify)
		{
			notifications.GET("", h.GetAllNotifications)                             // TODO
			notifications.POST("/:id", h.CheckRole("admin"), h.SendNotification)     // TODO
			notifications.PATCH("/:id", h.ViewNotification)                          // TODO
			notifications.DELETE("/:id", h.CheckRole("admin"), h.DeleteNotification) // TODO
		}

	}

	categories := router.Group("/categories")
	{
		categories.GET("/:id", h.GetCategory)                                             // TODO
		categories.GET("/route/:id", h.GetCategoryRoute)                                  // TODO
		categories.POST("/:id", h.UserIdentify, h.CheckRole("admin"), h.CreateCategory)   // TODO
		categories.PATCH("/:id", h.UserIdentify, h.CheckRole("admin"), h.ChangeCategory)  // TODO
		categories.DELETE("/:id", h.UserIdentify, h.CheckRole("admin"), h.DeleteCategory) // TODO
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
			stock.PUT("/:id", h.ChangeStock) // TODO
		}

		infos := items.Group("/infos", h.UserIdentify, h.CheckRole("admin"))
		{
			infos.POST("/:id", h.AddInfo)      // TODO
			infos.PATCH("/:id", h.ChangeInfo)  // TODO
			infos.DELETE("/:id", h.DeleteInfo) // TODO
		}

		images := items.Group("/images", h.UserIdentify, h.CheckRole("admin"))
		{
			images.POST("/:id", h.AddImages)      // TODO
			images.DELETE("/:id", h.DeleteImages) // TODO
		}
	}

	stores := router.Group("/stores")
	{
		stores.GET("", h.GetAllStores)                                             // TODO
		stores.POST("", h.UserIdentify, h.CheckRole("admin"), h.CreateStore)       // TODO
		stores.PATCH("/:id", h.UserIdentify, h.CheckRole("admin"), h.ChangeStore)  // TODO
		stores.DELETE("/:id", h.UserIdentify, h.CheckRole("admin"), h.DeleteStore) // TODO
	}

	return router
}
