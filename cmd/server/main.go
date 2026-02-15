package main

import (
	"food-catalog-api/config"
	"food-catalog-api/database"
	"food-catalog-api/handlers"
	"food-catalog-api/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
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
	
	// Crear router
	router := mux.NewRouter()
	
	// Aplicar middlewares globales
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.LoggingMiddleware)
	
	// Ruta de health check (sin autenticación)
	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	
	// Crear subrouter con autenticación
	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(middleware.AuthMiddleware)
	
	// Rutas de Alimentos
	apiRouter.HandleFunc("/alimentos", handlers.GetAlimentos).Methods("GET")
	apiRouter.HandleFunc("/alimentos/{id:[0-9]+}", handlers.GetAlimentoByID).Methods("GET")
	apiRouter.HandleFunc("/alimentos/batch", handlers.GetAlimentosByIDs).Methods("GET")
	
	// Rutas de Tags
	apiRouter.HandleFunc("/tags", handlers.GetTags).Methods("GET")
	apiRouter.HandleFunc("/tags/{id:[0-9]+}", handlers.GetTagByID).Methods("GET")
	apiRouter.HandleFunc("/tags/batch", handlers.GetTagsByIDs).Methods("GET")
	apiRouter.HandleFunc("/tags/nombre/{nombre}", handlers.GetTagByNombre).Methods("GET")
	
	// Iniciar servidor
	port := ":" + config.AppConfig.Server.Port
	log.Printf("🚀 Servidor iniciado en http://localhost%s", port)
	log.Printf("📚 API disponible en http://localhost%s/api/v1", port)
	log.Fatal(http.ListenAndServe(port, router))
}