package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"open-hah-cookbook/internal/api"
	l "open-hah-cookbook/internal/logger"
	"open-hah-cookbook/internal/storage"

	"github.com/joho/godotenv"
)

var (
	ENV_FILE = flag.String("env", "", ".env config file")
	PORT     = flag.Int("port", 8080, "app port")
)

func init() {
	flag.Parse()

	env := *ENV_FILE
	if len(env) == 0 {
		return
	}

	if err := godotenv.Load(env); err != nil {
		l.Default.Fatalf("cannot load .env file: %v", err)
	}
}

func requireEnv(keys ...string) {
	for _, key := range keys {
		if os.Getenv(key) == "" {
			l.Default.Fatalf("required env var is not set: %s", key)
		}
	}
}

func main() {
	requireEnv("DEBUG", "DATABASE_FILE", "ADMIN_USERNAME", "ADMIN_PASSWORD")

	// Parse DEBUG-mode from env
	var DEBUG bool
	switch os.Getenv("DEBUG") {
	case "true":
		DEBUG = true
	case "false":
		DEBUG = false
	default:
		l.Default.Fatal("wrong DEBUG param")
	}

	// Database
	var db, err = storage.NewDB(os.Getenv("DATABASE_FILE"))
	if err != nil {
		l.Default.Fatal(err)
	}
	if err := db.InitSchema(); err != nil {
		l.Default.Fatal(err)
	}
	defer db.Close()

	// Default admin
	db.AddUser(os.Getenv("ADMIN_USERNAME"), os.Getenv("ADMIN_PASSWORD"), true)
	if err != nil {
		l.Default.Println(err)
	}

	// Server
	server := api.NewAPIServer(db)
	go server.UpdateSessions()
	go server.RunRecipeListener()

	rootMux := http.NewServeMux()

	// .../api/...
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("POST /login", server.LoginHandler)

	// .../api/... with USER auth
	sharedMux := http.NewServeMux()
	sharedMux.HandleFunc("POST /recipe", server.RecipeHandler)
	sharedMux.HandleFunc("GET /recipes", server.FilteredQueryHandler)
	sharedMux.HandleFunc("GET /export", server.ExportHandler)

	// .../api/... with ADMIN auth
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("DELETE /recipe/{id}", server.RecipeHandler)
	adminMux.HandleFunc("POST /users", server.TODOHandler)
	adminMux.HandleFunc("GET /users", server.TODOHandler)
	adminMux.HandleFunc("GET /users/{id}", server.TODOHandler)
	adminMux.HandleFunc("DELETE /users/{id}", server.TODOHandler)

	sharedHandler := server.AuthMiddleware(sharedMux)
	adminHandler := server.AdminMiddleware(sharedHandler)

	apiMux.Handle("/", sharedHandler)
	apiMux.Handle("/admin", adminHandler)
	rootMux.Handle("/api/", http.StripPrefix("/api", apiMux))

	coreHandler := server.RecoveryMiddleware(server.LoggingMiddleware(rootMux))
	if DEBUG {
		coreHandler = server.AllowCORSMiddleware(coreHandler)
	}

	l.Default.Debugf("server running at :%d", *PORT)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *PORT), coreHandler); err != nil {
		l.Default.Fatal(err)
	}
}
