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

// GetTags obtiene una lista de tags con filtros
func GetTags(w http.ResponseWriter, r *http.Request) {
	params := utils.ParseQueryParams(r)
	
	db := database.GetDB()
	var tags []models.Tag
	var total int64
	
	// Construir la query base
	query := db.Model(&models.Tag{})
	
	// Aplicar filtros
	query = applyTagFilters(query, params)
	
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
	if err := query.Find(&tags).Error; err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error al obtener tags")
		return
	}
	
	// Calcular metadata
	meta := utils.CalculateMeta(params.Page, params.Limit, total)
	
	utils.RespondJSONWithMeta(w, http.StatusOK, tags, meta)
}

// GetTagByID obtiene un tag por su ID
func GetTagByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	
	db := database.GetDB()
	var tag models.Tag
	
	if err := db.First(&tag, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(w, http.StatusNotFound, "Tag no encontrado")
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Error al obtener tag")
		return
	}
	
	utils.RespondJSON(w, http.StatusOK, tag)
}

// GetTagsByIDs obtiene múltiples tags por sus IDs
func GetTagsByIDs(w http.ResponseWriter, r *http.Request) {
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
	var tags []models.Tag
	
	if err := db.Where("id IN ?", ids).Find(&tags).Error; err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error al obtener tags")
		return
	}
	
	utils.RespondJSON(w, http.StatusOK, tags)
}

// GetTagByNombre obtiene un tag por su nombre
func GetTagByNombre(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nombre := vars["nombre"]
	
	if nombre == "" {
		utils.RespondError(w, http.StatusBadRequest, "Nombre requerido")
		return
	}
	
	db := database.GetDB()
	var tag models.Tag
	
	if err := db.Where("LOWER(nombre) = ?", strings.ToLower(nombre)).First(&tag).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(w, http.StatusNotFound, "Tag no encontrado")
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Error al obtener tag")
		return
	}
	
	utils.RespondJSON(w, http.StatusOK, tag)
}

// applyTagFilters aplica filtros a la query de tags
func applyTagFilters(query *gorm.DB, params utils.QueryParams) *gorm.DB {
	// Filtro de búsqueda por nombre
	if params.Search != "" {
		query = query.Where("LOWER(nombre) LIKE ?", "%"+strings.ToLower(params.Search)+"%")
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