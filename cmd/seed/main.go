package main

import (
	"encoding/json"
	"flag"
	"food-catalog-api/config"
	"food-catalog-api/database"
	"food-catalog-api/models"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
)

// SeedData estructura del archivo JSON
type SeedData struct {
	Alimentos []AlimentoSeed `json:"alimentos"`
	Tags      []TagSeed      `json:"tags"`
}

// AlimentoSeed estructura de alimento en el JSON
type AlimentoSeed struct {
	ID            uint     `json:"id"`
	Nombre        string   `json:"nombre"`
	Calorias      float64  `json:"calorias"`
	Proteinas     float64  `json:"proteinas"`
	Carbohidratos float64  `json:"carbohidratos"`
	Grasas        float64  `json:"grasas"`
	Fibra         float64  `json:"fibra"`
	ImagenURL     string   `json:"imagen_url"`
	IconURL       string   `json:"icon_url"`
	Activo        bool     `json:"activo"`
	Tags          []string `json:"tags"` // Nombres de tags
	UpdatedAt     string   `json:"updated_at"` // Formato: "2006-01-02T15:04:05Z"
}

// TagSeed estructura de tag en el JSON
type TagSeed struct {
	ID        uint   `json:"id"`
	Nombre    string `json:"nombre"`
	ImagenURL string `json:"imagen_url"`
	IconURL   string `json:"icon_url"`
	Color     string `json:"color"`
	UpdatedAt string `json:"updated_at"` // Formato: "2006-01-02T15:04:05Z"
}

func main() {
	// Definir flags
	filePath := flag.String("file", "data.json", "Ruta al archivo JSON con los datos")
	flag.Parse()
	
	// Cargar configuración
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Error al cargar configuración:", err)
	}
	
	// Conectar a la base de datos
	if err := database.Connect(config.AppConfig.GetDSN()); err != nil {
		log.Fatal("Error al conectar con la base de datos:", err)
	}
	
	// Ejecutar migraciones
	if err := database.AutoMigrate(); err != nil {
		log.Fatal("Error al ejecutar migraciones:", err)
	}
	
	// Leer archivo JSON
	log.Printf("Leyendo archivo: %s", *filePath)
	data, err := readSeedFile(*filePath)
	if err != nil {
		log.Fatal("Error al leer archivo JSON:", err)
	}
	
	db := database.GetDB()
	
	// Procesar tags primero (por las relaciones)
	log.Println("\n📦 Procesando tags...")
	if err := processTags(db, data.Tags); err != nil {
		log.Fatal("Error al procesar tags:", err)
	}
	
	// Procesar alimentos
	log.Println("\n🍎 Procesando alimentos...")
	if err := processAlimentos(db, data.Alimentos); err != nil {
		log.Fatal("Error al procesar alimentos:", err)
	}
	
	log.Println("\n✅ Seed completado exitosamente")
}

// readSeedFile lee y parsea el archivo JSON
func readSeedFile(filePath string) (*SeedData, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error al leer archivo: %w", err)
	}
	
	var data SeedData
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, fmt.Errorf("error al parsear JSON: %w", err)
	}
	
	return &data, nil
}

// processTags procesa y actualiza los tags
func processTags(db *gorm.DB, tagSeeds []TagSeed) error {
	created := 0
	updated := 0
	skipped := 0
	
	for _, tagSeed := range tagSeeds {
		var existingTag models.Tag
		result := db.Where("nombre = ?", tagSeed.Nombre).First(&existingTag)
		
		// Parsear la fecha de actualización del JSON
		jsonUpdatedAt, err := time.Parse(time.RFC3339, tagSeed.UpdatedAt)
		if err != nil {
			log.Printf("⚠️  Tag '%s': fecha inválida, usando fecha actual", tagSeed.Nombre)
			jsonUpdatedAt = time.Now()
		}
		
		if result.Error == gorm.ErrRecordNotFound {
			// Tag no existe, crear nuevo
			newTag := models.Tag{
				Nombre:    tagSeed.Nombre,
				ImagenURL: tagSeed.ImagenURL,
				IconURL:   tagSeed.IconURL,
				Color:     tagSeed.Color,
				UpdatedAt: jsonUpdatedAt,
			}
			
			if err := db.Create(&newTag).Error; err != nil {
				return fmt.Errorf("error al crear tag '%s': %w", tagSeed.Nombre, err)
			}
			
			log.Printf("✓ Tag creado: %s (ID: %d)", newTag.Nombre, newTag.ID)
			created++
			
		} else if result.Error != nil {
			return fmt.Errorf("error al buscar tag '%s': %w", tagSeed.Nombre, result.Error)
			
		} else {
			// Tag existe, verificar si necesita actualización
			if jsonUpdatedAt.After(existingTag.UpdatedAt) {
				existingTag.ImagenURL = tagSeed.ImagenURL
				existingTag.IconURL = tagSeed.IconURL
				existingTag.Color = tagSeed.Color
				existingTag.UpdatedAt = jsonUpdatedAt
				
				if err := db.Save(&existingTag).Error; err != nil {
					return fmt.Errorf("error al actualizar tag '%s': %w", tagSeed.Nombre, err)
				}
				
				log.Printf("↻ Tag actualizado: %s (ID: %d)", existingTag.Nombre, existingTag.ID)
				updated++
			} else {
				log.Printf("− Tag sin cambios: %s (ID: %d)", existingTag.Nombre, existingTag.ID)
				skipped++
			}
		}
	}
	
	log.Printf("\n📊 Resumen Tags: %d creados, %d actualizados, %d sin cambios", created, updated, skipped)
	return nil
}

