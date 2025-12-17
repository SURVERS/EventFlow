package main

import (
	"eventflow/internal/database"
	"eventflow/internal/handlers"
	"eventflow/internal/middleware"
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "eventflow/docs" // Swagger docs
)

// @title EventFlow API
// @version 1.0
// @description API для системы управления событиями EventFlow
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@eventflow.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	database.Connect()

	router := gin.Default()

	config := cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Total-Count", "Range", "Content-Range", "Accept"},
		ExposeHeaders:    []string{"X-Total-Count", "Content-Range"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}

	router.Use(cors.New(config))
	router.Use(middleware.LoggerMiddleware())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		v1.POST("/auth/register", handlers.Register)
		v1.POST("/auth/login", handlers.Login)
		v1.POST("/auth/refresh", handlers.RefreshAccessToken)
		v1.GET("/auth/me", middleware.AuthMiddleware(), handlers.GetCurrentUser)

		v1.GET("/categories", handlers.GetCategories)
		v1.POST("/categories", handlers.PostCategory)
		v1.PUT("/categories/:id", handlers.UpdateCategory)
		v1.DELETE("/categories/:id", handlers.DeleteCategory)
		v1.GET("/categories/:id", handlers.GetCategoryById)

		v1.GET("/events", handlers.GetEvents)
		v1.POST("/events", handlers.PostEvent)
		v1.PUT("/events/:id", handlers.UpdateEvent)
		v1.DELETE("/events/:id", handlers.DeleteEvent)
		v1.GET("/events/:id", handlers.GetEventById)

		v1.GET("/event_types", handlers.GetEventTypes)
		v1.POST("/event_types", handlers.PostEventType)
		v1.PUT("/event_types/:id", handlers.UpdateEventType)
		v1.DELETE("/event_types/:id", handlers.DeleteEventType)
		v1.GET("/event_types/:id", handlers.GetEventTypeById)

		v1.GET("/participants", handlers.GetParticipants)
		v1.POST("/participants", handlers.PostParticipant)
		v1.PUT("/participants/:id", handlers.UpdateParticipant)
		v1.DELETE("/participants/:id", handlers.DeleteParticipant)
		v1.GET("/participants/:id/statistics", handlers.GetParticipantStatistics)
		v1.GET("/participants/:id", handlers.GetParticipantById)

		v1.GET("/event_registrations", handlers.GetEventRegistrations)
		v1.POST("/event_registrations", handlers.PostEventRegistration)
		v1.PUT("/event_registrations/:id", handlers.UpdateEventRegistration)
		v1.DELETE("/event_registrations/:id", handlers.DeleteEventRegistration)
		v1.GET("/event_registrations/:id", handlers.GetEventRegistrationById)

		v1.GET("/organizers", handlers.GetOrganizers)
		v1.POST("/organizers", handlers.PostOrganizer)
		v1.PUT("/organizers/:id", handlers.UpdateOrganizer)
		v1.DELETE("/organizers/:id", handlers.DeleteOrganizer)
		v1.GET("/organizers/:id", handlers.GetOrganizerById)

		v1.GET("/tickets", handlers.GetTickets)
		v1.POST("/tickets", handlers.PostTicket)
		v1.PUT("/tickets/:id", handlers.UpdateTicket)
		v1.DELETE("/tickets/:id", handlers.DeleteTicket)
		v1.GET("/tickets/qr/:qrcode", handlers.GetTicketByQRCode)
		v1.POST("/tickets/qr/:qrcode/use", handlers.MarkTicketAsUsed)
		v1.GET("/tickets/:id", handlers.GetTicketById)

		v1.GET("/dashboard/statistics", handlers.GetDashboardStatistics)
		v1.GET("/dashboard/popular-categories", handlers.GetPopularCategories)
		v1.GET("/dashboard/events/:id/statistics", handlers.GetEventStatistics)
	}

	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server started: https://localhost:8080/")
}
