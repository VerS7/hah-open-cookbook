package filter

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

type Percentage float64

func (p Percentage) String() string {
	return fmt.Sprintf("%.2f%%", float64(p)*100)
}

func (p Percentage) Value() float64 {
	return float64(p)
}

type Condition struct {
	Name     string
	Operator string
	Value    any
}

type Expectations map[string][]string

func Parse(input string) ([]Condition, error) {
	if input == "" {
		return []Condition{}, nil
	}

	var conditions []Condition
	segments := strings.SplitSeq(input, ";")

	for segment := range segments {
		segment = strings.TrimSpace(segment)
		if segment == "" {
			continue
		}

		operator := ""
		operatorPos := -1
		operators := []string{"<=", ">=", "!=", "=", "<", ">"}

		for _, op := range operators {
			if pos := strings.Index(segment, op); pos != -1 {
				operator = op
				operatorPos = pos
				break
			}
		}

		if operator == "" {
			return nil, fmt.Errorf("operator not found at segment: '%s'", segment)
		}

		name := strings.TrimSpace(segment[:operatorPos])
		if name == "" {
			return nil, fmt.Errorf("arg name not found at segment: '%s'", segment)
		}

		valueStr := strings.TrimSpace(segment[operatorPos+len(operator):])

		value, err := parseValue(valueStr)
		if err != nil {
			return nil, fmt.Errorf("cant parse value for '%s': '%s'", name, err)
		}

		conditions = append(conditions, Condition{
			Name:     name,
			Operator: operator,
			Value:    value,
		})
	}

	return conditions, nil
}

func parseValue(valueStr string) (any, error) {
	valueStr = strings.TrimSpace(valueStr)
	if valueStr == "" {
		return "", nil
	}

	if len(valueStr) > 0 && valueStr[0] == '"' {
		return parseStringOrList(valueStr)
	}

	if len(valueStr) > 0 && valueStr[len(valueStr)-1] == '%' {
		numStr := strings.TrimSpace(valueStr[:len(valueStr)-1])
		val, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return nil, fmt.Errorf("wrong percentage format: '%s'", valueStr)
		}
		return Percentage(val / 100.0), nil
	}

	valInt, err := strconv.ParseInt(valueStr, 10, 64)
	if err == nil {
		return valInt, nil
	}

	valFloat, err := strconv.ParseFloat(valueStr, 64)
	if err == nil {
		return valFloat, nil
	}

	return valueStr, nil
}

func parseStringOrList(listStr string) (any, error) {
	var result []string
	var current strings.Builder
	inQuotes := false
	escape := false

	listStr = strings.TrimSpace(listStr)

	if len(listStr) == 0 || listStr[0] != '"' {
		return listStr, nil
	}

	for i, ch := range listStr {
		if escape {
			current.WriteRune(ch)
			escape = false
			continue
		}

		if ch == '\\' {
			escape = true
			continue
		}

		if ch == '"' {
			if inQuotes {
				str := current.String()
				result = append(result, str)
				current.Reset()
			}
			inQuotes = !inQuotes
			continue
		}

		if ch == ',' && !inQuotes {
			continue
		}

		if inQuotes || !unicode.IsSpace(ch) {
			current.WriteRune(ch)
		}

		if i == len(listStr)-1 && inQuotes {
			str := current.String()
			if str != "" {
				result = append(result, str)
			}
		}
	}

	if len(result) == 1 {
		return result[0], nil
	}

	if len(result) == 0 {
		return "", nil
	}

	return result, nil
}

func ValidateConditions(conditions []Condition, expectedArgs Expectations) error {
	for _, cond := range conditions {
		allowedTypes, exists := expectedArgs[cond.Name]
		if !exists {
			return fmt.Errorf("unknown arg: '%s'", cond.Name)
		}

		valueType := getValueType(cond.Value)

		typeAllowed := false
		for _, allowedType := range allowedTypes {
			if valueType == allowedType {
				typeAllowed = true
				break
			}

			if valueType == "int" && (allowedType == "float" || allowedType == "percent") {
				typeAllowed = true
				break
			}

			if valueType == "string" && allowedType == "string_list" {
				typeAllowed = true
				break
			}

			if valueType == "float" && allowedType == "percent" {
				if f, ok := cond.Value.(float64); ok && f >= 0 && f <= 1 {
					typeAllowed = true
					break
				}
			}
		}

		if !typeAllowed {
			return fmt.Errorf("arg '%s' has type '%s', but expected one of: %v",
				cond.Name, valueType, allowedTypes)
		}
	}

	return nil
}

