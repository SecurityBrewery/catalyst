package caql

import (
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Logical operators https://www.arangodb.com/docs/3.7/aql/operators.html#logical-operators

func or(left, right interface{}) interface{} {
	if toBool(left) {
		return left
	}
	return right
}

func and(left, right interface{}) interface{} {
	if !toBool(left) {
		return left
	}
	return right
}

func toBool(i interface{}) bool {
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
	case []interface{}:
		return true
	case map[string]interface{}:
		return true
	default:
		panic("bool conversion failed")
	}
}

// Arithmetic operators https://www.arangodb.com/docs/3.7/aql/operators.html#arithmetic-operators

func plus(left, right interface{}) float64 {
	return toNumber(left) + toNumber(right)
}

func minus(left, right interface{}) float64 {
	return toNumber(left) - toNumber(right)
}

func times(left, right interface{}) float64 {
	return round(toNumber(left) * toNumber(right))
}

func round(r float64) float64 {
	return math.Round(r*100000) / 100000
}

func div(left, right interface{}) float64 {
	b := toNumber(right)
	if b == 0 {
		return 0
	}
	return round(toNumber(left) / b)
}

func mod(left, right interface{}) float64 {
	return math.Mod(toNumber(left), toNumber(right))
}

func toNumber(i interface{}) float64 {
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
	case []interface{}:
		if len(v) == 0 {
			return 0
		}
		if len(v) == 1 {
			return toNumber(v[0])
		}
		return 0
	case map[string]interface{}:
		return 0
	default:
		panic("number conversion error")
	}
}

// Logical operators https://www.arangodb.com/docs/3.7/aql/operators.html#logical-operators
// Order https://www.arangodb.com/docs/3.7/aql/fundamentals-type-value-order.html

func eq(left, right interface{}) bool {
	leftV, rightV := typeValue(left), typeValue(right)
	if leftV != rightV {
		return false
	}
	switch l := left.(type) {
	case nil:
		return true
	case bool, float64, string:
		return left == right
	case []interface{}:
		ra := right.([]interface{})
		max := len(l)
		if len(ra) > max {
			max = len(ra)
		}
		for i := 0; i < max; i++ {
			var li interface{} = nil
			var rai interface{} = nil
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
	case map[string]interface{}:
		ro := right.(map[string]interface{})

		for _, key := range keys(l, ro) {
			var li interface{} = nil
			var rai interface{} = nil
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

func ne(left, right interface{}) bool {
	return !eq(left, right)
}

func lt(left, right interface{}) bool {
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
	case []interface{}:
		ra := right.([]interface{})
		max := len(l)
		if len(ra) > max {
			max = len(ra)
		}
		for i := 0; i < max; i++ {
			var li interface{} = nil
			var rai interface{} = nil
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
	case map[string]interface{}:
		ro := right.(map[string]interface{})

		for _, key := range keys(l, ro) {
			var li interface{} = nil
			var rai interface{} = nil
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

func keys(l map[string]interface{}, ro map[string]interface{}) []string {
	var keys []string
	seen := map[string]bool{}
	for _, a := range []map[string]interface{}{l, ro} {
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

func gt(left, right interface{}) bool {
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
	case []interface{}:
		ra := right.([]interface{})
		max := len(l)
		if len(ra) > max {
			max = len(ra)
		}
		for i := 0; i < max; i++ {
			var li interface{} = nil
			var rai interface{} = nil
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
	case map[string]interface{}:
		ro := right.(map[string]interface{})

		for _, key := range keys(l, ro) {
			var li interface{} = nil
			var rai interface{} = nil
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

func le(left, right interface{}) bool {
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
	case []interface{}:
		ra := right.([]interface{})
		max := len(l)
		if len(ra) > max {
			max = len(ra)
		}
		for i := 0; i < max; i++ {
			var li interface{} = nil
			var rai interface{} = nil
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
	case map[string]interface{}:
		ro := right.(map[string]interface{})

		for _, key := range keys(l, ro) {
			var li interface{} = nil
			var rai interface{} = nil
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

func ge(left, right interface{}) bool {
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
	case []interface{}:
		ra := right.([]interface{})
		max := len(l)
		if len(ra) > max {
			max = len(ra)
		}
		for i := 0; i < max; i++ {
			var li interface{} = nil
			var rai interface{} = nil
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
	case map[string]interface{}:
		ro := right.(map[string]interface{})

		for _, key := range keys(l, ro) {
			var li interface{} = nil
			var rai interface{} = nil
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

func in(left, right interface{}) bool {
	a, ok := right.([]interface{})
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

func like(left, right interface{}) (bool, error) {
	return match(right.(string), left.(string))
}

func regexMatch(left, right interface{}) (bool, error) {
	return regexp.Match(right.(string), []byte(left.(string)))
}

func regexNonMatch(left, right interface{}) (bool, error) {
	m, err := regexp.Match(right.(string), []byte(left.(string)))
	return !m, err
}

func typeValue(v interface{}) int {
	switch v.(type) {
	case nil:
		return 0
	case bool:
		return 1
	case float64, int:
		return 2
	case string:
		return 3
	case []interface{}:
		return 4
	case map[string]interface{}:
		return 5
	default:
		panic("unknown type")
	}
}

// Ternary operator https://www.arangodb.com/docs/3.7/aql/operators.html#ternary-operator

func ternary(left, middle, right interface{}) interface{} {
	if toBool(left) {
		if middle != nil {
			return middle
		}
		return left
	}
	return right
}

// Range operators https://www.arangodb.com/docs/3.7/aql/operators.html#range-operator

func aqlrange(left, right interface{}) []float64 {
	var v []float64
	for i := int(left.(float64)); i <= int(right.(float64)); i++ {
		v = append(v, float64(i))
	}
	return v
}
