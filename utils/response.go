package utils

import (
	"encoding/json"
	"net/http"
)

// Response estructura para respuestas exitosas
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta información de paginación
type Meta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalPages  int   `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
}

// ErrorResponse estructura para respuestas de error
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// RespondJSON envía una respuesta JSON exitosa
func RespondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := Response{
		Success: true,
		Data:    data,
	}
	
	json.NewEncoder(w).Encode(response)
}

// RespondJSONWithMeta envía una respuesta JSON con metadata de paginación
func RespondJSONWithMeta(w http.ResponseWriter, statusCode int, data interface{}, meta *Meta) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	}
	
	json.NewEncoder(w).Encode(response)
}

// RespondError envía una respuesta de error
func RespondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := ErrorResponse{
		Success: false,
		Error:   message,
	}
	
	json.NewEncoder(w).Encode(response)
}

// RespondMessage envía una respuesta con un mensaje
func RespondMessage(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := Response{
		Success: true,
		Message: message,
	}
	
	json.NewEncoder(w).Encode(response)
}