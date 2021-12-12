package caql

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/imdario/mergo"

	"github.com/SecurityBrewery/catalyst/generated/caql/parser"
)

func (s *aqlInterpreter) function(ctx *parser.Function_callContext) {
	switch strings.ToUpper(ctx.T_STRING().GetText()) {

	default:
		s.appendErrors(errors.New("unknown function"))

	// Array https://www.arangodb.com/docs/stable/aql/functions-array.html
	case "APPEND":
		u := false
		if len(ctx.AllExpression()) == 3 {
			u = s.pop().(bool)
		}
		seen := map[interface{}]bool{}
		values, anyArray := s.pop().([]interface{}), s.pop().([]interface{})

		if u {
			for _, e := range anyArray {
				seen[e] = true
			}
		}

		for _, e := range values {
			_, ok := seen[e]
			if !ok || !u {
				seen[e] = true
				anyArray = append(anyArray, e)
			}
		}
		s.push(anyArray)
	case "COUNT_DISTINCT", "COUNT_UNIQUE":
		count := 0
		seen := map[interface{}]bool{}
		array := s.pop().([]interface{})
		for _, e := range array {
			_, ok := seen[e]
			if !ok {
				seen[e] = true
				count += 1
			}
		}
		s.push(float64(count))
	case "FIRST":
		array := s.pop().([]interface{})
		if len(array) == 0 {
			s.push(nil)
		} else {
			s.push(array[0])
		}
	// case "FLATTEN":
	// case "INTERLEAVE":
	case "INTERSECTION":
		iset := New(s.pop().([]interface{})...)

		for i := 1; i < len(ctx.AllExpression()); i++ {
			iset = iset.Intersection(New(s.pop().([]interface{})...))
		}

		s.push(iset.Values())
	// case "JACCARD":
	case "LAST":
		array := s.pop().([]interface{})
		if len(array) == 0 {
			s.push(nil)
		} else {
			s.push(array[len(array)-1])
		}
	case "COUNT", "LENGTH":
		switch v := s.pop().(type) {
		case nil:
			s.push(float64(0))
		case bool:
			if v {
				s.push(float64(1))
			} else {
				s.push(float64(0))
			}
		case float64:
			s.push(float64(len(fmt.Sprint(v))))
		case string:
			s.push(float64(utf8.RuneCountInString(v)))
		case []interface{}:
			s.push(float64(len(v)))
		case map[string]interface{}:
			s.push(float64(len(v)))
		default:
			panic("unknown type")
		}
	case "MINUS":
		var sets []*Set
		for i := 0; i < len(ctx.AllExpression()); i++ {
			sets = append(sets, New(s.pop().([]interface{})...))
		}

		iset := sets[len(sets)-1]
		// for i := len(sets)-1; i > 0; i-- {
		for i := 0; i < len(sets)-1; i++ {
			iset = iset.Minus(sets[i])
		}

		s.push(iset.Values())
	case "NTH":
		pos := s.pop().(float64)
		array := s.pop().([]interface{})
		if int(pos) >= len(array) || pos < 0 {
			s.push(nil)
		} else {
			s.push(array[int64(pos)])
		}
	// case "OUTERSECTION":
	// 	array := s.pop().([]interface{})
	// 	union := New(array...)
	// 	intersection := New(s.pop().([]interface{})...)
	// 	for i := 1; i < len(ctx.AllExpression()); i++ {
	// 		array = s.pop().([]interface{})
	// 		union = union.Union(New(array...))
	// 		intersection = intersection.Intersection(New(array...))
	// 	}
	// 	s.push(union.Minus(intersection).Values())
	case "POP":
		array := s.pop().([]interface{})
		s.push(array[:len(array)-1])
	case "POSITION", "CONTAINS_ARRAY":
		returnIndex := false
		if len(ctx.AllExpression()) == 3 {
			returnIndex = s.pop().(bool)
		}
		search := s.pop()
		array := s.pop().([]interface{})

		for idx, e := range array {
			if e == search {
				if returnIndex {
					s.push(float64(idx))
				} else {
					s.push(true)
				}
			}
		}

		if returnIndex {
			s.push(float64(-1))
		} else {
			s.push(false)
		}
	case "PUSH":
		u := false
		if len(ctx.AllExpression()) == 3 {
			u = s.pop().(bool)
		}
		element := s.pop()
		array := s.pop().([]interface{})

		if u && contains(array, element) {
			s.push(array)
		} else {
			s.push(append(array, element))
		}
	case "REMOVE_NTH":
		position := s.pop().(float64)
		anyArray := s.pop().([]interface{})

		if position < 0 {
			position = float64(len(anyArray) + int(position))
		}

		result := []interface{}{}
		for idx, e := range anyArray {
			if idx != int(position) {
				result = append(result, e)
			}
		}
		s.push(result)
	case "REPLACE_NTH":
		defaultPaddingValue := ""
		if len(ctx.AllExpression()) == 4 {
			defaultPaddingValue = s.pop().(string)
		}
		replaceValue := s.pop().(string)
		position := s.pop().(float64)
		anyArray := s.pop().([]interface{})

		if position < 0 {
			position = float64(len(anyArray) + int(position))
			if position < 0 {
				position = 0
			}
		}

		switch {
		case int(position) < len(anyArray):
			anyArray[int(position)] = replaceValue
		case int(position) == len(anyArray):
			anyArray = append(anyArray, replaceValue)
		default:
			if defaultPaddingValue == "" {
				panic("missing defaultPaddingValue")
			}
			for len(anyArray) < int(position) {
				anyArray = append(anyArray, defaultPaddingValue)
			}
			anyArray = append(anyArray, replaceValue)
		}

		s.push(anyArray)
	case "REMOVE_VALUE":
		limit := math.Inf(1)
		if len(ctx.AllExpression()) == 3 {
			limit = s.pop().(float64)
		}
		value := s.pop()
		array := s.pop().([]interface{})
		result := []interface{}{}
		for idx, e := range array {
			if e != value || float64(idx) > limit {
				result = append(result, e)
			}
		}
		s.push(result)
	case "REMOVE_VALUES":
		values := s.pop().([]interface{})
		array := s.pop().([]interface{})
		result := []interface{}{}
		for _, e := range array {
			if !contains(values, e) {
				result = append(result, e)
			}
		}
		s.push(result)
	case "REVERSE":
		array := s.pop().([]interface{})
		var reverse []interface{}
		for _, e := range array {
			reverse = append([]interface{}{e}, reverse...)
		}
		s.push(reverse)
	case "SHIFT":
		s.push(s.pop().([]interface{})[1:])
	case "SLICE":
		length := float64(-1)
		full := true
		if len(ctx.AllExpression()) == 3 {
			length = s.pop().(float64)
			full = false
		}
		start := int64(s.pop().(float64))
		array := s.pop().([]interface{})

		if start < 0 {
			start = int64(len(array)) + start
		}
		if full {
			length = float64(int64(len(array)) - start)
		}

		end := int64(0)
		if length < 0 {
			end = int64(len(array)) + int64(length)
		} else {
			end = start + int64(length)
		}
		s.push(array[start:end])
	case "SORTED":
		array := s.pop().([]interface{})
		sort.Slice(array, func(i, j int) bool { return lt(array[i], array[j]) })
		s.push(array)
	case "SORTED_UNIQUE":
		array := s.pop().([]interface{})
		sort.Slice(array, func(i, j int) bool { return lt(array[i], array[j]) })
		s.push(unique(array))
	case "UNION":
		array := s.pop().([]interface{})

		for i := 1; i < len(ctx.AllExpression()); i++ {
			array = append(array, s.pop().([]interface{})...)
		}

		sort.Slice(array, func(i, j int) bool { return lt(array[i], array[j]) })
		s.push(array)
	case "UNION_DISTINCT":
		iset := New(s.pop().([]interface{})...)

		for i := 1; i < len(ctx.AllExpression()); i++ {
			iset = iset.Union(New(s.pop().([]interface{})...))
		}

		s.push(unique(iset.Values()))
	case "UNIQUE":
		s.push(unique(s.pop().([]interface{})))
	case "UNSHIFT":
		u := false
		if len(ctx.AllExpression()) == 3 {
			u = s.pop().(bool)
		}
		element := s.pop()
		array := s.pop().([]interface{})
		if u && contains(array, element) {
			s.push(array)
		} else {
			s.push(append([]interface{}{element}, array...))
		}

	// Bit https://www.arangodb.com/docs/stable/aql/functions-bit.html
	// case "BIT_AND":
	// case "BIT_CONSTRUCT":
	// case "BIT_DECONSTRUCT":
	// case "BIT_FROM_STRING":
	// case "BIT_NEGATE":
	// case "BIT_OR":
	// case "BIT_POPCOUNT":
	// case "BIT_SHIFT_LEFT":
	// case "BIT_SHIFT_RIGHT":
	// case "BIT_TEST":
	// case "BIT_TO_STRING":
	// case "BIT_XOR":

	// Date https://www.arangodb.com/docs/stable/aql/functions-date.html
	// case "DATE_NOW":
	// case "DATE_ISO8601":
	// case "DATE_TIMESTAMP":
	// case "IS_DATESTRING":

	// case "DATE_DAYOFWEEK":
	// case "DATE_YEAR":
	// case "DATE_MONTH":
	// case "DATE_DAY":
	// case "DATE_HOUR":
	// case "DATE_MINUTE":
	// case "DATE_SECOND":
	// case "DATE_MILLISECOND":

	// case "DATE_DAYOFYEAR":
	// case "DATE_ISOWEEK":
	// case "DATE_LEAPYEAR":
	// case "DATE_QUARTER":
	// case "DATE_DAYS_IN_MONTH":
	// case "DATE_TRUNC":
	// case "DATE_ROUND":
	// case "DATE_FORMAT":

	// case "DATE_ADD":
	// case "DATE_SUBTRACT":
	// case "DATE_DIFF":
	// case "DATE_COMPARE":

	// Document https://www.arangodb.com/docs/stable/aql/functions-document.html
	case "ATTRIBUTES":
		if len(ctx.AllExpression()) == 3 {
			s.pop() // always sort
		}
		removeInternal := false
		if len(ctx.AllExpression()) >= 2 {
			removeInternal = s.pop().(bool)
		}
		var keys []interface{}
		for k := range s.pop().(map[string]interface{}) {
			isInternalKey := strings.HasPrefix(k, "_")
			if !removeInternal || !isInternalKey {
				keys = append(keys, k)
			}
		}
		sort.Slice(keys, func(i, j int) bool { return lt(keys[i], keys[j]) })
		s.push(keys)
	// case "COUNT":
	case "HAS":
		right, left := s.pop(), s.pop()
		_, ok := left.(map[string]interface{})[right.(string)]
		s.push(ok)
	// case "KEEP":
	// case "LENGTH":
	// case "MATCHES":
	case "MERGE":
		var docs []map[string]interface{}
		if len(ctx.AllExpression()) == 1 {
			for _, doc := range s.pop().([]interface{}) {
				docs = append([]map[string]interface{}{doc.(map[string]interface{})}, docs...)
			}
		} else {
			for i := 0; i < len(ctx.AllExpression()); i++ {
				docs = append(docs, s.pop().(map[string]interface{}))
			}
		}

		doc := docs[len(docs)-1]
		for i := len(docs) - 2; i >= 0; i-- {
			for k, v := range docs[i] {
				doc[k] = v
			}
		}
		s.push(doc)
	case "MERGE_RECURSIVE":
		var doc map[string]interface{}
		for i := 0; i < len(ctx.AllExpression()); i++ {
			err := mergo.Merge(&doc, s.pop().(map[string]interface{}))
			if err != nil {
				panic(err)
			}
		}
		s.push(doc)
	// case "PARSE_IDENTIFIER":
	// case "TRANSLATE":
	// case "UNSET":
	// case "UNSET_RECURSIVE":
	case "VALUES":
		removeInternal := false
		if len(ctx.AllExpression()) == 2 {
			removeInternal = s.pop().(bool)
		}
		var values []interface{}
		for k, v := range s.pop().(map[string]interface{}) {
			isInternalKey := strings.HasPrefix(k, "_")
			if !removeInternal || !isInternalKey {
				values = append(values, v)
			}
		}
		sort.Slice(values, func(i, j int) bool { return lt(values[i], values[j]) })
		s.push(values)
	// case "ZIP":

	// Numeric https://www.arangodb.com/docs/stable/aql/functions-numeric.html
	case "ABS":
		s.push(math.Abs(s.pop().(float64)))
	case "ACOS":
		v := s.pop().(float64)
		asin := math.Acos(v)
		if v > 1 || v < -1 {
			s.push(nil)
		} else {
			s.push(asin)
		}
	case "ASIN":
		v := s.pop().(float64)
		asin := math.Asin(v)
		if v > 1 || v < -1 {
			s.push(nil)
		} else {
			s.push(asin)
		}
	case "ATAN":
		s.push(math.Atan(s.pop().(float64)))
	case "ATAN2":
		s.push(math.Atan2(s.pop().(float64), s.pop().(float64)))
	case "AVERAGE", "AVG":
		count := 0
		sum := float64(0)
		array := s.pop().([]interface{})
		for _, element := range array {
			if element != nil {
				count += 1
				sum += toNumber(element)
			}
		}
		if count == 0 {
			s.push(nil)
		} else {
			s.push(sum / float64(count))
		}
	case "CEIL":
		s.push(math.Ceil(s.pop().(float64)))
	case "COS":
		s.push(math.Cos(s.pop().(float64)))
	case "DEGREES":
		s.push(s.pop().(float64) * 180 / math.Pi)
	case "EXP":
		s.push(math.Exp(s.pop().(float64)))
	case "EXP2":
		s.push(math.Exp2(s.pop().(float64)))
	case "FLOOR":
		s.push(math.Floor(s.pop().(float64)))
	case "LOG":
		l := math.Log(s.pop().(float64))
		if l <= 0 {
			s.push(nil)
		} else {
			s.push(l)
		}
	case "LOG2":
		l := math.Log2(s.pop().(float64))
		if l <= 0 {
			s.push(nil)
		} else {
			s.push(l)
		}
	case "LOG10":
		l := math.Log10(s.pop().(float64))
		if l <= 0 {
			s.push(nil)
		} else {
			s.push(l)
		}
	case "MAX":
		var set bool
		var max float64
		array := s.pop().([]interface{})
		for _, element := range array {
			if element != nil {
				if !set || toNumber(element) > max {
					max = toNumber(element)
					set = true
				}
			}
		}
		if set {
			s.push(max)
		} else {
			s.push(nil)
		}
	case "MEDIAN":
		array := s.pop().([]interface{})
		var numbers []float64
		for _, element := range array {
			if f, ok := element.(float64); ok {
				numbers = append(numbers, f)
			}
		}

		sort.Float64s(numbers) // sort the numbers

		middlePos := len(numbers) / 2

		switch {
		case len(numbers) == 0:
			s.push(nil)
		case len(numbers)%2 == 1:
			s.push(numbers[middlePos])
		default:
			s.push((numbers[middlePos-1] + numbers[middlePos]) / 2)
		}
	case "MIN":
		var set bool
		var min float64
		array := s.pop().([]interface{})
		for _, element := range array {
			if element != nil {
				if !set || toNumber(element) < min {
					min = toNumber(element)
					set = true
				}
			}
		}
		if set {
			s.push(min)
		} else {
			s.push(nil)
		}
	// case "PERCENTILE":
	case "PI":
		s.push(math.Pi)
	case "POW":
		right, left := s.pop(), s.pop()
		s.push(math.Pow(left.(float64), right.(float64)))
	case "PRODUCT":
		product := float64(1)
		array := s.pop().([]interface{})
		for _, element := range array {
			if element != nil {
				product *= toNumber(element)
			}
		}
		s.push(product)
	case "RADIANS":
		s.push(s.pop().(float64) * math.Pi / 180)
	case "RAND":
		s.push(rand.Float64())
	case "RANGE":
		var array []interface{}
		var start, end, step float64
		if len(ctx.AllExpression()) == 2 {
			right, left := s.pop(), s.pop()
			start = math.Trunc(left.(float64))
			end = math.Trunc(right.(float64))
			step = 1
		} else {
			middle, right, left := s.pop(), s.pop(), s.pop()
			start = left.(float64)
			end = right.(float64)
			step = middle.(float64)
		}
		for i := start; i <= end; i += step {
			array = append(array, i)
		}
		s.push(array)
	case "ROUND":
		x := s.pop().(float64)
		t := math.Trunc(x)
		if math.Abs(x-t) == 0.5 {
			s.push(x + 0.5)
		} else {
			s.push(math.Round(x))
		}
	case "SIN":
		s.push(math.Sin(s.pop().(float64)))
	case "SQRT":
		s.push(math.Sqrt(s.pop().(float64)))
	// case "STDDEV_POPULATION":
	// case "STDDEV_SAMPLE":
	// case "STDDEV":
	case "SUM":
		sum := float64(0)
		array := s.pop().([]interface{})
		for _, element := range array {
			sum += toNumber(element)
		}
		s.push(sum)
	case "TAN":
		s.push(math.Tan(s.pop().(float64)))
	// case "VARIANCE_POPULATION", "VARIANCE":
	// case "VARIANCE_SAMPLE":

	// String https://www.arangodb.com/docs/stable/aql/functions-string.html
	// case "CHAR_LENGTH":
	// case "CONCAT":
	// case "CONCAT_SEPARATOR":
	// case "CONTAINS":
	// case "CRC32":
	// case "ENCODE_URI_COMPONENT":
	// case "FIND_FIRST":
	// case "FIND_LAST":
	// case "FNV64":
	// case "IPV4_FROM_NUMBER":
	// case "IPV4_TO_NUMBER":
	// case "IS_IPV4":
	// case "JSON_PARSE":
	// case "JSON_STRINGIFY":
	// case "LEFT":
	// case "LENGTH":
	// case "LEVENSHTEIN_DISTANCE":
	// case "LIKE":
	case "LOWER":
		s.push(strings.ToLower(s.pop().(string)))
	// case "LTRIM":
	// case "MD5":
	// case "NGRAM_POSITIONAL_SIMILARITY":
	// case "NGRAM_SIMILARITY":
	// case "RANDOM_TOKEN":
	// case "REGEX_MATCHES":
	// case "REGEX_SPLIT":
	// case "REGEX_TEST":
	// case "REGEX_REPLACE":
	// case "REVERSE":
	// case "RIGHT":
	// case "RTRIM":
	// case "SHA1":
	// case "SHA512":
	// case "SOUNDEX":
	// case "SPLIT":
	// case "STARTS_WITH":
	// case "SUBSTITUTE":
	// case "SUBSTRING":
	// case "TOKENS":
	// case "TO_BASE64":
	// case "TO_HEX":
	// case "TRIM":
	case "UPPER":
		s.push(strings.ToUpper(s.pop().(string)))
	// case "UUID":

	// Type cast https://www.arangodb.com/docs/stable/aql/functions-type-cast.html
	case "TO_BOOL":
		s.push(toBool(s.pop()))
	case "TO_NUMBER":
		s.push(toNumber(s.pop()))
		// case "TO_STRING":
		// case "TO_ARRAY":
		// case "TO_LIST":

		// case "IS_NULL":
		// case "IS_BOOL":
		// case "IS_NUMBER":
		// case "IS_STRING":
		// case "IS_ARRAY":
		// case "IS_LIST":
		// case "IS_OBJECT":
		// case "IS_DOCUMENT":
		// case "IS_DATESTRING":
		// case "IS_IPV4":
		// case "IS_KEY":
		// case "TYPENAME":

	}
}

