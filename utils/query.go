package utils

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

// QueryParams estructura para parsear parámetros de query
type QueryParams struct {
	Page       int
	Limit      int
	Search     string
	Tags       []string
	DateFrom   *time.Time
	DateTo     *time.Time
	SortBy     string
	SortOrder  string
	Activo     *bool
}

// ParseQueryParams extrae y valida los parámetros de la URL
func ParseQueryParams(r *http.Request) QueryParams {
	query := r.URL.Query()
	
	params := QueryParams{
		Page:      parseIntParam(query.Get("page"), 1),
		Limit:     parseIntParam(query.Get("limit"), 20),
		Search:    strings.TrimSpace(query.Get("search")),
		Tags:      parseArrayParam(query.Get("tags")),
		SortBy:    query.Get("sort_by"),
		SortOrder: parseSortOrder(query.Get("sort_order")),
	}
	
	// Limitar el máximo de resultados por página
	if params.Limit > 100 {
		params.Limit = 100
	}
	if params.Limit < 1 {
		params.Limit = 20
	}
	
	// Parsear fechas
	if dateFrom := query.Get("date_from"); dateFrom != "" {
		if t, err := time.Parse("2006-01-02", dateFrom); err == nil {
			params.DateFrom = &t
		}
	}
	
	if dateTo := query.Get("date_to"); dateTo != "" {
		if t, err := time.Parse("2006-01-02", dateTo); err == nil {
			params.DateTo = &t
		}
	}
	
	// Parsear activo
	if activoStr := query.Get("activo"); activoStr != "" {
		if activo, err := strconv.ParseBool(activoStr); err == nil {
			params.Activo = &activo
		}
	}
	
	return params
}

// parseIntParam parsea un parámetro string a int con valor por defecto
func parseIntParam(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return defaultValue
	}
	
	return parsed
}

// parseArrayParam parsea una string separada por comas en un slice
func parseArrayParam(value string) []string {
	if value == "" {
		return []string{}
	}
	
	parts := strings.Split(value, ",")
	var result []string
	
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	
	return result
}

// parseSortOrder valida el orden de ordenamiento
func parseSortOrder(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "asc" || value == "desc" {
		return value
	}
	return "asc"
}

// CalculateMeta calcula la metadata de paginación
func CalculateMeta(page, limit int, total int64) *Meta {
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}
	
	return &Meta{
		CurrentPage: page,
		PerPage:     limit,
		TotalPages:  totalPages,
		TotalItems:  total,
	}
}