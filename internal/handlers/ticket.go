package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"eventflow/internal/database"
	"eventflow/internal/models"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func generateQRCode() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GetTicketById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(400, gin.H{"error": "ID parameter is required"})
		return
	}

	var ticket models.Ticket

	result := database.DB.First(&ticket, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Ticket not found"})
		} else {
			c.JSON(500, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(200, ticket)
}

// @Summary Получить список тикетов
// @Description Возвращает список всех тикетов с пагинацией
// @Tags Tickets
// @Accept json
// @Produce json
// @Param range query string false "Пагинация [start, end]"
// @Param sort query string false "Сортировка [field, order]"
// @Success 200 {array} models.Ticket
// @Header 200 {string} X-Total-Count "Общее количество записей"
// @Header 200 {string} Content-Range "Диапазон записей"
// @Router /tickets [get]
func GetTickets(c *gin.Context) {
	var tickets []models.Ticket
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

	countResult := database.DB.Model(&models.Ticket{}).Count(&total)
	if countResult.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve total record count"})
		return
	}

	contentRange := fmt.Sprintf("tickets %d-%d/%d", start, end, total)

	result := database.DB.
		Limit(limit).
		Offset(offset).
		Order(sortField + " " + sortOrder).
		Find(&tickets)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	c.Header("Content-Range", contentRange)
	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(200, tickets)
}

func UpdateTicket(c *gin.Context) {
	id := c.Param("id")

	var input models.UpdateTicketRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if input.TicketType != "free" && input.TicketType != "paid" {
		c.JSON(400, gin.H{"error": "Ticket type must be either 'free' or 'paid'"})
		return
	}

	if input.Status != "active" && input.Status != "canceled" {
		c.JSON(400, gin.H{"error": "Status must be either 'active' or 'canceled'"})
		return
	}

	var ticket models.Ticket

	result := database.DB.Model(&ticket).Where("id = ?", id).Updates(models.Ticket{
		TicketType: input.TicketType,
		Status:     input.Status,
	})

	if result.Error != nil {
		log.Printf("Database Error (Update): %v", result.Error)
		c.JSON(500, gin.H{"error": "Failed to update ticket. Database error."})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Ticket not found."})
		return
	}

	database.DB.First(&ticket, id)

	c.JSON(200, ticket)
}

func DeleteTicket(c *gin.Context) {
	id := c.Param("id")

	result := database.DB.Delete(&models.Ticket{}, id)

	if result.Error != nil {
		log.Printf("Database Error (Delete): %v", result.Error)
		c.JSON(500, gin.H{"error": "Failed to delete ticket. Database error."})
		return
	}

	c.JSON(200, gin.H{})
}

// @Summary Создать тикет
// @Description Создает новый тикет с уникальным QR-кодом
// @Tags Tickets
// @Accept json
// @Produce json
// @Param ticket body models.CreateTicketRequest true "Данные тикета"
// @Success 201 {object} models.Ticket
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tickets [post]
func PostTicket(c *gin.Context) {
	var newTicket models.CreateTicketRequest

	if err := c.ShouldBindJSON(&newTicket); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if newTicket.TicketType != "free" && newTicket.TicketType != "paid" {
		c.JSON(400, gin.H{"error": "Ticket type must be either 'free' or 'paid'"})
		return
	}

	if newTicket.Status != "active" && newTicket.Status != "canceled" {
		c.JSON(400, gin.H{"error": "Status must be either 'active' or 'canceled'"})
		return
	}

	qrCode, err := generateQRCode()
	if err != nil {
		log.Printf("QR Code Generation Error: %v", err)
		c.JSON(500, gin.H{"error": "Failed to generate QR code"})
		return
	}

	ticket := models.Ticket{
		EventID:       newTicket.EventID,
		ParticipantID: newTicket.ParticipantID,
		TicketType:    newTicket.TicketType,
		Status:        newTicket.Status,
		QRCode:        qrCode,
	}

	result := database.DB.Create(&ticket)

	if result.Error != nil {
		log.Printf("Database Error (Create): %v", result.Error)
		c.JSON(500, gin.H{"error": "Failed to create ticket. Database error."})
		return
	}

	c.JSON(201, ticket)
}

// @Summary Получить тикет по QR-коду
// @Description Получение тикета по QR-коду для отслеживания посещаемости
// @Tags Tickets
// @Accept json
// @Produce json
// @Param qrcode path string true "QR-код тикета"
// @Success 200 {object} models.Ticket
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tickets/qr/{qrcode} [get]
func GetTicketByQRCode(c *gin.Context) {
	qrCode := c.Param("qrcode")

	if qrCode == "" {
		c.JSON(400, gin.H{"error": "QR code parameter is required"})
		return
	}

	var ticket models.Ticket

	result := database.DB.Where("qr_code = ?", qrCode).First(&ticket)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Ticket not found"})
		} else {
			c.JSON(500, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(200, ticket)
}

// @Summary Отметить тикет как использованный
// @Description Отметить тикет как использованный для отслеживания посещаемости
// @Tags Tickets
// @Accept json
// @Produce json
// @Param qrcode path string true "QR-код тикета"
// @Success 200 {object} models.Ticket
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tickets/qr/{qrcode}/use [post]
func MarkTicketAsUsed(c *gin.Context) {
	qrCode := c.Param("qrcode")

	if qrCode == "" {
		c.JSON(400, gin.H{"error": "QR code parameter is required"})
		return
	}

	var ticket models.Ticket

	result := database.DB.Where("qr_code = ?", qrCode).First(&ticket)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Ticket not found"})
		} else {
			c.JSON(500, gin.H{"error": "Database error"})
		}
		return
	}

	if ticket.Status == "canceled" {
		c.JSON(400, gin.H{"error": "Ticket is canceled and cannot be used"})
		return
	}

	var registration models.EventRegistration
	regResult := database.DB.Where("event_id = ? AND participant_id = ?", ticket.EventID, ticket.ParticipantID).First(&registration)

	if regResult.Error == nil {
		database.DB.Model(&registration).Update("status", "attended")
	}

	c.JSON(200, gin.H{
		"message":        "Ticket successfully used. Attendance recorded.",
		"ticket":         ticket,
		"attendance_marked": regResult.Error == nil,
	})
}
