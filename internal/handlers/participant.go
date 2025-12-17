package handlers

import (
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

func GetParticipantById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(400, gin.H{"error": "ID parameter is required"})
		return
	}

	var participant models.Participant

	result := database.DB.First(&participant, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Participant not found"})
		} else {
			c.JSON(500, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(200, participant)
}

// @Summary Получить список участников
// @Description Возвращает список всех участников событий с пагинацией
// @Tags Participants
// @Accept json
// @Produce json
// @Param range query string false "Пагинация [start, end]"
// @Param sort query string false "Сортировка [field, order]"
// @Success 200 {array} models.Participant
// @Header 200 {string} X-Total-Count "Общее количество записей"
// @Header 200 {string} Content-Range "Диапазон записей"
// @Router /participants [get]
func GetParticipants(c *gin.Context) {
	var participants []models.Participant
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

	countResult := database.DB.Model(&models.Participant{}).Count(&total)
	if countResult.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve total record count"})
		return
	}

	contentRange := fmt.Sprintf("participants %d-%d/%d", start, end, total)

	result := database.DB.
		Limit(limit).
		Offset(offset).
		Order(sortField + " " + sortOrder).
		Find(&participants)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	c.Header("Content-Range", contentRange)
	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(200, participants)
}

func UpdateParticipant(c *gin.Context) {
	id := c.Param("id")

	var input models.CreateParticipantRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var participant models.Participant

	result := database.DB.Model(&participant).Where("id = ?", id).Updates(models.Participant{FullName: input.FullName, Email: input.Email, Phone: input.Phone})

	if result.Error != nil {
		log.Printf("Database Error (Update): %v", result.Error)
		c.JSON(500, gin.H{"error": "Failed to update participant. Database error."})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Participant not found."})
		return
	}

	database.DB.First(&participant, id)

	c.JSON(200, participant)
}

func DeleteParticipant(c *gin.Context) {
	id := c.Param("id")

	result := database.DB.Delete(&models.Participant{}, id)

	if result.Error != nil {
		log.Printf("Database Error (Delete): %v", result.Error)
		c.JSON(500, gin.H{"error": "Failed to delete participant. Database error."})
		return
	}

	c.JSON(200, gin.H{})
}

// @Summary Создать участника
// @Description Регистрирует нового участника событий
// @Tags Participants
// @Accept json
// @Produce json
// @Param participant body models.CreateParticipantRequest true "Данные участника"
// @Success 201 {object} models.Participant
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /participants [post]
func PostParticipant(c *gin.Context) {
	var newParticipant models.CreateParticipantRequest

	if err := c.ShouldBindJSON(&newParticipant); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	participant := models.Participant{
		FullName: newParticipant.FullName,
		Email:    newParticipant.Email,
		Phone:    newParticipant.Phone,
	}

	result := database.DB.Create(&participant)

	if result.Error != nil {
		log.Printf("Database Error (Create): %v", result.Error)
		c.JSON(500, gin.H{"error": "Failed to create participant. Database error."})
		return
	}

	c.JSON(201, participant)
}

func GetParticipantStatistics(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(400, gin.H{"error": "ID parameter is required"})
		return
	}

	var participant models.Participant
	result := database.DB.First(&participant, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Participant not found"})
		} else {
			c.JSON(500, gin.H{"error": "Database error"})
		}
		return
	}

	var totalRegistrations int64
	var registeredCount int64
	var attendedCount int64
	var noShowCount int64

	database.DB.Model(&models.EventRegistration{}).Where("participant_id = ?", id).Count(&totalRegistrations)
	database.DB.Model(&models.EventRegistration{}).Where("participant_id = ? AND status = ?", id, "registered").Count(&registeredCount)
	database.DB.Model(&models.EventRegistration{}).Where("participant_id = ? AND status = ?", id, "attended").Count(&attendedCount)
	database.DB.Model(&models.EventRegistration{}).Where("participant_id = ? AND status = ?", id, "no-show").Count(&noShowCount)

	var attendanceRate float64 = 0
	if totalRegistrations > 0 {
		attendanceRate = (float64(attendedCount) / float64(totalRegistrations)) * 100
	}

	statistics := gin.H{
		"participant_id":      participant.ID,
		"full_name":           participant.FullName,
		"total_registrations": totalRegistrations,
		"registered_count":    registeredCount,
		"attended_count":      attendedCount,
		"no_show_count":       noShowCount,
		"attendance_rate":     attendanceRate,
	}

	c.JSON(200, statistics)
}
