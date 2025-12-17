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

func GetEventById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(400, gin.H{"error": "ID parameter is required"})
		return
	}

	var event models.Event

	result := database.DB.First(&event, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Event not found"})
		} else {
			c.JSON(500, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(200, event)
}

// @Summary Получить список событий
// @Description Возвращает список событий с поддержкой пагинации и фильтрации
// @Tags Events
// @Accept json
// @Produce json
// @Param range query string false "Пагинация [start, end]" example([0, 24])
// @Param sort query string false "Сортировка [field, order]" example(["id", "ASC"])
// @Param filter query string false "Фильтрация {field: value}" example({"category_id": 1})
// @Param category_id query int false "Фильтр по категории"
// @Param start_date query string false "Фильтр по дате начала (>= start_date)"
// @Param end_date query string false "Фильтр по дате окончания (<= end_date)"
// @Param publish_status query string false "Фильтр по статусу публикации (published/draft)"
// @Success 200 {array} models.Event "Список событий"
// @Header 200 {string} X-Total-Count "Общее количество записей"
// @Header 200 {string} Content-Range "Диапазон записей"
// @Router /events [get]
func GetEvents(c *gin.Context) {
	var events []models.Event
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

	query := database.DB.Model(&models.Event{})

	categoryID := c.Query("category_id")
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	startDate := c.Query("start_date")
	if startDate != "" {
		query = query.Where("start_time >= ?", startDate)
	}

	endDate := c.Query("end_date")
	if endDate != "" {
		query = query.Where("end_time <= ?", endDate)
	}

	publishStatus := c.Query("publish_status")
	if publishStatus != "" {
		query = query.Where("publish_status = ?", publishStatus)
	}

	countResult := query.Count(&total)
	if countResult.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve total record count"})
		return
	}

	contentRange := fmt.Sprintf("events %d-%d/%d", start, end, total)

	result := query.
		Limit(limit).
		Offset(offset).
		Order(sortField + " " + sortOrder).
		Find(&events)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	c.Header("Content-Range", contentRange)
	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(200, events)
}

func UpdateEvent(c *gin.Context) {
	id := c.Param("id")

	var input models.CreateEventRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if input.PublishStatus != "draft" && input.PublishStatus != "published" {
		c.JSON(400, gin.H{"error": "Publish status must be either 'draft' or 'published'"})
		return
	}

	var event models.Event

	result := database.DB.Model(&event).Where("id = ?", id).Updates(models.Event{
		Title:         input.Title,
		Description:   input.Description,
		StartTime:     input.StartTime,
		EndTime:       input.EndTime,
		EventType:     input.EventType,
		Status:        input.Status,
		PublishStatus: input.PublishStatus,
		CategoryID:    input.CategoryID,
	})

	if result.Error != nil {
		log.Printf("Database Error (Update): %v", result.Error)
		c.JSON(500, gin.H{"error": "Failed to update Event. Database error."})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Event not found."})
		return
	}

	database.DB.First(&event, id)

	c.JSON(200, event)
}

func DeleteEvent(c *gin.Context) {
	id := c.Param("id")

	result := database.DB.Delete(&models.Event{}, id)

	if result.Error != nil {
		log.Printf("Database Error (Delete): %v", result.Error)
		c.JSON(500, gin.H{"error": "Failed to delete event. Database error."})
		return
	}

	c.JSON(200, gin.H{})
}

// @Summary Создать новое событие
// @Description Создает новое событие в системе
// @Tags Events
// @Accept json
// @Produce json
// @Param event body models.CreateEventRequest true "Данные события"
// @Success 201 {object} models.Event "Созданное событие"
// @Failure 400 {object} map[string]string "Ошибка валидации"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /events [post]
func PostEvent(c *gin.Context) {
	var newPostEvent models.CreateEventRequest

	if err := c.ShouldBindJSON(&newPostEvent); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if newPostEvent.PublishStatus != "draft" && newPostEvent.PublishStatus != "published" {
		c.JSON(400, gin.H{"error": "Publish status must be either 'draft' or 'published'"})
		return
	}

	status := newPostEvent.Status
	if status == "" {
		status = "scheduled"
	}

	Event := models.Event{
		Title:         newPostEvent.Title,
		Description:   newPostEvent.Description,
		StartTime:     newPostEvent.StartTime,
		EndTime:       newPostEvent.EndTime,
		EventType:     newPostEvent.EventType,
		Status:        status,
		PublishStatus: newPostEvent.PublishStatus,
		CategoryID:    newPostEvent.CategoryID,
	}

	result := database.DB.Create(&Event)
	if result.Error != nil {
		log.Printf("Database Error (Create): %v", result.Error)
		c.JSON(500, gin.H{"error": "Failed to create event. Database error."})
		return
	}

	c.JSON(201, Event)
}
