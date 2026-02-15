package handlers

import (
	"food-catalog-api/database"
	"food-catalog-api/models"
	"food-catalog-api/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// GetAlimentos obtiene una lista de alimentos con filtros
func GetAlimentos(w http.ResponseWriter, r *http.Request) {
	params := utils.ParseQueryParams(r)
	
	db := database.GetDB()
	var alimentos []models.Alimento
	var total int64
	
	// Construir la query base
	query := db.Model(&models.Alimento{}).Preload("Tags")
	
	// Aplicar filtros
	query = applyAlimentoFilters(query, params)
	
	// Contar total de registros
	query.Count(&total)
	
	// Aplicar ordenamiento
	if params.SortBy != "" {
		orderClause := params.SortBy + " " + params.SortOrder
		query = query.Order(orderClause)
	} else {
		query = query.Order("nombre asc")
	}
	
	// Aplicar paginación
	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)
	
	// Ejecutar query
	if err := query.Find(&alimentos).Error; err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error al obtener alimentos")
		return
	}
	
	// Calcular metadata
	meta := utils.CalculateMeta(params.Page, params.Limit, total)
	
	utils.RespondJSONWithMeta(w, http.StatusOK, alimentos, meta)
}

// GetAlimentoByID obtiene un alimento por su ID
func GetAlimentoByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	
	db := database.GetDB()
	var alimento models.Alimento
	
	if err := db.Preload("Tags").First(&alimento, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(w, http.StatusNotFound, "Alimento no encontrado")
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Error al obtener alimento")
		return
	}
	
	utils.RespondJSON(w, http.StatusOK, alimento)
}

// GetAlimentosByIDs obtiene múltiples alimentos por sus IDs
func GetAlimentosByIDs(w http.ResponseWriter, r *http.Request) {
	idsParam := r.URL.Query().Get("ids")
	
	if idsParam == "" {
		utils.RespondError(w, http.StatusBadRequest, "Parámetro 'ids' requerido")
		return
	}
	
	// Parsear los IDs
	idsStr := strings.Split(idsParam, ",")
	var ids []int
	
	for _, idStr := range idsStr {
		id, err := strconv.Atoi(strings.TrimSpace(idStr))
		if err == nil && id > 0 {
			ids = append(ids, id)
		}
	}
	
	if len(ids) == 0 {
		utils.RespondError(w, http.StatusBadRequest, "No se proporcionaron IDs válidos")
		return
	}
	
	db := database.GetDB()
	var alimentos []models.Alimento
	
	if err := db.Preload("Tags").Where("id IN ?", ids).Find(&alimentos).Error; err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error al obtener alimentos")
		return
	}
	
	utils.RespondJSON(w, http.StatusOK, alimentos)
}

// applyAlimentoFilters aplica filtros a la query de alimentos
func applyAlimentoFilters(query *gorm.DB, params utils.QueryParams) *gorm.DB {
	// Filtro de búsqueda por nombre
	if params.Search != "" {
		query = query.Where("LOWER(nombre) LIKE ?", "%"+strings.ToLower(params.Search)+"%")
	}
	
	// Filtro por tags
	if len(params.Tags) > 0 {
		// Subconsulta para encontrar alimentos que tienen todos los tags especificados
		query = query.Joins("JOIN alimento_tags ON alimento_tags.alimento_id = alimentos.id").
			Joins("JOIN tags ON tags.id = alimento_tags.tag_id").
			Where("LOWER(tags.nombre) IN ?", toLowerSlice(params.Tags)).
			Group("alimentos.id").
			Having("COUNT(DISTINCT tags.id) = ?", len(params.Tags))
	}
	
	// Filtro por estado activo
	if params.Activo != nil {
		query = query.Where("activo = ?", *params.Activo)
	}
	
	// Filtro por fecha de creación
	if params.DateFrom != nil {
		query = query.Where("created_at >= ?", params.DateFrom)
	}
	
	if params.DateTo != nil {
		query = query.Where("created_at <= ?", params.DateTo)
	}
	
	return query
}

// toLowerSlice convierte un slice de strings a minúsculas
func toLowerSlice(slice []string) []string {
	result := make([]string, len(slice))
	for i, s := range slice {
		result[i] = strings.ToLower(s)
	}
	return result
}