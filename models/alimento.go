package models

import (
	"time"
)

type Alimento struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Nombre        string    `json:"nombre" gorm:"not null;index"`
	
	// Información nutricional por 100g
	Calorias      float64   `json:"calorias"`
	Proteinas     float64   `json:"proteinas"`     // gramos
	Carbohidratos float64   `json:"carbohidratos"` // gramos
	Grasas        float64   `json:"grasas"`        // gramos
	Fibra         float64   `json:"fibra"`         // gramos
	
	// Relación con Tags (muchos a muchos)
	Tags          []Tag     `json:"tags" gorm:"many2many:alimento_tags;"`
	
	// Recursos visuales
	ImagenURL     string    `json:"imagen_url"`
	IconURL       string    `json:"icon_url"` // URL del SVG
	
	// Metadatos
	Activo        bool      `json:"activo" gorm:"default:true"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName especifica el nombre de la tabla
func (Alimento) TableName() string {
	return "alimentos"
}