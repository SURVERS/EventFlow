package handlers

import (
	"eventflow/internal/database"
	"eventflow/internal/models"

	"github.com/gin-gonic/gin"
)

// @Summary Общая статистика для dashboard
// @Description Возвращает общую статистику по событиям, участникам, регистрациям и посещаемости
// @Tags Dashboard
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Статистика dashboard"
// @Router /dashboard/statistics [get]
func GetDashboardStatistics(c *gin.Context) {
	var totalEvents, totalParticipants, totalOrganizers, totalTickets int64
	var publishedEvents, draftEvents int64

	database.DB.Model(&models.Event{}).Count(&totalEvents)
	database.DB.Model(&models.Event{}).Where("publish_status = ?", "published").Count(&publishedEvents)
	database.DB.Model(&models.Event{}).Where("publish_status = ?", "draft").Count(&draftEvents)

	database.DB.Model(&models.Participant{}).Count(&totalParticipants)

	database.DB.Model(&models.Organizer{}).Count(&totalOrganizers)

	database.DB.Model(&models.Ticket{}).Count(&totalTickets)

	var registeredCount, attendedCount, noShowCount int64
	database.DB.Model(&models.EventRegistration{}).Where("status = ?", "registered").Count(&registeredCount)
	database.DB.Model(&models.EventRegistration{}).Where("status = ?", "attended").Count(&attendedCount)
	database.DB.Model(&models.EventRegistration{}).Where("status = ?", "no-show").Count(&noShowCount)

	totalRegistrations := registeredCount + attendedCount + noShowCount
	var attendanceRate float64 = 0
	if totalRegistrations > 0 {
		attendanceRate = (float64(attendedCount) / float64(totalRegistrations)) * 100
	}

	statistics := gin.H{
		"total_events":       totalEvents,
		"published_events":   publishedEvents,
		"draft_events":       draftEvents,
		"total_participants": totalParticipants,
		"total_organizers":   totalOrganizers,
		"total_tickets":      totalTickets,
		"total_registrations": totalRegistrations,
		"registered_count":   registeredCount,
		"attended_count":     attendedCount,
		"no_show_count":      noShowCount,
		"attendance_rate":    attendanceRate,
	}

	c.JSON(200, statistics)
}

func GetPopularCategories(c *gin.Context) {
	type CategoryStat struct {
		CategoryID   uint   `json:"category_id"`
		CategoryName string `json:"category_name"`
		EventCount   int64  `json:"event_count"`
	}

	var results []CategoryStat

	database.DB.Table("events").
		Select("events.category_id, categories.name as category_name, COUNT(events.id) as event_count").
		Joins("LEFT JOIN categories ON categories.id = events.category_id").
		Group("events.category_id, categories.name").
		Order("event_count DESC").
		Limit(10).
		Scan(&results)

	c.JSON(200, results)
}

func GetEventStatistics(c *gin.Context) {
	eventID := c.Param("id")

	var totalTickets, activeTickets, canceledTickets int64
	var totalRegistrations, attendedCount int64

	database.DB.Model(&models.Ticket{}).Where("event_id = ?", eventID).Count(&totalTickets)
	database.DB.Model(&models.Ticket{}).Where("event_id = ? AND status = ?", eventID, "active").Count(&activeTickets)
	database.DB.Model(&models.Ticket{}).Where("event_id = ? AND status = ?", eventID, "canceled").Count(&canceledTickets)

	database.DB.Model(&models.EventRegistration{}).Where("event_id = ?", eventID).Count(&totalRegistrations)
	database.DB.Model(&models.EventRegistration{}).Where("event_id = ? AND status = ?", eventID, "attended").Count(&attendedCount)

	var attendanceRate float64 = 0
	if totalRegistrations > 0 {
		attendanceRate = (float64(attendedCount) / float64(totalRegistrations)) * 100
	}

	statistics := gin.H{
		"event_id":           eventID,
		"total_tickets":      totalTickets,
		"active_tickets":     activeTickets,
		"canceled_tickets":   canceledTickets,
		"total_registrations": totalRegistrations,
		"attended_count":     attendedCount,
		"attendance_rate":    attendanceRate,
	}

	c.JSON(200, statistics)
}
