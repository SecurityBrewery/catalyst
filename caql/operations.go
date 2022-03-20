package caql

import (
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Logical operators https://www.arangodb.com/docs/3.7/aql/operators.html#logical-operators

func or(left, right any) any {
	if toBool(left) {
		return left
	}

	return right
}

func and(left, right any) any {
	if !toBool(left) {
		return left
	}

	return right
}

func toBool(i any) bool {
	switch v := i.(type) {
	case nil:
		return false
	case bool:
		return v
	case int:
		return v != 0
	case float64:
		return v != 0
	case string:
		return v != ""
	case []any:
		return true
	case map[string]any:
		return true
	default:
		panic("bool conversion failed")
	}
}

// Arithmetic operators https://www.arangodb.com/docs/3.7/aql/operators.html#arithmetic-operators

func plus(left, right any) float64 {
	return toNumber(left) + toNumber(right)
}

func minus(left, right any) float64 {
	return toNumber(left) - toNumber(right)
}

func times(left, right any) float64 {
	return round(toNumber(left) * toNumber(right))
}

func round(r float64) float64 {
	return math.Round(r*100000) / 100000
}

func div(left, right any) float64 {
	b := toNumber(right)
	if b == 0 {
		return 0
	}

	return round(toNumber(left) / b)
}

func mod(left, right any) float64 {
	return math.Mod(toNumber(left), toNumber(right))
}

func toNumber(i any) float64 {
	switch v := i.(type) {
	case nil:
		return 0
	case bool:
		if v {
			return 1
		}

		return 0
	case float64:
		switch {
		case math.IsNaN(v):
			return 0
		case math.IsInf(v, 0):
			return 0
		}

		return v
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return 0
		}

		return f
	case []any:
		if len(v) == 0 {
			return 0
		}
		if len(v) == 1 {
			return toNumber(v[0])
		}

		return 0
	case map[string]any:
		return 0
	default:
		panic("number conversion error")
	}
}

// Logical operators https://www.arangodb.com/docs/3.7/aql/operators.html#logical-operators
// Order https://www.arangodb.com/docs/3.7/aql/fundamentals-type-value-order.html

func eq(left, right any) bool {
	leftV, rightV := typeValue(left), typeValue(right)
	if leftV != rightV {
		return false
	}
	switch l := left.(type) {
	case nil:
		return true
	case bool, float64, string:
		return left == right
	case []any:
		ra := right.([]any)
		max := len(l)
		if len(ra) > max {
			max = len(ra)
		}
		for i := 0; i < max; i++ {
			var li any
			var rai any
			if len(l) > i {
				li = l[i]
			}
			if len(ra) > i {
				rai = ra[i]
			}

			if !eq(li, rai) {
				return false
			}
		}

		return true
	case map[string]any:
		ro := right.(map[string]any)

		for _, key := range keys(l, ro) {
			var li any
			var rai any
			if lv, ok := l[key]; ok {
				li = lv
			}
			if rv, ok := ro[key]; ok {
				rai = rv
			}

			if !eq(li, rai) {
				return false
			}
		}

		return true
	default:
		panic("unknown type")
	}
}

func ne(left, right any) bool {
	return !eq(left, right)
}

func lt(left, right any) bool {
	leftV, rightV := typeValue(left), typeValue(right)
	if leftV != rightV {
		return leftV < rightV
	}
	switch l := left.(type) {
	case nil:
		return false
	case bool:
		return toNumber(l) < toNumber(right)
	case int:
		return l < right.(int)
	case float64:
		return l < right.(float64)
	case string:
		return l < right.(string)
	case []any:
		ra := right.([]any)
		max := len(l)
		if len(ra) > max {
			max = len(ra)
		}
		for i := 0; i < max; i++ {
			var li any
			var rai any
			if len(l) > i {
				li = l[i]
			}
			if len(ra) > i {
				rai = ra[i]
			}

			if !eq(li, rai) {
				return lt(li, rai)
			}
		}

		return false
	case map[string]any:
		ro := right.(map[string]any)

		for _, key := range keys(l, ro) {
			var li any
			var rai any
			if lv, ok := l[key]; ok {
				li = lv
			}
			if rv, ok := ro[key]; ok {
				rai = rv
			}

			if !eq(li, rai) {
				return lt(li, rai)
			}
		}

		return false
	default:
		panic("unknown type")
	}
}

func keys(l map[string]any, ro map[string]any) []string {
	var keys []string
	seen := map[string]bool{}
	for _, a := range []map[string]any{l, ro} {
		for k := range a {
			if _, ok := seen[k]; !ok {
				seen[k] = true
				keys = append(keys, k)
			}
		}
	}
	sort.Strings(keys)

	return keys
}

