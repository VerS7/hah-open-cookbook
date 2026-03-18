package main

import (
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"

	"open-hah-cookbook/internal/api"
	l "open-hah-cookbook/internal/logger"
	"open-hah-cookbook/internal/storage"
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
	requireEnv("DEBUG", "DATA_DIR")

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

		l.Default.Infof("initializing storage with name: '%s' at '%s'...", filename, datapath)

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

			continue
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

	defer aggregator.CloseAll()

	// Recipes servers
	recipeServers := make(map[string]*api.RecipesAPIServer)

	// Cookbook versions
	cookbookVersions := make([]string, 0)

	for _, aggregatorEntry := range aggregator.GetEntries() {
		recipeServers[aggregatorEntry.Alias] = api.NewRecipesAPIServer(aggregatorEntry.Storage)
		cookbookVersions = append(cookbookVersions, aggregatorEntry.Alias)
	}

	l.Default.Infof("cookbook versions found: '%s'", strings.Join(cookbookVersions, ", "))

	rootMux := http.NewServeMux()

	// Cookbook versions (e.g. w161, w16, old...)
	rootMux.HandleFunc("GET /versions", func(w http.ResponseWriter, r *http.Request) {
		api.WriteJSON(w, 200, cookbookVersions)
	})

	for name, recipeServer := range recipeServers {
		// .../api/[cookbook version]/... with USER auth
		rootMux.HandleFunc(fmt.Sprintf("GET /api/%s/recipes", name), recipeServer.FilteredQueryHandler)
		rootMux.HandleFunc(fmt.Sprintf("GET /api/%s/export", name), recipeServer.ExportHandler)

		// Users can send recipes only to non-archived cookbooks
		if !strings.Contains(name, "archived") {
			go recipeServer.RunRecipeListener()

			rootMux.HandleFunc(fmt.Sprintf("POST /api/%s/recipe", name), recipeServer.RecipeHandler)
		}

		rootMux.HandleFunc(fmt.Sprintf("DELETE /api/%s/recipe/{id}", name), recipeServer.RecipeHandler)
	}

	coreHandler := api.RecoveryMiddleware(api.LoggingMiddleware(rootMux))
	if DEBUG {
		coreHandler = api.AllowCORSMiddleware(coreHandler)
	}

	l.Default.Debugf("server running at :%d", *PORT)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *PORT), coreHandler); err != nil {
		l.Default.Fatal(err)
	}
}