// processAlimentos procesa y actualiza los alimentos
func processAlimentos(db *gorm.DB, alimentoSeeds []AlimentoSeed) error {
	created := 0
	updated := 0
	skipped := 0
	
	for _, alimentoSeed := range alimentoSeeds {
		var existingAlimento models.Alimento
		result := db.Preload("Tags").Where("nombre = ?", alimentoSeed.Nombre).First(&existingAlimento)
		
		// Parsear la fecha de actualización del JSON
		jsonUpdatedAt, err := time.Parse(time.RFC3339, alimentoSeed.UpdatedAt)
		if err != nil {
			log.Printf("⚠️  Alimento '%s': fecha inválida, usando fecha actual", alimentoSeed.Nombre)
			jsonUpdatedAt = time.Now()
		}
		
		// Buscar los tags por nombre
		var tags []models.Tag
		if len(alimentoSeed.Tags) > 0 {
			if err := db.Where("nombre IN ?", alimentoSeed.Tags).Find(&tags).Error; err != nil {
				return fmt.Errorf("error al buscar tags para alimento '%s': %w", alimentoSeed.Nombre, err)
			}
		}
		
		if result.Error == gorm.ErrRecordNotFound {
			// Alimento no existe, crear nuevo
			newAlimento := models.Alimento{
				Nombre:        alimentoSeed.Nombre,
				Calorias:      alimentoSeed.Calorias,
				Proteinas:     alimentoSeed.Proteinas,
				Carbohidratos: alimentoSeed.Carbohidratos,
				Grasas:        alimentoSeed.Grasas,
				Fibra:         alimentoSeed.Fibra,
				ImagenURL:     alimentoSeed.ImagenURL,
				IconURL:       alimentoSeed.IconURL,
				Activo:        alimentoSeed.Activo,
				Tags:          tags,
				UpdatedAt:     jsonUpdatedAt,
			}
			
			if err := db.Create(&newAlimento).Error; err != nil {
				return fmt.Errorf("error al crear alimento '%s': %w", alimentoSeed.Nombre, err)
			}
			
			log.Printf("✓ Alimento creado: %s (ID: %d) con %d tags", newAlimento.Nombre, newAlimento.ID, len(tags))
			created++
			
		} else if result.Error != nil {
			return fmt.Errorf("error al buscar alimento '%s': %w", alimentoSeed.Nombre, result.Error)
			
		} else {
			// Alimento existe, verificar si necesita actualización
			if jsonUpdatedAt.After(existingAlimento.UpdatedAt) {
				existingAlimento.Calorias = alimentoSeed.Calorias
				existingAlimento.Proteinas = alimentoSeed.Proteinas
				existingAlimento.Carbohidratos = alimentoSeed.Carbohidratos
				existingAlimento.Grasas = alimentoSeed.Grasas
				existingAlimento.Fibra = alimentoSeed.Fibra
				existingAlimento.ImagenURL = alimentoSeed.ImagenURL
				existingAlimento.IconURL = alimentoSeed.IconURL
				existingAlimento.Activo = alimentoSeed.Activo
				existingAlimento.UpdatedAt = jsonUpdatedAt
				
				// Actualizar relación con tags
				if err := db.Model(&existingAlimento).Association("Tags").Replace(tags); err != nil {
					return fmt.Errorf("error al actualizar tags del alimento '%s': %w", alimentoSeed.Nombre, err)
				}
				
				if err := db.Save(&existingAlimento).Error; err != nil {
					return fmt.Errorf("error al actualizar alimento '%s': %w", alimentoSeed.Nombre, err)
				}
				
				log.Printf("↻ Alimento actualizado: %s (ID: %d) con %d tags", existingAlimento.Nombre, existingAlimento.ID, len(tags))
				updated++
			} else {
				log.Printf("− Alimento sin cambios: %s (ID: %d)", existingAlimento.Nombre, existingAlimento.ID)
				skipped++
			}
		}
	}
	
	log.Printf("\n📊 Resumen Alimentos: %d creados, %d actualizados, %d sin cambios", created, updated, skipped)
	return nil
}