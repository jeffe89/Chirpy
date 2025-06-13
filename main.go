package main

import (
	"chirpy/internal/database"
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Struct for in-memory data
type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
}

func main() {

	// Define constants
	const filepathRoot = "."
	const port = "8080"

	// Load .env file into environment variables
	godotenv.Load()

	// Get DB_URL from environment
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	// Get platform from environment
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	// Get jwtSecret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	// Open a connection to database
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	// Initialize an apiConfig struct
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		jwtSecret:      jwtSecret,
	}

	// Create a new http.ServeMux
	mux := http.NewServeMux()

	// Setup file server handler with the /app/ path
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	// *** API ***
	// Register a handler function for the /api/healthz path to display status of server
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	// Register a handler function for the /api/login path to login a user with credentials
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	// Register a handler function for the /api/refresh path to refresh token
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	// Register a handler function for the /api/revoke path to revoke a token
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	// Register a handler function for the /api/users path allowing users to be created
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	// Register a handler function for the /api/Chirps path to create chirps
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	// Register a handler function for the /api/Chirps path to retreive all chirps
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
	// Register a handler function for the /api/Chirps path to retreive one specified chirp
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpsGet)

	// *** ADMIN ***
	// Register a handler function for the /admin/reset path to reset hit count
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	// Register a handler function for the /admin/metrics path to display hit count
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	// Create a new HTTP server struct
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Log information on files being served on particular port
	log.Printf("Serving on port: %s\n", port)
	// Use servers ListenAndServe method to start server
	log.Fatal(srv.ListenAndServe())
}
