package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type exportFoodRecipe struct {
	ItemName     string           `json:"itemName"`
	ResourceName string           `json:"resourceName"`
	Hunger       float32          `json:"hunger"`
	Energy       int              `json:"energy"`
	FEPs         []FoodFEP        `json:"feps"`
	Ingredients  []FoodIngredient `json:"ingredients"`
}

const (
	Default  = "default"
	JSON     = "json"
	Nurgling = "nurgling"
)

func (st *Storage) ExportAs(ex string) (string, bytes.Buffer, error) {
	switch ex {
	case Default:
		b, err := st.exportDefault()
		return "x-sqlite", b, err
	case JSON:
		b, err := st.exportJSON()
		return "json", b, err
	default:
		return "", bytes.Buffer{}, errors.New("wrong export format")
	}
}

func (st *Storage) exportDefault() (bytes.Buffer, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	const exportFileName = "export_cookbook.db"
	var buf bytes.Buffer

	// Remove export file
	_, err := os.Stat(exportFileName)
	if !errors.Is(err, os.ErrNotExist) {
		if err := os.Remove("export_cookbook.db"); err != nil {
			return buf, err
		}
	}
	defer os.Remove(exportFileName)

	// Select recipes
	recipeRows, err := st.db.Queryx("SELECT * FROM recipes")
	if err != nil {
		return buf, err
	}

	var recipeResults []map[string]any

	for recipeRows.Next() {
		row := make(map[string]any)

		if err := recipeRows.MapScan(row); err != nil {
			return buf, err
		}

		recipeResults = append(recipeResults, row)
	}
	recipeRows.Close()

	recipeKeys := make([]string, 0, len(recipeResults[0]))
	for key := range recipeResults[0] {
		recipeKeys = append(recipeKeys, key)
	}

	recipeQuery := fmt.Sprintf(
		"INSERT INTO recipes (%s) VALUES (:%s)",
		strings.Join(recipeKeys, ", "),
		strings.Join(recipeKeys, ", :"),
	)

	// Select ingredients
	ingredientRows, err := st.db.Queryx("SELECT * FROM ingredients")
	if err != nil {
		return buf, err
	}

	var ingredientResults []map[string]any

	for ingredientRows.Next() {
		row := make(map[string]any)

		if err := ingredientRows.MapScan(row); err != nil {
			return buf, err
		}

		ingredientResults = append(ingredientResults, row)
	}
	ingredientRows.Close()

	ingredientKeys := make([]string, 0, len(ingredientResults[0]))
	for key := range ingredientResults[0] {
		ingredientKeys = append(ingredientKeys, key)
	}

	ingredientQuery := fmt.Sprintf(
		"INSERT INTO ingredients (%s) VALUES (:%s)",
		strings.Join(ingredientKeys, ", "),
		strings.Join(ingredientKeys, ", :"),
	)

	// Exported database
	newDB, err := sqlx.Open("sqlite", exportFileName)
	if err != nil {
		return buf, err
	}
	defer newDB.Close()

	// Schema init
	tx, err := newDB.Beginx()
	if err != nil {
		return buf, err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(RECIPES); err != nil {
		return buf, err
	}

	if _, err := tx.Exec(INGREDIENTS); err != nil {
		return buf, err
	}

	// Chunk insert
	const chunkSize int = 100

	for i := 0; i < len(recipeResults); i += chunkSize {
		end := min(i+chunkSize, len(recipeResults))

		chunk := recipeResults[i:end]
		if _, err := tx.NamedQuery(recipeQuery, chunk); err != nil {
			return buf, err
		}
	}

	for i := 0; i < len(ingredientResults); i += chunkSize {
		end := min(i+chunkSize, len(ingredientResults))

		chunk := ingredientResults[i:end]
		if _, err := tx.NamedQuery(ingredientQuery, chunk); err != nil {
			return buf, err
		}
	}

	tx.Commit()

	if err := newDB.Close(); err != nil {
		return buf, err
	}

	exportFile, err := os.Open(exportFileName)
	if err != nil {
		return buf, err
	}
	defer exportFile.Close()

	if _, err := io.Copy(&buf, exportFile); err != nil {
		return buf, err
	}

	return buf, nil
}

func (st *Storage) exportJSON() (bytes.Buffer, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	var buf bytes.Buffer

	recipes, err := st.GetAllRecipes()
	if err != nil {
		return buf, err
	}

	exportRecipes := make([]exportFoodRecipe, len(recipes))

	reverseShortcuts := map[string]string{}
	for k, v := range Shortcuts {
		reverseShortcuts[v] = k
	}

	for i, rc := range recipes {
		newFeps := make([]FoodFEP, len(rc.FEPs))
		for i, f := range rc.FEPs {
			newFeps[i] = FoodFEP{reverseShortcuts[f.Name], f.Value}
		}

		exportRecipes[i] = exportFoodRecipe{
			rc.ItemName,
			rc.ResourceName,
			rc.Hunger,
			rc.Energy,
			newFeps,
			rc.Ingredients,
		}
	}

	encoded, err := json.Marshal(exportRecipes)
	if err != nil {
		return buf, err
	}

	buf.Write(encoded)

	return buf, nil
}
