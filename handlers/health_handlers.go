package handlers

import (
	"food-catalog-api/database"
	"food-catalog-api/utils"
	"net/http"
)

// HealthCheck verifica el estado del servidor y la base de datos
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Verificar conexión a la base de datos
	db := database.GetDB()
	sqlDB, err := db.DB()
	
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error en la base de datos")
		return
	}
	
	if err := sqlDB.Ping(); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Base de datos no disponible")
		return
	}
	
	response := map[string]interface{}{
		"status":   "healthy",
		"database": "connected",
	}
	
	utils.RespondJSON(w, http.StatusOK, response)
}