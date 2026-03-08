package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"

	"open-hah-cookbook/internal/auth"
	l "open-hah-cookbook/internal/logger"
	"open-hah-cookbook/internal/storage"
)

type RecipesAPIServer struct {
	Storage *storage.Storage
	recipes chan storage.FoodRecipe
	mu      sync.RWMutex
}

func NewRecipesAPIServer(st *storage.Storage) *RecipesAPIServer {
	return &RecipesAPIServer{
		st,
		make(chan storage.FoodRecipe),
		sync.RWMutex{},
	}
}

func (sr *RecipesAPIServer) RunRecipeListener() {
	for r := range sr.recipes {
		if err := sr.Storage.AddRecipe(r); err != nil {
			l.Default.Error(err)
		}
		l.Default.Infof("recieved recipe for '%s'", r.ItemName)
	}
}

func (sr *RecipesAPIServer) RecipeHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, H{"error": "cant read request body"})
		return
	}
	defer r.Body.Close()

	if len(body) == 0 {
		WriteJSON(w, http.StatusBadRequest, H{"error": "empty request body"})
		return
	}

	var jsonData []storage.FoodRecipe
	if err := json.Unmarshal(body, &jsonData); err != nil {
		WriteJSON(w, http.StatusBadRequest, H{"error": "cant read recipe json"})
		return
	}

	for _, data := range jsonData {
		sr.recipes <- data
	}

	w.WriteHeader(http.StatusOK)
}

func (sr *RecipesAPIServer) FilteredQueryHandler(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	by := r.URL.Query().Get("by")
	filter := r.URL.Query().Get("filter")
	if filter == "" {
		WriteJSON(w, http.StatusBadRequest, H{"error": "empty filter"})
		return
	}

	length, err := strconv.Atoi(r.URL.Query().Get("l"))
	if err != nil || length < 1 || length > 500 {
		WriteJSON(w, http.StatusBadRequest, H{"error": "bad length value"})
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("p"))
	if err != nil || page < 1 {
		WriteJSON(w, http.StatusBadRequest, H{"error": "bad page value"})
		return
	}

	recipes, err := sr.Storage.GetFilteredRecipes(filter, sort, by, page, length)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, H{"error": err.Error()})
		return
	}

	WriteJSON(w, http.StatusOK, *recipes)
}

func (sr *RecipesAPIServer) ExportHandler(w http.ResponseWriter, r *http.Request) {
	exportType := r.URL.Query().Get("type")
	if exportType == "" {
		WriteJSON(w, http.StatusBadRequest, H{"error": "empty export type"})
		return
	}

	ct, data, err := sr.Storage.ExportAs(exportType)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError,
			H{"error": fmt.Sprintf("can't export with type '%s'", exportType)},
		)
		return
	}

	w.Header().Set("Content-Disposition", "attachment")
	w.Header().Set("Content-Type", "application/"+ct)
	w.Header().Set("Content-Length", strconv.Itoa(data.Len()))

	data.WriteTo(w)
}

type UsersAPIServer struct {
	Storage  *storage.Storage
	sessions []storage.Session
	mu       sync.RWMutex
}

func NewUsersAPIServer(st *storage.Storage) *UsersAPIServer {
	return &UsersAPIServer{
		st,
		make([]storage.Session, 0),
		sync.RWMutex{},
	}
}

func (sr *UsersAPIServer) UpdateSessions() error {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	sessions, err := sr.Storage.GetSessions()
	if err != nil {
		return err
	}
	sr.sessions = sessions
	return nil
}

func (sr *UsersAPIServer) HasSession(token string) *storage.Session {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	for _, s := range sr.sessions {
		if s.AccessToken == token {
			return &s
		}
	}
	return nil
}

func (sr *UsersAPIServer) LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if len(username) == 0 || len(password) == 0 {
		WriteJSON(w, http.StatusBadRequest, H{"error": "invalid login form"})
		return
	}

	ud, err := sr.Storage.GetUserByName(username)
	if err != nil || ud.HashedPassword != auth.Hash(password) {
		WriteJSON(w, http.StatusUnauthorized, H{"error": "invalid username or password"})
		return
	}

	sr.UpdateSessions()

	WriteJSON(w, http.StatusOK, H{"username": ud.Name, "token": ud.AccessToken, "isAdmin": ud.IsAdmin})
}

func TODOHandler(w http.ResponseWriter, r *http.Request) {
	l.Default.Warn("not implemented")
	w.WriteHeader(http.StatusNotImplemented)
}

type H map[string]any

func WriteJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		l.Default.Errorf("Cant encode JSON: '%v'", err)
	}
}