func getValueType(value any) string {
	switch value.(type) {
	case int64:
		return "int"
	case float64:
		return "float"
	case string:
		return "string"
	case Percentage:
		return "percent"
	case []string:
		return "string_list"
	default:
		return "unknown"
	}
}

func CheckCondition(cond Condition, value any) bool {
	switch cond.Operator {
	case "=":
		return equals(value, cond.Value)
	case "!=":
		return !equals(value, cond.Value)
	case "<":
		return lessThan(value, cond.Value)
	case ">":
		return greaterThan(value, cond.Value)
	case "<=":
		return lessThanOrEqual(value, cond.Value)
	case ">=":
		return greaterThanOrEqual(value, cond.Value)
	default:
		return false
	}
}

func equals(a, b any) bool {
	switch aVal := a.(type) {
	case int64:
		switch bVal := b.(type) {
		case int64:
			return aVal == bVal
		case float64:
			return float64(aVal) == bVal
		case Percentage:
			return float64(aVal) == float64(bVal)
		}
	case float64:
		switch bVal := b.(type) {
		case int64:
			return aVal == float64(bVal)
		case float64:
			return aVal == bVal
		case Percentage:
			return aVal == float64(bVal)
		}
	case Percentage:
		switch bVal := b.(type) {
		case int64:
			return float64(aVal) == float64(bVal)
		case float64:
			return float64(aVal) == bVal
		case Percentage:
			return aVal == bVal
		}
	}

	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func lessThan(a, b any) bool {
	switch aVal := a.(type) {
	case int64:
		switch bVal := b.(type) {
		case int64:
			return aVal < bVal
		case float64:
			return float64(aVal) < bVal
		case Percentage:
			return float64(aVal) < float64(bVal)
		}
	case float64:
		switch bVal := b.(type) {
		case int64:
			return aVal < float64(bVal)
		case float64:
			return aVal < bVal
		case Percentage:
			return aVal < float64(bVal)
		}
	case Percentage:
		switch bVal := b.(type) {
		case int64:
			return float64(aVal) < float64(bVal)
		case float64:
			return float64(aVal) < bVal
		case Percentage:
			return aVal < bVal
		}
	}
	return false
}

func greaterThan(a, b any) bool {
	return !lessThanOrEqual(a, b)
}

func lessThanOrEqual(a, b any) bool {
	return equals(a, b) || lessThan(a, b)
}

func greaterThanOrEqual(a, b any) bool {
	return equals(a, b) || greaterThan(a, b)
}

func NewExpectations() Expectations {
	return make(Expectations)
}

func (e Expectations) SetTypes(argName string, types ...string) {
	e[argName] = types
}

func (e Expectations) AddType(argName, typeName string) {
	if _, exists := e[argName]; !exists {
		e[argName] = []string{}
	}

	if slices.Contains(e[argName], typeName) {
		return
	}

	e[argName] = append(e[argName], typeName)
}

func IsPercentage(value any) bool {
	_, ok := value.(Percentage)

	return ok
}

func GetPercentageValue(value any) (float64, bool) {
	switch v := value.(type) {
	case Percentage:
		return float64(v), true
	case float64:
		if v >= 0 && v <= 1 {
			return v, true
		}
	}
	return 0, false
}

func GetIntValue(value any) (int64, bool) {
	switch v := value.(type) {
	case int64:
		return v, true
	case float64:
		if v == float64(int64(v)) {
			return int64(v), true
		}
	case Percentage:
		return int64(v), true
	}
	return 0, false
}

func GetFloatValue(value any) (float64, bool) {
	switch v := value.(type) {
	case int64:
		return float64(v), true
	case float64:
		return v, true
	case Percentage:
		return float64(v), true
	}
	return 0, false
}

func GetStringValue(value any) (string, bool) {
	switch v := value.(type) {
	case string:
		return v, true
	case []string:
		if len(v) == 1 {
			return v[0], true
		}
		return strings.Join(v, ","), true
	default:
		return fmt.Sprintf("%v", v), true
	}
}

func GetStringListValue(value any) ([]string, bool) {
	switch v := value.(type) {
	case []string:
		return v, true
	case string:
		return []string{v}, true
	default:
		return []string{fmt.Sprintf("%v", v)}, true
	}
}
