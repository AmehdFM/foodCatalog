package database

import (
	"food-catalog-api/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect establece la conexión con la base de datos
func Connect(dsn string) error {
	var err error
	
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	
	if err != nil {
		return err
	}
	
	log.Println("✓ Conexión exitosa con la base de datos")
	return nil
}

// AutoMigrate ejecuta las migraciones automáticas
func AutoMigrate() error {
	err := DB.AutoMigrate(
		&models.Alimento{},
		&models.Tag{},
	)
	
	if err != nil {
		return err
	}
	
	log.Println("✓ Migraciones ejecutadas correctamente")
	return nil
}

// GetDB retorna la instancia de la base de datos
func GetDB() *gorm.DB {
	return DB
}