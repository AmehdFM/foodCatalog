package models

import (
	"time"
)

type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Nombre    string    `json:"nombre" gorm:"not null;uniqueIndex"`
	
	// Recursos visuales
	ImagenURL string    `json:"imagen_url"` // URL de imagen
	IconURL   string    `json:"icon_url"`   // URL de icono/SVG
	Color     string    `json:"color"`      // para UI, ej: "#FF5733"
	
	// Relación inversa (opcional para queries)
	Alimentos []Alimento `json:"-" gorm:"many2many:alimento_tags;"`
	
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// TableName especifica el nombre de la tabla
func (Tag) TableName() string {
	return "tags"
}