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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetOrganizerById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(400, gin.H{"error": "ID parameter is required"})
		return
	}

	var organizer models.Organizer

	result := database.DB.First(&organizer, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Organizer not found"})
		} else {
			c.JSON(500, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(200, organizer)
}

// @Summary Получить список организаторов
// @Description Возвращает список всех организаторов событий с пагинацией
// @Tags Organizers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param range query string false "Пагинация [start, end]"
// @Param sort query string false "Сортировка [field, order]"
// @Success 200 {array} models.Organizer
// @Header 200 {string} X-Total-Count "Общее количество записей"
// @Header 200 {string} Content-Range "Диапазон записей"
// @Router /organizers [get]
func GetOrganizers(c *gin.Context) {
	var organizers []models.Organizer
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

	countResult := database.DB.Model(&models.Organizer{}).Count(&total)
	if countResult.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve total record count"})
		return
	}

	contentRange := fmt.Sprintf("organizers %d-%d/%d", start, end, total)

	result := database.DB.
		Limit(limit).
		Offset(offset).
		Order(sortField + " " + sortOrder).
		Find(&organizers)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	c.Header("Content-Range", contentRange)
	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(200, organizers)
}

func UpdateOrganizer(c *gin.Context) {
	id := c.Param("id")

	var input models.CreateOrganizerRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if input.Role != "admin" && input.Role != "organizer" {
		c.JSON(400, gin.H{"error": "Role must be either 'admin' or 'organizer'"})
		return
	}

	var organizer models.Organizer

	result := database.DB.Model(&organizer).Where("id = ?", id).Updates(models.Organizer{Name: input.Name, Email: input.Email, Role: input.Role})

	if result.Error != nil {
		log.Printf("Database Error (Update): %v", result.Error)
		c.JSON(500, gin.H{"error": "Failed to update organizer. Database error."})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Organizer not found."})
		return
	}

	database.DB.First(&organizer, id)

	c.JSON(200, organizer)
}

func DeleteOrganizer(c *gin.Context) {
	id := c.Param("id")

	result := database.DB.Delete(&models.Organizer{}, id)

	if result.Error != nil {
		log.Printf("Database Error (Delete): %v", result.Error)
		c.JSON(500, gin.H{"error": "Failed to delete organizer. Database error."})
		return
	}

	c.JSON(200, gin.H{})
}

// @Summary Создать организатора
// @Description Создает нового организатора (требуется авторизация admin)
// @Tags Organizers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizer body models.CreateOrganizerRequest true "Данные организатора"
// @Success 201 {object} models.Organizer
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /organizers [post]
func PostOrganizer(c *gin.Context) {
	var newOrganizer models.CreateOrganizerRequest

	if err := c.ShouldBindJSON(&newOrganizer); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if newOrganizer.Role != "admin" && newOrganizer.Role != "organizer" {
		c.JSON(400, gin.H{"error": "Role must be either 'admin' or 'organizer'"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newOrganizer.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Password hashing error: %v", err)
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	organizer := models.Organizer{
		Name:     newOrganizer.Name,
		Email:    newOrganizer.Email,
		Password: string(hashedPassword),
		Role:     newOrganizer.Role,
	}

	result := database.DB.Create(&organizer)

	if result.Error != nil {
		log.Printf("Database Error (Create): %v", result.Error)
		c.JSON(500, gin.H{"error": "Failed to create organizer. Database error."})
		return
	}

	c.JSON(201, organizer)
}
