package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"open-hah-cookbook/internal/auth"
	"open-hah-cookbook/internal/filter"
)

type FoodFEP struct {
	Name  string  `json:"name"`
	Value float32 `json:"value"`
}

type FoodIngredient struct {
	Name       string `json:"name"`
	Percentage int    `json:"percentage"`
}

type FoodRecipe struct {
	Id           int              `json:"id"`
	ItemName     string           `json:"itemName"`
	ResourceName string           `json:"resourceName"`
	Hunger       float32          `json:"hunger"`
	Energy       int              `json:"energy"`
	FEPs         []FoodFEP        `json:"feps"`
	Ingredients  []FoodIngredient `json:"ingredients"`
	CreatedAt    *time.Time       `json:"created_at,omitempty"`
}

type FoodRecipeTimestamps struct {
	First time.Time `json:"first"`
	Last  time.Time `json:"last"`
}

type FoodRecipePage struct {
	Recipes    []FoodRecipe         `json:"recipes"`
	Count      int                  `json:"count"`
	Total      int                  `json:"total"`
	Page       int                  `json:"page"`
	Pages      int                  `json:"pages"`
	Timestamps FoodRecipeTimestamps `json:"timestamps"`
}

var Shortcuts = map[string]string{
	"Strength +1":     "str1",
	"Strength +2":     "str2",
	"Agility +1":      "agi1",
	"Agility +2":      "agi2",
	"Intelligence +1": "int1",
	"Intelligence +2": "int2",
	"Constitution +1": "con1",
	"Constitution +2": "con2",
	"Perception +1":   "prc1",
	"Perception +2":   "prc2",
	"Charisma +1":     "csm1",
	"Charisma +2":     "csm2",
	"Dexterity +1":    "dex1",
	"Dexterity +2":    "dex2",
	"Will +1":         "wil1",
	"Will +2":         "wil2",
	"Psyche +1":       "psy1",
	"Psyche +2":       "psy2",
}

var CombinedShortcuts = map[string][2]string{
	"str": {"str1", "str2"},
	"agi": {"agi1", "agi2"},
	"int": {"int1", "int2"},
	"con": {"con1", "con2"},
	"prc": {"prc1", "prc2"},
	"csm": {"csm1", "csm2"},
	"dex": {"dex1", "dex2"},
	"wil": {"wil1", "wil2"},
	"psy": {"psy1", "psy2"},
}

func (f *FoodRecipe) AsMap() map[string]any {
	recipe := make(map[string]any)

	// Base
	recipe["name"] = f.ItemName
	recipe["resource"] = f.ResourceName
	recipe["hunger"] = f.Hunger
	recipe["energy"] = f.Energy
	recipe["ingredients"] = f.Ingredients

	// FEPs
	for _, fep := range f.FEPs {
		recipe[Shortcuts[fep.Name]] = fep.Value
	}

	return recipe
}

func (f *FoodRecipe) Hash() string {
	hashable := make([]string, 4)

	hashable[0] = f.ItemName
	hashable[1] = f.ResourceName
	hashable[2] = fmt.Sprintf("%f", f.Hunger)
	hashable[3] = fmt.Sprintf("%d", f.Energy)

	for _, ingr := range f.Ingredients {
		hashable = append(hashable, ingr.Name)
		hashable = append(hashable, fmt.Sprintf("%d", ingr.Percentage))
	}

	for _, fep := range f.FEPs {
		hashable = append(hashable, fep.Name)
		hashable = append(hashable, fmt.Sprintf("%f", fep.Value))
	}

	return auth.HashSum(hashable...)
}

