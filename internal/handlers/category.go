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

func GetCategoryById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(400, gin.H{"error": "ID parameter is required"})
		return
	}

	var category models.Category

	result := database.DB.First(&category, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Category not found"})
		} else {
			c.JSON(500, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(200, category)
}

// @Summary Получить список категорий
// @Description Возвращает список всех категорий событий с пагинацией
// @Tags Categories
// @Accept json
// @Produce json
// @Param range query string false "Пагинация [start, end]"
// @Param sort query string false "Сортировка [field, order]"
// @Success 200 {array} models.Category
// @Header 200 {string} X-Total-Count "Общее количество записей"
// @Header 200 {string} Content-Range "Диапазон записей"
// @Router /categories [get]
func GetCategories(c *gin.Context) {
	var categories []models.Category
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

	countResult := database.DB.Model(&models.Category{}).Count(&total)
	if countResult.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve total record count"})
		return
	}

	contentRange := fmt.Sprintf("categories %d-%d/%d", start, end, total)

	result := database.DB.
		Limit(limit).
		Offset(offset).
		Order(sortField + " " + sortOrder).
		Find(&categories)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	c.Header("Content-Range", contentRange)
	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(200, categories)
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var input models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var category models.Category

	result := database.DB.Model(&category).Where("id = ?", id).Updates(models.Category{Name: input.Name})

	if result.Error != nil {
		log.Printf("Database Error (Update): %v", result.Error)
		c.JSON(500, gin.H{"error": "Failed to update category. Database error."})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Category not found."})
		return
	}

	database.DB.First(&category, id)

	c.JSON(200, category)
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	result := database.DB.Delete(&models.Category{}, id)

	if result.Error != nil {
		log.Printf("Database Error (Delete): %v", result.Error)
		c.JSON(500, gin.H{"error": "Failed to delete category. Database error."})
		return
	}

	c.JSON(200, gin.H{})
}

// @Summary Создать категорию
// @Description Создает новую категорию событий
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body models.CreateCategoryRequest true "Данные категории"
// @Success 201 {object} models.Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories [post]
func PostCategory(c *gin.Context) {
	var newCategories models.CreateCategoryRequest

	if err := c.ShouldBindJSON(&newCategories); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{
		Name: newCategories.Name,
	}

	result := database.DB.Create(&category)

	if result.Error != nil {
		log.Printf("Database Error (Create): %v", result.Error)
		c.JSON(500, gin.H{"error": "Failed to create category. Database error."})
		return
	}

	c.JSON(201, category)
}
