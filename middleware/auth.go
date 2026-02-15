package middleware

import (
	"food-catalog-api/config"
	"food-catalog-api/utils"
	"net/http"
	"strings"
)

// AuthMiddleware valida el token de autorización
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtener el token del header Authorization
		authHeader := r.Header.Get("Authorization")
		
		if authHeader == "" {
			utils.RespondError(w, http.StatusUnauthorized, "Token de autorización requerido")
			return
		}
		
		// Extraer el token (formato: "Bearer TOKEN" o solo "TOKEN")
		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		
		// Validar el token
		if !config.AppConfig.IsValidToken(token) {
			utils.RespondError(w, http.StatusUnauthorized, "Token de autorización inválido")
			return
		}
		
		// Token válido, continuar con la petición
		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware agrega headers CORS
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// Manejar preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware registra las peticiones
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log básico de la petición
		// log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}