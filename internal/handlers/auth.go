package handlers

import (
	"errors"
	"eventflow/internal/database"
	"eventflow/internal/models"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "your-secret-key-change-in-production"
	}
	return secret
}

// @Summary Регистрация нового организатора
// @Description Создание нового аккаунта организатора с email и паролем
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.CreateOrganizerRequest true "Данные регистрации"
// @Success 201 {object} models.LoginResponse "Успешная регистрация"
// @Failure 400 {object} map[string]string "Ошибка валидации"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var req models.CreateOrganizerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if req.Role != "admin" && req.Role != "organizer" {
		c.JSON(400, gin.H{"error": "Role must be either 'admin' or 'organizer'"})
		return
	}

	var existingOrganizer models.Organizer
	result := database.DB.Where("email = ?", req.Email).First(&existingOrganizer)
	if result.Error == nil {
		c.JSON(400, gin.H{"error": "Email already registered"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Password hashing error: %v", err)
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	organizer := models.Organizer{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	createResult := database.DB.Create(&organizer)
	if createResult.Error != nil {
		log.Printf("Database Error (Create): %v", createResult.Error)
		c.JSON(500, gin.H{"error": "Failed to create organizer"})
		return
	}

	accessToken, err := generateAccessToken(organizer.ID, organizer.Email, organizer.Role)
	if err != nil {
		log.Printf("Access token generation error: %v", err)
		c.JSON(500, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := generateRefreshToken(organizer.ID)
	if err != nil {
		log.Printf("Refresh token generation error: %v", err)
		c.JSON(500, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	organizer.RefreshToken = refreshToken
	database.DB.Save(&organizer)

	c.JSON(201, models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    300,
		User:         organizer,
	})
}

// @Summary Вход в систему
// @Description Аутентификация пользователя по email и паролю
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Данные для входа"
// @Success 200 {object} models.LoginResponse "Успешный вход"
// @Failure 401 {object} map[string]string "Неверные учетные данные"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var organizer models.Organizer
	result := database.DB.Where("email = ?", req.Email).First(&organizer)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(401, gin.H{"error": "Invalid email or password"})
		} else {
			c.JSON(500, gin.H{"error": "Database error"})
		}
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(organizer.Password), []byte(req.Password))
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	accessToken, err := generateAccessToken(organizer.ID, organizer.Email, organizer.Role)
	if err != nil {
		log.Printf("Access token generation error: %v", err)
		c.JSON(500, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := generateRefreshToken(organizer.ID)
	if err != nil {
		log.Printf("Refresh token generation error: %v", err)
		c.JSON(500, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	organizer.RefreshToken = refreshToken
	database.DB.Save(&organizer)

	c.JSON(200, models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    300,
		User:         organizer,
	})
}

func generateAccessToken(userID uint, email string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"type":    "access",
		"exp":     time.Now().Add(time.Minute * 5).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func generateRefreshToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// @Summary Обновление access токена
// @Description Обновляет access token используя refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} models.LoginResponse
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]
func RefreshAccessToken(c *gin.Context) {
	var req models.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"error": "Invalid refresh token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(401, gin.H{"error": "Invalid token claims"})
		return
	}

	tokenType, _ := claims["type"].(string)
	if tokenType != "refresh" {
		c.JSON(401, gin.H{"error": "Not a refresh token"})
		return
	}

	userID := uint(claims["user_id"].(float64))

	var organizer models.Organizer
	if err := database.DB.First(&organizer, userID).Error; err != nil {
		c.JSON(401, gin.H{"error": "User not found"})
		return
	}

	if organizer.RefreshToken != req.RefreshToken {
		c.JSON(401, gin.H{"error": "Refresh token does not match"})
		return
	}

	accessToken, err := generateAccessToken(organizer.ID, organizer.Email, organizer.Role)
	if err != nil {
		log.Printf("Access token generation error: %v", err)
		c.JSON(500, gin.H{"error": "Failed to generate access token"})
		return
	}

	c.JSON(200, models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken,
		ExpiresIn:    300,
		User:         organizer,
	})
}

// @Summary Получение текущего пользователя
// @Description Возвращает информацию о текущем авторизованном пользователе
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Organizer "Информация о пользователе"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 404 {object} map[string]string "Пользователь не найден"
// @Router /auth/me [get]
func GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var organizer models.Organizer
	result := database.DB.First(&organizer, userID)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, organizer)
}
