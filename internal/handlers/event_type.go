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

func GetEventTypeById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(400, gin.H{"error": "ID parameter is required"})
		return
	}

	var eventType models.EventType

	result := database.DB.First(&eventType, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Event type not found"})
		} else {
			c.JSON(500, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(200, eventType)
}

func GetEventTypes(c *gin.Context) {
	var eventTypes []models.EventType
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

	countResult := database.DB.Model(&models.EventType{}).Count(&total)
	if countResult.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve total record count"})
		return
	}

	contentRange := fmt.Sprintf("event_types %d-%d/%d", start, end, total)

	result := database.DB.
		Limit(limit).
		Offset(offset).
		Order(sortField + " " + sortOrder).
		Find(&eventTypes)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	c.Header("Content-Range", contentRange)
	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(200, eventTypes)
}

func UpdateEventType(c *gin.Context) {
	id := c.Param("id")

	var input models.CreateEventTypeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var eventType models.EventType

	result := database.DB.Model(&eventType).Where("id = ?", id).Updates(models.EventType{Name: input.Name})

	if result.Error != nil {
		log.Printf("Database Error (Update): %v", result.Error)
		c.JSON(500, gin.H{"error": "Failed to update event_type. Database error."})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Event type not found."})
		return
	}

	database.DB.First(&eventType, id)

	c.JSON(200, eventType)
}

func DeleteEventType(c *gin.Context) {
	id := c.Param("id")

	result := database.DB.Delete(&models.EventType{}, id)

	if result.Error != nil {
		log.Printf("Database Error (Delete): %v", result.Error)
		c.JSON(500, gin.H{"error": "Failed to delete event type. Database error."})
		return
	}

	c.JSON(200, gin.H{})
}

func PostEventType(c *gin.Context) {
	var newPostEventType models.CreateEventTypeRequest

	if err := c.ShouldBindJSON(&newPostEventType); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	eventType := models.EventType{
		Name: newPostEventType.Name,
	}

	result := database.DB.Create(&eventType)
	if result.Error != nil {
		log.Printf("Database Error (Create): %v", result.Error)
		c.JSON(500, gin.H{"error": "Failed to create event type. Database error."})
		return
	}

	c.JSON(201, eventType)
}