func (st *Storage) AddRecipe(recipe FoodRecipe) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	rm := recipe.AsMap()

	params := make(map[string]any)
	names := make([]string, 0)
	placeholders := make([]string, 0)

	params["hash"] = recipe.Hash()
	names = append(names, "hash")
	placeholders = append(placeholders, ":hash")

	for k, v := range rm {
		switch k {
		case "ingredients":
			continue
		case "feps":
			feps, _ := v.([]FoodFEP)
			for _, f := range feps {
				params[f.Name] = f.Value
				names = append(names, f.Name)
			}
		default:
			params[k] = v
			names = append(names, k)
			placeholders = append(placeholders, ":"+k)
		}
	}

	tx, _ := st.db.Beginx()
	defer tx.Rollback()
	{
		query := fmt.Sprintf(
			"INSERT OR IGNORE INTO recipes (%s) VALUES (%s)",
			strings.Join(names, ", "),
			strings.Join(placeholders, ", "),
		)

		r, err := tx.NamedExec(query, params)
		if err != nil {
			return err
		}

		rowsAffected, err := r.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected > 0 {
			id, err := r.LastInsertId()
			if err != nil {
				return err
			}

			ingrs, ok := rm["ingredients"].([]FoodIngredient)
			if ok && len(ingrs) > 0 {
				for _, ing := range recipe.Ingredients {
					if _, err := tx.Exec("INSERT INTO ingredients (name, rate, recipe) VALUES (?, ?, ?)",
						ing.Name,
						ing.Percentage,
						id,
					); err != nil {
						return err
					}
				}
			}
		}

	}
	return tx.Commit()
}

