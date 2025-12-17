package handlers

import (
	"encoding/json"
	"errors"
	"eventflow/internal/database"
	"eventflow/internal/models"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetEventRegistrationById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(400, gin.H{"error": "ID parameter is required"})
		return
	}

	var eventRegistration models.EventRegistration

	result := database.DB.First(&eventRegistration, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "EventRegistration not found"})
		} else {
			c.JSON(500, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(200, eventRegistration)
}

// @Summary Получить список регистраций на события
// @Description Возвращает список всех регистраций участников на события с пагинацией
// @Tags EventRegistrations
// @Accept json
// @Produce json
// @Param range query string false "Пагинация [start, end]"
// @Param sort query string false "Сортировка [field, order]"
// @Success 200 {array} models.EventRegistration
// @Header 200 {string} X-Total-Count "Общее количество записей"
// @Header 200 {string} Content-Range "Диапазон записей"
// @Router /event_registrations [get]
func GetEventRegistrations(c *gin.Context) {
	var eventRegistrations []models.EventRegistration
	var total int64

	rangeParam := c.Query("range")
	var start, end int = 0, 25
	if rangeParam != "" {
		var rangeArray []int
		if err := json.Unmarshal([]byte(rangeParam), &rangeArray); err == nil && len(rangeArray) == 2 {
			start = rangeArray[0]
			end = rangeArray[1]
		}
	}

	sortParam := c.Query("sort")
	var sortField, sortOrder string = "id", "ASC"
	if sortParam != "" {
		var sortArray []string
		if err := json.Unmarshal([]byte(sortParam), &sortArray); err == nil && len(sortArray) == 2 {
			sortField = sortArray[0]
			sortOrder = sortArray[1]
		}
	}

	limit := end - start + 1
	offset := start

	countResult := database.DB.Model(&models.EventRegistration{}).Count(&total)
	if countResult.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve total record count"})
		return
	}

	contentRange := fmt.Sprintf("event_registrations %d-%d/%d", start, end, total)

	result := database.DB.
		Limit(limit).
		Offset(offset).
		Order(sortField + " " + sortOrder).
		Find(&eventRegistrations)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	c.Header("Content-Range", contentRange)
	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(200, eventRegistrations)
}

func UpdateEventRegistration(c *gin.Context) {
	id := c.Param("id")

	var input models.CreateEventRegistrationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var eventRegistration models.EventRegistration

	result := database.DB.Model(&eventRegistration).Where("id = ?", id).Updates(models.EventRegistration{
		EventID:       input.EventID,
		ParticipantID: input.ParticipantID,
		Status:        input.Status,
	})

	if result.Error != nil {
		log.Printf("Database Error (Update): %v", result.Error)
		c.JSON(500, gin.H{"error": "Failed to update event registration. Database error."})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Event registration not found."})
		return
	}

	database.DB.First(&eventRegistration, id)

	c.JSON(200, eventRegistration)
}

func DeleteEventRegistration(c *gin.Context) {
	id := c.Param("id")

	result := database.DB.Delete(&models.EventRegistration{}, id)

	if result.Error != nil {
		log.Printf("Database Error (Delete): %v", result.Error)
		c.JSON(500, gin.H{"error": "Failed to delete event registration. Database error."})
		return
	}

	c.JSON(200, gin.H{})
}

// @Summary Зарегистрировать участника на событие
// @Description Создает регистрацию участника на событие и автоматически создает тикет
// @Tags EventRegistrations
// @Accept json
// @Produce json
// @Param registration body models.CreateEventRegistrationRequest true "Данные регистрации"
// @Success 201 {object} models.EventRegistration
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /event_registrations [post]
func PostEventRegistration(c *gin.Context) {
	var newEventRegistration models.CreateEventRegistrationRequest

	if err := c.ShouldBindJSON(&newEventRegistration); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	eventRegistration := models.EventRegistration{
		EventID:       newEventRegistration.EventID,
		ParticipantID: newEventRegistration.ParticipantID,
		Status:        newEventRegistration.Status,
		RegisteredAt:  time.Now(),
	}

	result := database.DB.Create(&eventRegistration)
	if result.Error != nil {
		log.Printf("Database Error (Create): %v", result.Error)
		c.JSON(500, gin.H{"error": "Failed to create event registration. Database error."})
		return
	}

	qrCode, err := generateQRCode()
	if err != nil {
		log.Printf("QR Code Generation Error: %v", err)
	} else {
		ticket := models.Ticket{
			EventID:       newEventRegistration.EventID,
			ParticipantID: newEventRegistration.ParticipantID,
			TicketType:    "free", // По умолчанию бесплатный
			Status:        "active",
			QRCode:        qrCode,
		}

		ticketResult := database.DB.Create(&ticket)
		if ticketResult.Error != nil {
			log.Printf("Auto-ticket creation failed: %v", ticketResult.Error)
		}
	}

	c.JSON(201, eventRegistration)
}