func gt(left, right any) bool {
	leftV, rightV := typeValue(left), typeValue(right)
	if leftV != rightV {
		return leftV > rightV
	}
	switch l := left.(type) {
	case nil:
		return false
	case bool:
		return toNumber(l) > toNumber(right)
	case int:
		return l > right.(int)
	case float64:
		return l > right.(float64)
	case string:
		return l > right.(string)
	case []any:
		ra := right.([]any)
		max := len(l)
		if len(ra) > max {
			max = len(ra)
		}
		for i := 0; i < max; i++ {
			var li any
			var rai any
			if len(l) > i {
				li = l[i]
			}
			if len(ra) > i {
				rai = ra[i]
			}

			if !eq(li, rai) {
				return gt(li, rai)
			}
		}

		return false
	case map[string]any:
		ro := right.(map[string]any)

		for _, key := range keys(l, ro) {
			var li any
			var rai any
			if lv, ok := l[key]; ok {
				li = lv
			}
			if rv, ok := ro[key]; ok {
				rai = rv
			}

			if !eq(li, rai) {
				return gt(li, rai)
			}
		}

		return false
	default:
		panic("unknown type")
	}
}

func le(left, right any) bool {
	leftV, rightV := typeValue(left), typeValue(right)
	if leftV != rightV {
		return leftV <= rightV
	}
	switch l := left.(type) {
	case nil:
		return false
	case bool:
		return toNumber(l) <= toNumber(right)
	case int:
		return l <= right.(int)
	case float64:
		return l <= right.(float64)
	case string:
		return l <= right.(string)
	case []any:
		ra := right.([]any)
		max := len(l)
		if len(ra) > max {
			max = len(ra)
		}
		for i := 0; i < max; i++ {
			var li any
			var rai any
			if len(l) > i {
				li = l[i]
			}
			if len(ra) > i {
				rai = ra[i]
			}

			if !eq(li, rai) {
				return le(li, rai)
			}
		}

		return true
	case map[string]any:
		ro := right.(map[string]any)

		for _, key := range keys(l, ro) {
			var li any
			var rai any
			if lv, ok := l[key]; ok {
				li = lv
			}
			if rv, ok := ro[key]; ok {
				rai = rv
			}

			if !eq(li, rai) {
				return lt(li, rai)
			}
		}

		return true
	default:
		panic("unknown type")
	}
}

func ge(left, right any) bool {
	leftV, rightV := typeValue(left), typeValue(right)
	if leftV != rightV {
		return leftV >= rightV
	}
	switch l := left.(type) {
	case nil:
		return false
	case bool:
		return toNumber(l) >= toNumber(right)
	case int:
		return l >= right.(int)
	case float64:
		return l >= right.(float64)
	case string:
		return l >= right.(string)
	case []any:
		ra := right.([]any)
		max := len(l)
		if len(ra) > max {
			max = len(ra)
		}
		for i := 0; i < max; i++ {
			var li any
			var rai any
			if len(l) > i {
				li = l[i]
			}
			if len(ra) > i {
				rai = ra[i]
			}

			if !eq(li, rai) {
				return ge(li, rai)
			}
		}

		return true
	case map[string]any:
		ro := right.(map[string]any)

		for _, key := range keys(l, ro) {
			var li any
			var rai any
			if lv, ok := l[key]; ok {
				li = lv
			}
			if rv, ok := ro[key]; ok {
				rai = rv
			}

			if !eq(li, rai) {
				return gt(li, rai)
			}
		}

		return true
	default:
		panic("unknown type")
	}
}

func in(left, right any) bool {
	a, ok := right.([]any)
	if !ok {
		return false
	}
	for _, v := range a {
		if left == v {
			return true
		}
	}

	return false
}

func like(left, right any) (bool, error) {
	return match(right.(string), left.(string))
}

func regexMatch(left, right any) (bool, error) {
	return regexp.Match(right.(string), []byte(left.(string)))
}

func regexNonMatch(left, right any) (bool, error) {
	m, err := regexp.Match(right.(string), []byte(left.(string)))

	return !m, err
}

func typeValue(v any) int {
	switch v.(type) {
	case nil:
		return 0
	case bool:
		return 1
	case float64, int:
		return 2
	case string:
		return 3
	case []any:
		return 4
	case map[string]any:
		return 5
	default:
		panic("unknown type")
	}
}

// Ternary operator https://www.arangodb.com/docs/3.7/aql/operators.html#ternary-operator

func ternary(left, middle, right any) any {
	if toBool(left) {
		if middle != nil {
			return middle
		}

		return left
	}

	return right
}

// Range operators https://www.arangodb.com/docs/3.7/aql/operators.html#range-operator

func aqlrange(left, right any) []float64 {
	var v []float64
	for i := int(left.(float64)); i <= int(right.(float64)); i++ {
		v = append(v, float64(i))
	}

	return v
}
