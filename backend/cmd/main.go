package main

import (
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

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

func getDataFiles(dataDir string) []string {
	filepaths := make([]string, 0)

	if err := filepath.Walk(dataDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".db" {
			filepaths = append(filepaths, path)
		}

		return nil
	}); err != nil {
		l.Default.Fatal(err)
	}

	return filepaths
}

func main() {
	requireEnv("DEBUG", "DATA_DIR", "ADMIN_USERNAME", "ADMIN_PASSWORD")

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

	// Databases
	aggregator := storage.NewAggregator()

	dataFiles := getDataFiles(os.Getenv("DATA_DIR"))

	for _, datapath := range dataFiles {
		filename := strings.Replace(filepath.Base(datapath), filepath.Ext(datapath), "", 1)

		// Special case for users.db
		if filename == "users" {
			if _, err := aggregator.NewStorage(storage.StorageEntryConfig{
				Name:    filename,
				File:    datapath,
				Alias:   "users",
				Schemas: []string{storage.USERS},
			}); err != nil {
				l.Default.Fatal(err)
			}
		}

		// Special case for archived bases
		if strings.Contains(filename, "archived") {
			if _, err := aggregator.NewStorage(storage.StorageEntryConfig{
				Name:    "archived_cookbook_" + filename,
				File:    datapath,
				Alias:   filename,
				Schemas: []string{storage.RECIPES, storage.INGREDIENTS},
			}); err != nil {
				l.Default.Fatal(err)
			}
		}

		// Other cases interpret as regular cookbooks
		if _, err := aggregator.NewStorage(storage.StorageEntryConfig{
			Name:    "cookbook_" + filename,
			File:    datapath,
			Alias:   filename,
			Schemas: []string{storage.RECIPES, storage.INGREDIENTS},
		}); err != nil {
			l.Default.Fatal(err)
		}
	}

	// Fallback if users.db not presented
	if !slices.ContainsFunc(dataFiles, func(v string) bool {
		return strings.Contains(filepath.Base(v), "users")
	}) {
		if _, err := aggregator.NewStorage(storage.StorageEntryConfig{
			Name:    "users",
			File:    os.Getenv("DATA_DIR") + "users.db",
			Alias:   "users",
			Schemas: []string{storage.USERS},
		}); err != nil {
			l.Default.Fatal(err)
		}
	}

	defer aggregator.CloseAll()

	// Default admin
	users, _ := aggregator.GetEntry("users")
	if err := users.Storage.AddUser(os.Getenv("ADMIN_USERNAME"), os.Getenv("ADMIN_PASSWORD"), true); err != nil {
		l.Default.Println(err)
	}

	// Users server
	usersServer := api.NewUsersAPIServer(users.Storage)
	go usersServer.UpdateSessions()

	// Recipes servers
	recipeServers := make(map[string]*api.RecipesAPIServer)

	// Cookbook versions
	cookbookVersions := make([]string, 0)

	for _, aggregatorEntry := range aggregator.GetEntries() {
		// Skip Users storage, only recipes
		if aggregatorEntry.Alias != "users" {
			recipeServers[aggregatorEntry.Alias] = api.NewRecipesAPIServer(aggregatorEntry.Storage)
			cookbookVersions = append(cookbookVersions, aggregatorEntry.Alias)
		}
	}

	rootMux := http.NewServeMux()
	apiMux := http.NewServeMux()
	sharedMux := http.NewServeMux()
	adminMux := http.NewServeMux()

	// .../api/...
	apiMux.HandleFunc("POST /login", usersServer.LoginHandler)

	// Cookbook versions (e.g. w161, w16, old...)
	apiMux.HandleFunc("GET /versions", func(w http.ResponseWriter, r *http.Request) {

		api.WriteJSON(w, 200, cookbookVersions)
	})

	adminMux.HandleFunc("POST /users", api.TODOHandler)
	adminMux.HandleFunc("GET /users", api.TODOHandler)
	adminMux.HandleFunc("GET /users/{id}", api.TODOHandler)
	adminMux.HandleFunc("DELETE /users/{id}", api.TODOHandler)

	for name, recipeServer := range recipeServers {
		// .../api/[cookbook version]/... with USER auth
		sharedMux.HandleFunc(fmt.Sprintf("GET /%s/recipes", name), recipeServer.FilteredQueryHandler)
		sharedMux.HandleFunc(fmt.Sprintf("GET /%s/export", name), recipeServer.ExportHandler)
		// Users can send recipes only to non-archived cookbooks
		if !strings.Contains(name, "archived") {
			go recipeServer.RunRecipeListener()

			sharedMux.HandleFunc(fmt.Sprintf("POST /%s/recipe", name), recipeServer.RecipeHandler)
		}

		// .../api/[cookbook version]/... with ADMIN auth
		adminMux.HandleFunc(fmt.Sprintf("DELETE /%s/recipe/{id}", name), recipeServer.RecipeHandler)
	}

	sharedHandler := usersServer.AuthMiddleware(sharedMux)
	adminHandler := usersServer.AdminMiddleware(adminMux)
	apiMux.Handle("/", sharedHandler)
	apiMux.Handle("/admin", adminHandler)
	rootMux.Handle("/api/", http.StripPrefix("/api", apiMux))

	coreHandler := api.RecoveryMiddleware(api.LoggingMiddleware(rootMux))
	if DEBUG {
		coreHandler = api.AllowCORSMiddleware(coreHandler)
	}

	l.Default.Debugf("server running at :%d", *PORT)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *PORT), coreHandler); err != nil {
		l.Default.Fatal(err)
	}
}