func (st *Storage) RemoveRecipe(id int) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	tx, _ := st.db.Beginx()
	defer tx.Rollback()
	{
		if _, err := tx.Exec("DELETE FROM recipes WHERE id = ?", id); err != nil {
			return err
		}

		if _, err := tx.Exec("DELETE FROM ingredients WHERE recipe = ?", id); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (st *Storage) GetTimestamps() (*FoodRecipeTimestamps, error) {
	var first, last string
	var timestamps FoodRecipeTimestamps
	if err := st.db.QueryRowx("SELECT MIN(created_at) AS first, MAX(created_at) AS last FROM recipes").Scan(&first, &last); err != nil {
		return nil, err
	}

	firstTime, err := time.Parse("2006-01-02 15:04:05", first)
	lastTime, err := time.Parse("2006-01-02 15:04:05", last)
	if err != nil {
		return nil, err
	}

	timestamps.First = firstTime
	timestamps.Last = lastTime

	return &timestamps, nil
}

func (st *Storage) GetQueriedTimestamps(q *Query) (*FoodRecipeTimestamps, error) {
	var first, last string
	var timestamps FoodRecipeTimestamps

	query, values := q.AsTimestamp()

	if err := st.db.QueryRowx(query, values...).Scan(&first, &last); err != nil {
		return nil, err
	}

	firstTime, err := time.Parse("2006-01-02 15:04:05", first)
	lastTime, err := time.Parse("2006-01-02 15:04:05", last)
	if err != nil {
		return nil, err
	}

	timestamps.First = firstTime
	timestamps.Last = lastTime

	return &timestamps, nil
}

func (st *Storage) GetAllRecipes() ([]FoodRecipe, error) {
	query := `
		SELECT 
		r.id,
		r.created_at,
		r.energy,
		r.hunger,
		r.name,
		r.resource,

		(
		COALESCE(r.str1, 0) + COALESCE(r.str2, 0) + 
		COALESCE(r.agi1, 0) + COALESCE(r.agi2, 0) +
		COALESCE(r.int1, 0) + COALESCE(r.int2, 0) +
		COALESCE(r.con1, 0) + COALESCE(r.con2, 0) +
		COALESCE(r.prc1, 0) + COALESCE(r.prc2, 0) +
		COALESCE(r.csm1, 0) + COALESCE(r.csm2, 0) +
		COALESCE(r.dex1, 0) + COALESCE(r.dex2, 0) +
		COALESCE(r.wil1, 0) + COALESCE(r.wil2, 0) +
		COALESCE(r.psy1, 0) + COALESCE(r.psy2, 0)
		) as total,

		CASE
			WHEN COUNT(i.id) > 0 THEN
				json_group_array(
					json_object(
						'name', i.name,
						'percentage', i.rate
					)
				)
			ELSE json_array()
		END as ingredients,

		COALESCE(r.str1, 0) as str1,
		COALESCE(r.str2, 0) as str2,
		COALESCE(r.agi1, 0) as agi1,
		COALESCE(r.agi2, 0) as agi2,
		COALESCE(r.int1, 0) as int1,
		COALESCE(r.int2, 0) as int2,
		COALESCE(r.con1, 0) as con1,
		COALESCE(r.con2, 0) as con2,
		COALESCE(r.prc1, 0) as prc1,
		COALESCE(r.prc2, 0) as prc2,
		COALESCE(r.csm1, 0) as csm1,
		COALESCE(r.csm2, 0) as csm2,
		COALESCE(r.dex1, 0) as dex1,
		COALESCE(r.dex2, 0) as dex2,
		COALESCE(r.wil1, 0) as wil1,
		COALESCE(r.wil2, 0) as wil2,
		COALESCE(r.psy1, 0) as psy1,
		COALESCE(r.psy2, 0) as psy2
		
		FROM recipes r
		LEFT JOIN ingredients i ON r.id = i.recipe
		GROUP BY r.id, r.created_at, r.energy, r.hunger, r.name, r.resource,
				r.str1, r.str2, r.agi1, r.agi2, r.int1, r.int2,
				r.con1, r.con2, r.prc1, r.prc2, r.csm1, r.csm2,
				r.dex1, r.dex2, r.wil1, r.wil2, r.psy1, r.psy2
		ORDER BY r.id`

	result, err := st.db.Query(query)
	if err != nil {
		return nil, err
	}

	columns, err := result.Columns()
	if err != nil {
		return nil, err
	}

	const fepAlign = 8
	var count int
	recipes := make([]FoodRecipe, 0)

	for result.Next() {
		count += 1

		var id int
		var energy int
		var createdAt time.Time
		var hunger float32
		var total float32
		var name string
		var resource string
		var ingredients string

		feps := make([]any, len(columns)-fepAlign)
		fepsPtrs := make([]any, len(columns)-fepAlign)
		for i := range columns[fepAlign:] {
			fepsPtrs[i] = &feps[i]
		}

		scanPtrs := append([]any{&id, &createdAt, &energy, &hunger, &name, &resource, &total, &ingredients}, fepsPtrs...)
		err := result.Scan(scanPtrs...)
		if err != nil {
			return nil, err
		}

		var fi []FoodIngredient
		if err := json.Unmarshal([]byte(ingredients), &fi); err != nil {
			return nil, err
		}

		var rc FoodRecipe
		rc.Id = id
		rc.ItemName = name
		rc.ResourceName = resource
		rc.Energy = energy
		rc.Hunger = hunger
		rc.Ingredients = fi
		rc.CreatedAt = &createdAt

		fp := make([]FoodFEP, 0)
		for i, col := range columns[fepAlign:] {
			v := parseFepValue(feps[i])
			if v > 0 {
				fp = append(fp, FoodFEP{col, v})
			}
		}

		rc.FEPs = fp

		recipes = append(recipes, rc)
	}

	return recipes, nil
}

func (st *Storage) GetFilteredRecipes(filterString string, sort string, by string, page, pageSize int) (*FoodRecipePage, error) {
	var sortDirection string
	var sortBy string

	conditions, err := parseFilter(filterString)
	if err != nil {
		return nil, err
	}

	// Sort order
	switch sort {
	case "asc", "desc":
		sortDirection = sort
	default:
		return nil, fmt.Errorf("sorting: '%s' unavailable", sort)
	}

	nonZeroBy := make(map[string]string)

	// Sort category
	category := make(map[string]string)
	category["name"] = "name"
	category["energy"] = "energy"
	category["hunger"] = "hunger"
	category["total"] = "total"
	category["feph"] = "feph"

	for _, v := range Shortcuts {
		category[v] = v
		nonZeroBy[v] = fmt.Sprintf("COALESCE(%s, 0)", v)
	}

	for fep, values := range CombinedShortcuts {
		expr := fmt.Sprintf("COALESCE(%s, 0) + COALESCE(%s, 0)", values[0], values[1])
		category[fep] = expr
		nonZeroBy[fep] = expr
	}

	if expr, ok := category[by]; ok {
		sortBy = expr
	} else {
		return nil, fmt.Errorf("sorting by category: '%s' unavailable", by)
	}

	// Base query
	query, err := ConstructQuery(conditions)
	if err != nil {
		return nil, err
	}

	// When sorting by FEP key, skip rows where that key equals 0
	if expr, ok := nonZeroBy[by]; ok {
		base := query.builder.String()
		sb := strings.Builder{}
		sb.WriteString(base)
		if strings.Contains(strings.ToUpper(base), " WHERE ") {
			fmt.Fprintf(&sb, " AND (%s) != 0 ", expr)
		} else {
			fmt.Fprintf(&sb, " WHERE (%s) != 0 ", expr)
		}
		query.builder = sb
	}

	// Recipes count
	countQuery, countValues := query.AsCount()
	var total int
	if err := st.db.Get(&total, countQuery, countValues...); err != nil {
		return nil, err
	}

	// Paginated and sorted data
	pageQuery, pageValues := query.AsPageSorted(sortDirection, sortBy, page, pageSize)
	result, err := st.db.Query(pageQuery, pageValues...)
	if err != nil {
		return nil, err
	}

	// Timestamps
	var timestamps *FoodRecipeTimestamps
	t, err := st.GetQueriedTimestamps(&query)
	if err != nil {
		fmt.Println(err)

		t, err := st.GetTimestamps()
		if err != nil {
			return nil, err
		}
		timestamps = t
	} else {
		timestamps = t
	}

	columns, err := result.Columns()
	if err != nil {
		return nil, err
	}

	const fepAlign = 9
	var count int
	recipes := make([]FoodRecipe, 0)

	for result.Next() {
		count += 1

		var id int
		var energy int
		var createdAt time.Time
		var hunger float32
		var total float32
		var feph float32
		var name string
		var resource string
		var ingredients string

		feps := make([]any, len(columns)-fepAlign)
		fepsPtrs := make([]any, len(columns)-fepAlign)
		for i := range columns[fepAlign:] {
			fepsPtrs[i] = &feps[i]
		}

		scanPtrs := append([]any{&id, &createdAt, &energy, &hunger, &name, &resource, &total, &feph, &ingredients}, fepsPtrs...)
		err := result.Scan(scanPtrs...)
		if err != nil {
			return nil, err
		}

		var fi []FoodIngredient
		if err := json.Unmarshal([]byte(ingredients), &fi); err != nil {
			return nil, err
		}

		var rc FoodRecipe
		rc.Id = id
		rc.ItemName = name
		rc.ResourceName = resource
		rc.Energy = energy
		rc.Hunger = hunger
		rc.Ingredients = fi
		rc.CreatedAt = &createdAt

		fp := make([]FoodFEP, 0)
		for i, col := range columns[fepAlign:] {
			v := parseFepValue(feps[i])
			if v > 0 {
				fp = append(fp, FoodFEP{col, v})
			}
		}

		rc.FEPs = fp

		recipes = append(recipes, rc)
	}
	return &FoodRecipePage{
		Recipes:    recipes,
		Total:      total,
		Count:      count,
		Page:       page,
		Pages:      int(math.Ceil(float64(total) / float64(pageSize))),
		Timestamps: *timestamps,
	}, nil
}

type Query struct {
	builder strings.Builder
	values  []any
}

func ConstructQuery(conditions []filter.Condition) (Query, error) {
	var q Query
	sb := strings.Builder{}

	sb.WriteString(`
	WITH summary AS (
		SELECT 
		r.id,
		r.created_at,
		r.energy,
		r.hunger,
		r.name,
		r.resource,

		(
		COALESCE(r.str1, 0) + COALESCE(r.str2, 0) + 
		COALESCE(r.agi1, 0) + COALESCE(r.agi2, 0) +
		COALESCE(r.int1, 0) + COALESCE(r.int2, 0) +
		COALESCE(r.con1, 0) + COALESCE(r.con2, 0) +
		COALESCE(r.prc1, 0) + COALESCE(r.prc2, 0) +
		COALESCE(r.csm1, 0) + COALESCE(r.csm2, 0) +
		COALESCE(r.dex1, 0) + COALESCE(r.dex2, 0) +
		COALESCE(r.wil1, 0) + COALESCE(r.wil2, 0) +
		COALESCE(r.psy1, 0) + COALESCE(r.psy2, 0)
		) as total,

		(
		COALESCE(r.str1, 0) + COALESCE(r.str2, 0) + 
		COALESCE(r.agi1, 0) + COALESCE(r.agi2, 0) +
		COALESCE(r.int1, 0) + COALESCE(r.int2, 0) +
		COALESCE(r.con1, 0) + COALESCE(r.con2, 0) +
		COALESCE(r.prc1, 0) + COALESCE(r.prc2, 0) +
		COALESCE(r.csm1, 0) + COALESCE(r.csm2, 0) +
		COALESCE(r.dex1, 0) + COALESCE(r.dex2, 0) +
		COALESCE(r.wil1, 0) + COALESCE(r.wil2, 0) +
		COALESCE(r.psy1, 0) + COALESCE(r.psy2, 0)
		) / COALESCE(r.hunger, 1) as feph,

		CASE
			WHEN COUNT(i.id) > 0 THEN
				json_group_array(
						json_object(
							'name', i.name,
							'percentage', i.rate
						)
					)
			ELSE json_array()
		END as ingredients,

		COALESCE(r.str1, 0) as str1,
		COALESCE(r.str2, 0) as str2,
		COALESCE(r.agi1, 0) as agi1,
		COALESCE(r.agi2, 0) as agi2,
		COALESCE(r.int1, 0) as int1,
		COALESCE(r.int2, 0) as int2,
		COALESCE(r.con1, 0) as con1,
		COALESCE(r.con2, 0) as con2,
		COALESCE(r.prc1, 0) as prc1,
		COALESCE(r.prc2, 0) as prc2,
		COALESCE(r.csm1, 0) as csm1,
		COALESCE(r.csm2, 0) as csm2,
		COALESCE(r.dex1, 0) as dex1,
		COALESCE(r.dex2, 0) as dex2,
		COALESCE(r.wil1, 0) as wil1,
		COALESCE(r.wil2, 0) as wil2,
		COALESCE(r.psy1, 0) as psy1,
		COALESCE(r.psy2, 0) as psy2
		
	FROM recipes r
	LEFT JOIN ingredients i ON r.id = i.recipe
	GROUP BY r.id, r.created_at, r.energy, r.hunger, r.name, r.resource,
			r.str1, r.str2, r.agi1, r.agi2, r.int1, r.int2,
			r.con1, r.con2, r.prc1, r.prc2, r.csm1, r.csm2,
			r.dex1, r.dex2, r.wil1, r.wil2, r.psy1, r.psy2
	ORDER BY r.id
	) `)
	sb.WriteString("SELECT * FROM summary ")

	values := make([]any, 0)

	if len(conditions) > 0 {
		sb.WriteString("WHERE ")
	}

	for i, c := range conditions {
		switch c.Name {
		case "from":
			if c.Operator != "=" && c.Operator != "!=" {
				return q, fmt.Errorf("unexpected operator for 'from' condition: '%s'", c.Operator)
			}
			list, ok := filter.GetStringListValue(c.Value)
			if ok { // If value is list
				for j, el := range list {
					values = append(values, el)
					if c.Operator == "!=" {
						sb.WriteString("ingredients NOT LIKE CONCAT('%', ?, '%') ")
					} else {
						sb.WriteString("ingredients LIKE CONCAT('%', ?, '%') ")
					}
					if j < len(list)-1 {
						log.Println(j, len(list))
						sb.WriteString("AND ")
					}
				}
			} else { // If value is string
				val, ok := filter.GetStringValue(c.Value)
				values = append(values, val)
				if ok {
					if c.Operator == "!=" {
						sb.WriteString("AND ingredients NOT LIKE CONCAT('%', ?, '%') ")
					} else {
						sb.WriteString("AND ingredients LIKE CONCAT('%', ?, '%') ")
					}
				}
			}

		case "name":
			if c.Operator != "=" && c.Operator != "!=" {
				return q, fmt.Errorf("unexpected operator for 'name' condition: '%s'", c.Operator)
			}
			values = append(values, c.Value)
			if c.Operator == "!=" {
				sb.WriteString("name NOT LIKE CONCAT('%', ?, '%') ")
			} else {
				sb.WriteString("name LIKE CONCAT('%', ?, '%') ")
			}

		default:
			if filter.IsPercentage(c.Value) {
				fmt.Fprintf(&sb, "IFNULL(%s, 0) != 0 AND ", c.Name)
				v, _ := c.Value.(filter.Percentage)
				fmt.Fprintf(&sb, "(%s / total) * 100 %s ?", c.Name, c.Operator)
				values = append(values, v.Value()*100)
			} else {
				values = append(values, c.Value)
				fmt.Fprintf(&sb, "%s %s ? ", c.Name, c.Operator)
			}
		}

		if i < len(conditions)-1 {
			sb.WriteString("AND ")
		}
	}

	q.builder = sb
	q.values = values

	return q, nil
}

func (q *Query) AsPageSorted(sort, by string, page, pageSize int) (string, []any) {
	sb := strings.Builder{}

	sb.WriteString(q.builder.String())
	fmt.Fprintf(&sb, " ORDER BY %s %s ", by, strings.ToUpper(sort))
	sb.WriteString("LIMIT ? OFFSET (? - 1) * ?")
	values := append(q.values, pageSize, page, pageSize)

	return sb.String(), values
}

func (q *Query) AsCount() (string, []any) {
	sb := strings.Builder{}
	sb.WriteString("WITH full_summary AS (")
	sb.WriteString(q.builder.String())
	sb.WriteString(") SELECT COUNT(*) FROM full_summary")

	return sb.String(), q.values
}

func (q *Query) AsTimestamp() (string, []any) {
	sb := strings.Builder{}
	sb.WriteString("WITH full_summary AS (")
	sb.WriteString(q.builder.String())
	sb.WriteString(") SELECT MIN(created_at) AS first, MAX(created_at) AS last FROM full_summary")

	return sb.String(), q.values
}

func parseFepValue(fep any) float32 {
	var fepValue float32 = 0

	if fep != nil {
		switch fep := fep.(type) {
		case int32:
			fepValue = float32(fep)
		case int64:
			fepValue = float32(fep)
		case float64:
			fepValue = float32(fep)
		default:
			fepValue = fep.(float32)
		}
	}

	return fepValue
}

func parseFilter(filterString string) ([]filter.Condition, error) {
	ex := filter.NewExpectations()
	for _, s := range Shortcuts {
		ex.SetTypes(s, "int", "float", "percent")
	}
	ex.SetTypes("hunger", "float")
	ex.SetTypes("total", "int", "float")
	ex.SetTypes("energy", "int")
	ex.SetTypes("name", "string")
	ex.SetTypes("from", "string_list")

	cn, err := filter.Parse(filterString)
	if err != nil {
		return nil, err
	}

	if err := filter.ValidateConditions(cn, ex); err != nil {
		return nil, err
	}
	return cn, nil
}