func unique(array []interface{}) []interface{} {
	seen := map[interface{}]bool{}
	var filtered []interface{}
	for _, e := range array {
		_, ok := seen[e]
		if !ok {
			seen[e] = true
			filtered = append(filtered, e)
		}
	}
	return filtered
}

func contains(values []interface{}, e interface{}) bool {
	for _, v := range values {
		if e == v {
			return true
		}
	}
	return false
}

func stringSliceContains(values []string, e string) bool {
	for _, v := range values {
		if e == v {
			return true
		}
	}
	return false
}

var functionNames = []string{
	"APPEND", "COUNT_DISTINCT", "COUNT_UNIQUE", "FIRST", "FLATTEN", "INTERLEAVE", "INTERSECTION", "JACCARD", "LAST",
	"COUNT", "LENGTH", "MINUS", "NTH", "OUTERSECTION", "POP", "POSITION", "CONTAINS_ARRAY", "PUSH", "REMOVE_NTH",
	"REPLACE_NTH", "REMOVE_VALUE", "REMOVE_VALUES", "REVERSE", "SHIFT", "SLICE", "SORTED", "SORTED_UNIQUE", "UNION",
	"UNION_DISTINCT", "UNIQUE", "UNSHIFT", "BIT_AND", "BIT_CONSTRUCT", "BIT_DECONSTRUCT", "BIT_FROM_STRING",
	"BIT_NEGATE", "BIT_OR", "BIT_POPCOUNT", "BIT_SHIFT_LEFT", "BIT_SHIFT_RIGHT", "BIT_TEST", "BIT_TO_STRING",
	"BIT_XOR", "DATE_NOW", "DATE_ISO8601", "DATE_TIMESTAMP", "IS_DATESTRING", "DATE_DAYOFWEEK", "DATE_YEAR",
	"DATE_MONTH", "DATE_DAY", "DATE_HOUR", "DATE_MINUTE", "DATE_SECOND", "DATE_MILLISECOND", "DATE_DAYOFYEAR",
	"DATE_ISOWEEK", "DATE_LEAPYEAR", "DATE_QUARTER", "DATE_DAYS_IN_MONTH", "DATE_TRUNC", "DATE_ROUND", "DATE_FORMAT",
	"DATE_ADD", "DATE_SUBTRACT", "DATE_DIFF", "DATE_COMPARE", "ATTRIBUTES", "COUNT", "HAS", "KEEP", "LENGTH",
	"MATCHES", "MERGE", "MERGE_RECURSIVE", "PARSE_IDENTIFIER", "TRANSLATE", "UNSET", "UNSET_RECURSIVE", "VALUES",
	"ZIP", "ABS", "ACOS", "ASIN", "ATAN", "ATAN2", "AVERAGE", "AVG", "CEIL", "COS", "DEGREES", "EXP", "EXP2", "FLOOR",
	"LOG", "LOG2", "LOG10", "MAX", "MEDIAN", "MIN", "PERCENTILE", "PI", "POW", "PRODUCT", "RADIANS", "RAND", "RANGE",
	"ROUND", "SIN", "SQRT", "STDDEV_POPULATION", "STDDEV_SAMPLE", "STDDEV", "SUM", "TAN", "VARIANCE_POPULATION",
	"VARIANCE", "VARIANCE_SAMPLE", "CHAR_LENGTH", "CONCAT", "CONCAT_SEPARATOR", "CONTAINS", "CRC32",
	"ENCODE_URI_COMPONENT", "FIND_FIRST", "FIND_LAST", "FNV64", "IPV4_FROM_NUMBER", "IPV4_TO_NUMBER", "IS_IPV4",
	"JSON_PARSE", "JSON_STRINGIFY", "LEFT", "LENGTH", "LEVENSHTEIN_DISTANCE", "LIKE", "LOWER", "LTRIM", "MD5",
	"NGRAM_POSITIONAL_SIMILARITY", "NGRAM_SIMILARITY", "RANDOM_TOKEN", "REGEX_MATCHES", "REGEX_SPLIT", "REGEX_TEST",
	"REGEX_REPLACE", "REVERSE", "RIGHT", "RTRIM", "SHA1", "SHA512", "SOUNDEX", "SPLIT", "STARTS_WITH", "SUBSTITUTE",
	"SUBSTRING", "TOKENS", "TO_BASE64", "TO_HEX", "TRIM", "UPPER", "UUID", "TO_BOOL", "TO_NUMBER", "TO_STRING",
	"TO_ARRAY", "TO_LIST", "IS_NULL", "IS_BOOL", "IS_NUMBER", "IS_STRING", "IS_ARRAY", "IS_LIST", "IS_OBJECT",
	"IS_DOCUMENT", "IS_DATESTRING", "IS_IPV4", "IS_KEY", "TYPENAME"}
