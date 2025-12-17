package database

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		errEnv = godotenv.Load("../.env")
		if errEnv != nil {
			log.Println("Warning: .env file not found, using environment variables")
		}
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("error connecting the database. Error: %s", err)
	}

	log.Println("🚀 Database success the connected :3")

	RunMigrations()
}

func RunMigrations() {
	migrateURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	migrationsPath := "file://migrations"
	if _, err := os.Stat("migrations"); os.IsNotExist(err) {
		migrationsPath = "file://../migrations"
	}

	m, err := migrate.New(migrationsPath, migrateURL)
	if err != nil {
		log.Fatalf("❌ Failed to create migrate instance: %s", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("✅ No new migrations to apply")
		} else {
			log.Fatalf("❌ Migration failed: %s", err)
		}
	} else {
		log.Println("✅ Migrations applied successfully")
	}
}
