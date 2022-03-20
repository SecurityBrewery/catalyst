package caql_test

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"

	"github.com/SecurityBrewery/catalyst/caql"
)

func TestFunctions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		saql           string
		wantRebuild    string
		wantValue      any
		wantParseErr   bool
		wantRebuildErr bool
		wantEvalErr    bool
		values         string
	}{
		// https://www.arangodb.com/docs/3.7/aql/functions-array.html
		{name: "APPEND", saql: `APPEND([1, 2, 3], [5, 6, 9])`, wantRebuild: `APPEND([1, 2, 3], [5, 6, 9])`, wantValue: jsonParse(`[1, 2, 3, 5, 6, 9]`)},
		{name: "APPEND", saql: `APPEND([1, 2, 3], [3, 4, 5, 2, 9], true)`, wantRebuild: `APPEND([1, 2, 3], [3, 4, 5, 2, 9], true)`, wantValue: jsonParse(`[1, 2, 3, 4, 5, 9]`)},
		{name: "COUNT_DISTINCT", saql: `COUNT_DISTINCT([1, 2, 3])`, wantRebuild: `COUNT_DISTINCT([1, 2, 3])`, wantValue: 3},
		{name: "COUNT_DISTINCT", saql: `COUNT_DISTINCT(["yes", "no", "yes", "sauron", "no", "yes"])`, wantRebuild: `COUNT_DISTINCT(["yes", "no", "yes", "sauron", "no", "yes"])`, wantValue: 3},
		{name: "FIRST", saql: `FIRST([1, 2, 3])`, wantRebuild: `FIRST([1, 2, 3])`, wantValue: 1},
		{name: "FIRST", saql: `FIRST([])`, wantRebuild: `FIRST([])`, wantValue: nil},
		// {name: "FLATTEN", saql: `FLATTEN([1, 2, [3, 4], 5, [6, 7], [8, [9, 10]]])`, wantRebuild: `FLATTEN([1, 2, [3, 4], 5, [6, 7], [8, [9, 10]]])`, wantValue:},
		// {name: "FLATTEN", saql: `FLATTEN([1, 2, [3, 4], 5, [6, 7], [8, [9, 10]]], 2)`, wantRebuild: `FLATTEN([1, 2, [3, 4], 5, [6, 7], [8, [9, 10]]], 2)`, wantValue:},
		// {name: "INTERLEAVE", saql: `INTERLEAVE([1, 1, 1], [2, 2, 2], [3, 3, 3])`, wantRebuild: `INTERLEAVE([1, 1, 1], [2, 2, 2], [3, 3, 3])`, wantValue:},
		// {name: "INTERLEAVE", saql: `INTERLEAVE([1], [2, 2], [3, 3, 3])`, wantRebuild: `INTERLEAVE([1], [2, 2], [3, 3, 3])`, wantValue:},
		{name: "INTERSECTION", saql: `INTERSECTION([1,2,3,4,5], [2,3,4,5,6], [3,4,5,6,7])`, wantRebuild: `INTERSECTION([1, 2, 3, 4, 5], [2, 3, 4, 5, 6], [3, 4, 5, 6, 7])`, wantValue: jsonParse(`[3, 4, 5]`)},
		{name: "INTERSECTION", saql: `INTERSECTION([2,4,6], [8,10,12], [14,16,18])`, wantRebuild: `INTERSECTION([2, 4, 6], [8, 10, 12], [14, 16, 18])`, wantValue: jsonParse(`[]`)},
		// {name: "JACCARD", saql: `JACCARD([1,2,3,4], [3,4,5,6])`, wantRebuild: `JACCARD([1,2,3,4], [3,4,5,6])`, wantValue: 0.3333333333333333},
		// {name: "JACCARD", saql: `JACCARD([1,1,2,2,2,3], [2,2,3,4])`, wantRebuild: `JACCARD([1,1,2,2,2,3], [2,2,3,4])`, wantValue: 0.5},
		// {name: "JACCARD", saql: `JACCARD([1,2,3], [])`, wantRebuild: `JACCARD([1, 2, 3], [])`, wantValue: 0},
		// {name: "JACCARD", saql: `JACCARD([], [])`, wantRebuild: `JACCARD([], [])`, wantValue: 1},
		{name: "LAST", saql: `LAST([1,2,3,4,5])`, wantRebuild: `LAST([1, 2, 3, 4, 5])`, wantValue: 5},
		{name: "LENGTH", saql: `LENGTH("🥑")`, wantRebuild: `LENGTH("🥑")`, wantValue: 1},
		{name: "LENGTH", saql: `LENGTH(1234)`, wantRebuild: `LENGTH(1234)`, wantValue: 4},
		{name: "LENGTH", saql: `LENGTH([1,2,3,4,5,6,7])`, wantRebuild: `LENGTH([1, 2, 3, 4, 5, 6, 7])`, wantValue: 7},
		{name: "LENGTH", saql: `LENGTH(false)`, wantRebuild: `LENGTH(false)`, wantValue: 0},
		{name: "LENGTH", saql: `LENGTH({a:1, b:2, c:3, d:4, e:{f:5,g:6}})`, wantRebuild: `LENGTH({a: 1, b: 2, c: 3, d: 4, e: {f: 5, g: 6}})`, wantValue: 5},
		{name: "MINUS", saql: `MINUS([1,2,3,4], [3,4,5,6], [5,6,7,8])`, wantRebuild: `MINUS([1, 2, 3, 4], [3, 4, 5, 6], [5, 6, 7, 8])`, wantValue: jsonParse(`[1, 2]`)},
		{name: "NTH", saql: `NTH(["foo", "bar", "baz"], 2)`, wantRebuild: `NTH(["foo", "bar", "baz"], 2)`, wantValue: "baz"},
		{name: "NTH", saql: `NTH(["foo", "bar", "baz"], 3)`, wantRebuild: `NTH(["foo", "bar", "baz"], 3)`, wantValue: nil},
		{name: "NTH", saql: `NTH(["foo", "bar", "baz"], -1)`, wantRebuild: `NTH(["foo", "bar", "baz"], -1)`, wantValue: nil},
		// {name: "OUTERSECTION", saql: `OUTERSECTION([1, 2, 3], [2, 3, 4], [3, 4, 5])`, wantRebuild: `OUTERSECTION([1, 2, 3], [2, 3, 4], [3, 4, 5])`, wantValue: jsonParse(`[1, 5]`)},
		{name: "POP", saql: `POP([1, 2, 3, 4])`, wantRebuild: `POP([1, 2, 3, 4])`, wantValue: jsonParse(`[1, 2, 3]`)},
		{name: "POP", saql: `POP([1])`, wantRebuild: `POP([1])`, wantValue: jsonParse(`[]`)},
		{name: "POSITION", saql: `POSITION([2,4,6,8], 4)`, wantRebuild: `POSITION([2, 4, 6, 8], 4)`, wantValue: true},
		{name: "POSITION", saql: `POSITION([2,4,6,8], 4, true)`, wantRebuild: `POSITION([2, 4, 6, 8], 4, true)`, wantValue: 1},
		{name: "PUSH", saql: `PUSH([1, 2, 3], 4)`, wantRebuild: `PUSH([1, 2, 3], 4)`, wantValue: jsonParse(`[1, 2, 3, 4]`)},
		{name: "PUSH", saql: `PUSH([1, 2, 2, 3], 2, true)`, wantRebuild: `PUSH([1, 2, 2, 3], 2, true)`, wantValue: jsonParse(`[1, 2, 2, 3]`)},
		{name: "REMOVE_NTH", saql: `REMOVE_NTH(["a", "b", "c", "d", "e"], 1)`, wantRebuild: `REMOVE_NTH(["a", "b", "c", "d", "e"], 1)`, wantValue: jsonParse(`["a", "c", "d", "e"]`)},
		{name: "REMOVE_NTH", saql: `REMOVE_NTH(["a", "b", "c", "d", "e"], -2)`, wantRebuild: `REMOVE_NTH(["a", "b", "c", "d", "e"], -2)`, wantValue: jsonParse(`["a", "b", "c", "e"]`)},
		{name: "REPLACE_NTH", saql: `REPLACE_NTH(["a", "b", "c"], 1 , "z")`, wantRebuild: `REPLACE_NTH(["a", "b", "c"], 1, "z")`, wantValue: jsonParse(`["a", "z", "c"]`)},
		{name: "REPLACE_NTH", saql: `REPLACE_NTH(["a", "b", "c"], 3 , "z")`, wantRebuild: `REPLACE_NTH(["a", "b", "c"], 3, "z")`, wantValue: jsonParse(`["a", "b", "c", "z"]`)},
		{name: "REPLACE_NTH", saql: `REPLACE_NTH(["a", "b", "c"], 6, "z", "y")`, wantRebuild: `REPLACE_NTH(["a", "b", "c"], 6, "z", "y")`, wantValue: jsonParse(`["a", "b", "c", "y", "y", "y", "z"]`)},
		{name: "REPLACE_NTH", saql: `REPLACE_NTH(["a", "b", "c"], -1, "z")`, wantRebuild: `REPLACE_NTH(["a", "b", "c"], -1, "z")`, wantValue: jsonParse(`["a", "b", "z"]`)},
		{name: "REPLACE_NTH", saql: `REPLACE_NTH(["a", "b", "c"], -9, "z")`, wantRebuild: `REPLACE_NTH(["a", "b", "c"], -9, "z")`, wantValue: jsonParse(`["z", "b", "c"]`)},
		{name: "REMOVE_VALUE", saql: `REMOVE_VALUE(["a", "b", "b", "a", "c"], "a")`, wantRebuild: `REMOVE_VALUE(["a", "b", "b", "a", "c"], "a")`, wantValue: jsonParse(`["b", "b", "c"]`)},
		{name: "REMOVE_VALUE", saql: `REMOVE_VALUE(["a", "b", "b", "a", "c"], "a", 1)`, wantRebuild: `REMOVE_VALUE(["a", "b", "b", "a", "c"], "a", 1)`, wantValue: jsonParse(`["b", "b", "a", "c"]`)},
		{name: "REMOVE_VALUES", saql: `REMOVE_VALUES(["a", "a", "b", "c", "d", "e", "f"], ["a", "f", "d"])`, wantRebuild: `REMOVE_VALUES(["a", "a", "b", "c", "d", "e", "f"], ["a", "f", "d"])`, wantValue: jsonParse(`["b", "c", "e"]`)},
		{name: "REVERSE", saql: `REVERSE ([2,4,6,8,10])`, wantRebuild: `REVERSE([2, 4, 6, 8, 10])`, wantValue: jsonParse(`[10, 8, 6, 4, 2]`)},
		{name: "SHIFT", saql: `SHIFT([1, 2, 3, 4])`, wantRebuild: `SHIFT([1, 2, 3, 4])`, wantValue: jsonParse(`[2, 3, 4]`)},
		{name: "SHIFT", saql: `SHIFT([1])`, wantRebuild: `SHIFT([1])`, wantValue: jsonParse(`[]`)},
		{name: "SLICE", saql: `SLICE([1, 2, 3, 4, 5], 0, 1)`, wantRebuild: `SLICE([1, 2, 3, 4, 5], 0, 1)`, wantValue: jsonParse(`[1]`)},
		{name: "SLICE", saql: `SLICE([1, 2, 3, 4, 5], 1, 2)`, wantRebuild: `SLICE([1, 2, 3, 4, 5], 1, 2)`, wantValue: jsonParse(`[2, 3]`)},
		{name: "SLICE", saql: `SLICE([1, 2, 3, 4, 5], 3)`, wantRebuild: `SLICE([1, 2, 3, 4, 5], 3)`, wantValue: jsonParse(`[4, 5]`)},
		{name: "SLICE", saql: `SLICE([1, 2, 3, 4, 5], 1, -1)`, wantRebuild: `SLICE([1, 2, 3, 4, 5], 1, -1)`, wantValue: jsonParse(`[2, 3, 4]`)},
		{name: "SLICE", saql: `SLICE([1, 2, 3, 4, 5], 0, -2)`, wantRebuild: `SLICE([1, 2, 3, 4, 5], 0, -2)`, wantValue: jsonParse(`[1, 2, 3]`)},
		{name: "SLICE", saql: `SLICE([1, 2, 3, 4, 5], -3, 2)`, wantRebuild: `SLICE([1, 2, 3, 4, 5], -3, 2)`, wantValue: jsonParse(`[3, 4]`)},
		{name: "SORTED", saql: `SORTED([8,4,2,10,6])`, wantRebuild: `SORTED([8, 4, 2, 10, 6])`, wantValue: jsonParse(`[2, 4, 6, 8, 10]`)},
		{name: "SORTED_UNIQUE", saql: `SORTED_UNIQUE([8,4,2,10,6,2,8,6,4])`, wantRebuild: `SORTED_UNIQUE([8, 4, 2, 10, 6, 2, 8, 6, 4])`, wantValue: jsonParse(`[2, 4, 6, 8, 10]`)},
		{name: "UNION", saql: `UNION([1, 2, 3], [1, 2])`, wantRebuild: `UNION([1, 2, 3], [1, 2])`, wantValue: jsonParse(`[1, 1, 2, 2, 3]`)},
		{name: "UNION_DISTINCT", saql: `UNION_DISTINCT([1, 2, 3], [1, 2])`, wantRebuild: `UNION_DISTINCT([1, 2, 3], [1, 2])`, wantValue: jsonParse(`[1, 2, 3]`)},
		{name: "UNIQUE", saql: `UNIQUE([1,2,2,3,3,3,4,4,4,4,5,5,5,5,5])`, wantRebuild: `UNIQUE([1, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 5])`, wantValue: jsonParse(`[1, 2, 3, 4, 5]`)},
		{name: "UNSHIFT", saql: `UNSHIFT([1, 2, 3], 4)`, wantRebuild: `UNSHIFT([1, 2, 3], 4)`, wantValue: jsonParse(`[4, 1, 2, 3]`)},
		{name: "UNSHIFT", saql: `UNSHIFT([1, 2, 3], 2, true)`, wantRebuild: `UNSHIFT([1, 2, 3], 2, true)`, wantValue: jsonParse(`[1, 2, 3]`)},

		// https://www.arangodb.com/docs/3.7/aql/functions-bit.html
		// {name: "BIT_CONSTRUCT", saql: `BIT_CONSTRUCT([1, 2, 3])`, wantRebuild: `BIT_CONSTRUCT([1, 2, 3])`, wantValue: 14},
		// {name: "BIT_CONSTRUCT", saql: `BIT_CONSTRUCT([0, 4, 8])`, wantRebuild: `BIT_CONSTRUCT([0, 4, 8])`, wantValue: 273},
		// {name: "BIT_CONSTRUCT", saql: `BIT_CONSTRUCT([0, 1, 10, 31])`, wantRebuild: `BIT_CONSTRUCT([0, 1, 10, 31])`, wantValue: 2147484675},
		// {name: "BIT_DECONSTRUCT", saql: `BIT_DECONSTRUCT(14)`, wantRebuild: `BIT_DECONSTRUCT(14) `, wantValue: []interface{}{1, 2, 3}},
		// {name: "BIT_DECONSTRUCT", saql: `BIT_DECONSTRUCT(273)`, wantRebuild: `BIT_DECONSTRUCT(273)`, wantValue: []interface{}{0, 4, 8}},
		// {name: "BIT_DECONSTRUCT", saql: `BIT_DECONSTRUCT(2147484675)`, wantRebuild: `BIT_DECONSTRUCT(2147484675)`, wantValue: []interface{}{0, 1, 10, 31}},
		// {name: "BIT_FROM_STRING", saql: `BIT_FROM_STRING("0111")`, wantRebuild: `BIT_FROM_STRING("0111")`, wantValue: 7},
		// {name: "BIT_FROM_STRING", saql: `BIT_FROM_STRING("000000000000010")`, wantRebuild: `BIT_FROM_STRING("000000000000010")`, wantValue: 2},
		// {name: "BIT_FROM_STRING", saql: `BIT_FROM_STRING("11010111011101")`, wantRebuild: `BIT_FROM_STRING("11010111011101")`, wantValue: 13789},
		// {name: "BIT_FROM_STRING", saql: `BIT_FROM_STRING("100000000000000000000")`, wantRebuild: `BIT_FROM_STRING("100000000000000000000")`, wantValue: 1048756},
		// {name: "BIT_NEGATE", saql: `BIT_NEGATE(0, 8)`, wantRebuild: `BIT_NEGATE(0, 8)`, wantValue: 255},
		// {name: "BIT_NEGATE", saql: `BIT_NEGATE(0, 10)`, wantRebuild: `BIT_NEGATE(0, 10)`, wantValue: 1023},
		// {name: "BIT_NEGATE", saql: `BIT_NEGATE(3, 4)`, wantRebuild: `BIT_NEGATE(3, 4)`, wantValue: 12},
		// {name: "BIT_NEGATE", saql: `BIT_NEGATE(446359921, 32)`, wantRebuild: `BIT_NEGATE(446359921, 32)`, wantValue: 3848607374},
		// {name: "BIT_OR", saql: `BIT_OR([1, 4, 8, 16])`, wantRebuild: `BIT_OR([1, 4, 8, 16])`, wantValue: 29},
		// {name: "BIT_OR", saql: `BIT_OR([3, 7, 63])`, wantRebuild: `BIT_OR([3, 7, 63])`, wantValue: 63},
		// {name: "BIT_OR", saql: `BIT_OR([255, 127, null, 63])`, wantRebuild: `BIT_OR([255, 127, null, 63])`, wantValue: 255},
		// {name: "BIT_OR", saql: `BIT_OR(255, 127)`, wantRebuild: `BIT_OR(255, 127)`, wantValue: 255},
		// {name: "BIT_OR", saql: `BIT_OR("foo")`, wantRebuild: `BIT_OR("foo")`, wantValue: nil},
		// {name: "BIT_POPCOUNT", saql: `BIT_POPCOUNT(0)`, wantRebuild: `BIT_POPCOUNT(0)`, wantValue: 0},
		// {name: "BIT_POPCOUNT", saql: `BIT_POPCOUNT(255)`, wantRebuild: `BIT_POPCOUNT(255)`, wantValue: 8},
		// {name: "BIT_POPCOUNT", saql: `BIT_POPCOUNT(69399252)`, wantRebuild: `BIT_POPCOUNT(69399252)`, wantValue: 12},
		// {name: "BIT_POPCOUNT", saql: `BIT_POPCOUNT("foo")`, wantRebuild: `BIT_POPCOUNT("foo")`, wantValue: nil},
		// {name: "BIT_SHIFT_LEFT", saql: `BIT_SHIFT_LEFT(0, 1, 8)`, wantRebuild: `BIT_SHIFT_LEFT(0, 1, 8)`, wantValue: 0},
		// {name: "BIT_SHIFT_LEFT", saql: `BIT_SHIFT_LEFT(7, 1, 16)`, wantRebuild: `BIT_SHIFT_LEFT(7, 1, 16)`, wantValue: 14},
		// {name: "BIT_SHIFT_LEFT", saql: `BIT_SHIFT_LEFT(2, 10, 16)`, wantRebuild: `BIT_SHIFT_LEFT(2, 10, 16)`, wantValue: 2048},
		// {name: "BIT_SHIFT_LEFT", saql: `BIT_SHIFT_LEFT(878836, 16, 32)`, wantRebuild: `BIT_SHIFT_LEFT(878836, 16, 32)`, wantValue: 1760821248},
		// {name: "BIT_SHIFT_RIGHT", saql: `BIT_SHIFT_RIGHT(0, 1, 8)`, wantRebuild: `BIT_SHIFT_RIGHT(0, 1, 8)`, wantValue: 0},
		// {name: "BIT_SHIFT_RIGHT", saql: `BIT_SHIFT_RIGHT(33, 1, 16)`, wantRebuild: `BIT_SHIFT_RIGHT(33, 1, 16)`, wantValue: 16},
		// {name: "BIT_SHIFT_RIGHT", saql: `BIT_SHIFT_RIGHT(65536, 13, 16)`, wantRebuild: `BIT_SHIFT_RIGHT(65536, 13, 16)`, wantValue: 8},
		// {name: "BIT_SHIFT_RIGHT", saql: `BIT_SHIFT_RIGHT(878836, 4, 32)`, wantRebuild: `BIT_SHIFT_RIGHT(878836, 4, 32)`, wantValue: 54927},
		// {name: "BIT_TEST", saql: `BIT_TEST(0, 3)`, wantRebuild: `BIT_TEST(0, 3)`, wantValue: false},
		// {name: "BIT_TEST", saql: `BIT_TEST(255, 0)`, wantRebuild: `BIT_TEST(255, 0)`, wantValue: true},
		// {name: "BIT_TEST", saql: `BIT_TEST(7, 2)`, wantRebuild: `BIT_TEST(7, 2)`, wantValue: true},
		// {name: "BIT_TEST", saql: `BIT_TEST(255, 8)`, wantRebuild: `BIT_TEST(255, 8)`, wantValue: false},
		// {name: "BIT_TO_STRING", saql: `BIT_TO_STRING(7, 4)`, wantRebuild: `BIT_TO_STRING(7, 4)`, wantValue: "0111"},
		// {name: "BIT_TO_STRING", saql: `BIT_TO_STRING(255, 8)`, wantRebuild: `BIT_TO_STRING(255, 8)`, wantValue: "11111111"},
		// {name: "BIT_TO_STRING", saql: `BIT_TO_STRING(60, 8)`, wantRebuild: `BIT_TO_STRING(60, 8)`, wantValue: "00011110"},
		// {name: "BIT_TO_STRING", saql: `BIT_TO_STRING(1048576, 32)`, wantRebuild: `BIT_TO_STRING(1048576, 32)`, wantValue: "00000000000100000000000000000000"},
		// {name: "BIT_XOR", saql: `BIT_XOR([1, 4, 8, 16])`, wantRebuild: `BIT_XOR([1, 4, 8, 16])`, wantValue: 29},
		// {name: "BIT_XOR", saql: `BIT_XOR([3, 7, 63])`, wantRebuild: `BIT_XOR([3, 7, 63])`, wantValue: 59},
		// {name: "BIT_XOR", saql: `BIT_XOR([255, 127, null, 63])`, wantRebuild: `BIT_XOR([255, 127, null, 63])`, wantValue: 191},
		// {name: "BIT_XOR", saql: `BIT_XOR(255, 257)`, wantRebuild: `BIT_XOR(255, 257)`, wantValue: 510},
		// {name: "BIT_XOR", saql: `BIT_XOR("foo")`, wantRebuild: `BIT_XOR("foo")`, wantValue: nil},

		// https://www.arangodb.com/docs/3.7/aql/functions-date.html
		// DATE_TIMESTAMP("2014-05-07T14:19:09.522")
		// DATE_TIMESTAMP("2014-05-07T14:19:09.522Z")
		// DATE_TIMESTAMP("2014-05-07 14:19:09.522")
		// DATE_TIMESTAMP("2014-05-07 14:19:09.522Z")
		// DATE_TIMESTAMP(2014, 5, 7, 14, 19, 9, 522)
		// DATE_TIMESTAMP(1399472349522)
		// DATE_ISO8601("2014-05-07T14:19:09.522Z")
		// DATE_ISO8601("2014-05-07 14:19:09.522Z")
		// DATE_ISO8601(2014, 5, 7, 14, 19, 9, 522)
		// DATE_ISO8601(1399472349522)
		// {name: "DATE_TIMESTAMP", saql: `DATE_TIMESTAMP(2016, 12, -1)`, wantRebuild: `DATE_TIMESTAMP(2016, 12, -1)`, wantValue: nil},
		// {name: "DATE_TIMESTAMP", saql: `DATE_TIMESTAMP(2016, 2, 32)`, wantRebuild: `DATE_TIMESTAMP(2016, 2, 32)`, wantValue: 1456963200000},
		// {name: "DATE_TIMESTAMP", saql: `DATE_TIMESTAMP(1970, 1, 1, 26)`, wantRebuild: `DATE_TIMESTAMP(1970, 1, 1, 26)`, wantValue: 93600000},
		// {name: "DATE_TRUNC", saql: `DATE_TRUNC('2017-02-03', 'month')`, wantRebuild: `DATE_TRUNC('2017-02-03', 'month')`, wantValue: "2017-02-01T00:00:00.000Z"},
		// {name: "DATE_TRUNC", saql: `DATE_TRUNC('2017-02-03 04:05:06', 'hours')`, wantRebuild: `DATE_TRUNC('2017-02-03 04:05:06', 'hours')`, wantValue: "2017-02-03 04:00:00.000Z"},
		// {name: "DATE_ROUND", saql: `DATE_ROUND('2000-04-28T11:11:11.111Z', 1, 'day')`, wantRebuild: `DATE_ROUND('2000-04-28T11:11:11.111Z', 1, 'day')`, wantValue: "2000-04-28T00:00:00.000Z"},
		// {name: "DATE_ROUND", saql: `DATE_ROUND('2000-04-10T11:39:29Z', 15, 'minutes')`, wantRebuild: `DATE_ROUND('2000-04-10T11:39:29Z', 15, 'minutes')`, wantValue: "2000-04-10T11:30:00.000Z"},
		// {name: "DATE_FORMAT", saql: `DATE_FORMAT(DATE_NOW(), "%q/%yyyy")`, wantRebuild: `DATE_FORMAT(DATE_NOW(), "%q/%yyyy")`},
		// {name: "DATE_FORMAT", saql: `DATE_FORMAT(DATE_NOW(), "%dd.%mm.%yyyy %hh:%ii:%ss,%fff")`, wantRebuild: `DATE_FORMAT(DATE_NOW(), "%dd.%mm.%yyyy %hh:%ii:%ss,%fff")`, wantValue: "18.09.2015 15:30:49,374"},
		// {name: "DATE_FORMAT", saql: `DATE_FORMAT("1969", "Summer of '%yy")`, wantRebuild: `DATE_FORMAT("1969", "Summer of '%yy")`, wantValue: "Summer of '69"},
		// {name: "DATE_FORMAT", saql: `DATE_FORMAT("2016", "%%l = %l")`, wantRebuild: `DATE_FORMAT("2016", "%%l = %l")`, wantValue: "%l = 1"},
		// {name: "DATE_FORMAT", saql: `DATE_FORMAT("2016-03-01", "%xxx%")`, wantRebuild: `DATE_FORMAT("2016-03-01", "%xxx%")`, wantValue: "063, trailing % ignored"},
		// {name: "DATE_ADD", saql: `DATE_ADD(DATE_NOW(), -1, "day")`, wantRebuild: `DATE_ADD(DATE_NOW(), -1, "day")`, wantValue: "yesterday; also see DATE_SUBTRACT()"},
		// {name: "DATE_ADD", saql: `DATE_ADD(DATE_NOW(), 3, "months")`, wantRebuild: `DATE_ADD(DATE_NOW(), 3, "months")`, wantValue: "in three months"},
		// {name: "DATE_ADD", saql: `DATE_ADD(DATE_ADD("2015-04-01", 5, "years"), 1, "month")`, wantRebuild: `DATE_ADD(DATE_ADD("2015-04-01", 5, "years"), 1, "month")`, wantValue: "May 1st 2020"},
		// {name: "DATE_ADD", saql: `DATE_ADD("2015-04-01", 12*5 + 1, "months")`, wantRebuild: `DATE_ADD("2015-04-01", 12*5 + 1, "months")`, wantValue: "also May 1st 2020"},
		// {name: "DATE_ADD", saql: `DATE_ADD(DATE_TIMESTAMP(DATE_YEAR(DATE_NOW()), 12, 24), -4, "years")`, wantRebuild: `DATE_ADD(DATE_TIMESTAMP(DATE_YEAR(DATE_NOW()), 12, 24), -4, "years")`, wantValue: "Christmas four years ago"},
		// {name: "DATE_ADD", saql: `DATE_ADD(DATE_ADD("2016-02", "month", 1), -1, "day")`, wantRebuild: `DATE_ADD(DATE_ADD("2016-02", "month", 1), -1, "day")`, wantValue: "last day of February (29th, because 2016 is a leap year!)"},
		// {name: "DATE_ADD", saql: `DATE_ADD(DATE_NOW(), "P1Y")`, wantRebuild: `DATE_ADD(DATE_NOW(), "P1Y")`},
		// {name: "DATE_ADD", saql: `DATE_ADD(DATE_NOW(), "P3M2W")`, wantRebuild: `DATE_ADD(DATE_NOW(), "P3M2W")`},
		// {name: "DATE_ADD", saql: `DATE_ADD(DATE_NOW(), "P5DT26H")`, wantRebuild: `DATE_ADD(DATE_NOW(), "P5DT26H")`},
		// {name: "DATE_ADD", saql: `DATE_ADD("2000-01-01", "PT4H")`, wantRebuild: `DATE_ADD("2000-01-01", "PT4H")`},
		// {name: "DATE_ADD", saql: `DATE_ADD("2000-01-01", "PT30M44.4S"`, wantRebuild: `DATE_ADD("2000-01-01", "PT30M44.4S"`},
		// {name: "DATE_ADD", saql: `DATE_ADD("2000-01-01", "P1Y2M3W4DT5H6M7.89S"`, wantRebuild: `DATE_ADD("2000-01-01", "P1Y2M3W4DT5H6M7.89S"`},
		// {name: "DATE_SUBTRACT", saql: `DATE_SUBTRACT(DATE_NOW(), 1, "day")`, wantRebuild: `DATE_SUBTRACT(DATE_NOW(), 1, "day")`},
		// {name: "DATE_SUBTRACT", saql: `DATE_SUBTRACT(DATE_TIMESTAMP(DATE_YEAR(DATE_NOW()), 12, 24), 4, "years")`, wantRebuild: `DATE_SUBTRACT(DATE_TIMESTAMP(DATE_YEAR(DATE_NOW()), 12, 24), 4, "years")`},
		// {name: "DATE_SUBTRACT", saql: `DATE_SUBTRACT(DATE_ADD("2016-02", "month", 1), 1, "day")`, wantRebuild: `DATE_SUBTRACT(DATE_ADD("2016-02", "month", 1), 1, "day")`},
		// {name: "DATE_SUBTRACT", saql: `DATE_SUBTRACT(DATE_NOW(), "P4D")`, wantRebuild: `DATE_SUBTRACT(DATE_NOW(), "P4D")`},
		// {name: "DATE_SUBTRACT", saql: `DATE_SUBTRACT(DATE_NOW(), "PT1H3M")`, wantRebuild: `DATE_SUBTRACT(DATE_NOW(), "PT1H3M")`},
		// DATE_COMPARE("1985-04-04", DATE_NOW(), "months", "days")
		// DATE_COMPARE("1984-02-29", DATE_NOW(), "months", "days")
		// DATE_COMPARE("2001-01-01T15:30:45.678Z", "2001-01-01T08:08:08.008Z", "years", "days")

		// https://www.arangodb.com/docs/3.7/aql/functions-document.html
		{name: "ATTRIBUTES", saql: `ATTRIBUTES({"foo": "bar", "_key": "123", "_custom": "yes"})`, wantRebuild: `ATTRIBUTES({"foo": "bar", "_key": "123", "_custom": "yes"})`, wantValue: jsonParse(`["_custom", "_key", "foo"]`)},
		{name: "ATTRIBUTES", saql: `ATTRIBUTES({"foo": "bar", "_key": "123", "_custom": "yes"}, true)`, wantRebuild: `ATTRIBUTES({"foo": "bar", "_key": "123", "_custom": "yes"}, true)`, wantValue: jsonParse(`["foo"]`)},
		{name: "ATTRIBUTES", saql: `ATTRIBUTES({"foo": "bar", "_key": "123", "_custom": "yes"}, false, true)`, wantRebuild: `ATTRIBUTES({"foo": "bar", "_key": "123", "_custom": "yes"}, false, true)`, wantValue: jsonParse(`["_custom", "_key", "foo"]`)},
		{name: "HAS", saql: `HAS({name: "Jane"}, "name")`, wantRebuild: `HAS({name: "Jane"}, "name")`, wantValue: true},
		{name: "HAS", saql: `HAS({name: "Jane"}, "age")`, wantRebuild: `HAS({name: "Jane"}, "age")`, wantValue: false},
		{name: "HAS", saql: `HAS({name: null}, "name")`, wantRebuild: `HAS({name: null}, "name")`, wantValue: true},
		// KEEP(doc, "firstname", "name", "likes")
		// KEEP(doc, ["firstname", "name", "likes"])
		// MATCHES({name: "jane", age: 27, active: true}, {age: 27, active: true})
		// MATCHES({"test": 1}, [{"test": 1, "foo": "bar"}, {"foo": 1}, {"test": 1}], true)
		{name: "MERGE", saql: `MERGE({"user1": {"name": "Jane"}}, {"user2": {"name": "Tom"}})`, wantRebuild: `MERGE({"user1": {"name": "Jane"}}, {"user2": {"name": "Tom"}})`, wantValue: jsonParse(`{"user1": {"name": "Jane"}, "user2": {"name": "Tom"}}`)},
		{name: "MERGE", saql: `MERGE({"users": {"name": "Jane"}}, {"users": {"name": "Tom"}})`, wantRebuild: `MERGE({"users": {"name": "Jane"}}, {"users": {"name": "Tom"}})`, wantValue: jsonParse(`{"users": {"name": "Tom"}}`)},
		{name: "MERGE", saql: `MERGE([{foo: "bar"}, {quux: "quetzalcoatl", ruled: true}, {bar: "baz", foo: "done"}])`, wantRebuild: `MERGE([{foo: "bar"}, {quux: "quetzalcoatl", ruled: true}, {bar: "baz", foo: "done"}])`, wantValue: jsonParse(`{"foo": "done", "quux": "quetzalcoatl", "ruled": true, "bar": "baz"}`)},
		{name: "MERGE_RECURSIVE", saql: `MERGE_RECURSIVE({"user-1": {"name": "Jane", "livesIn": {"city": "LA"}}}, {"user-1": {"age": 42, "livesIn": {"state": "CA"}}})`, wantRebuild: `MERGE_RECURSIVE({"user-1": {"name": "Jane", "livesIn": {"city": "LA"}}}, {"user-1": {"age": 42, "livesIn": {"state": "CA"}}})`, wantValue: jsonParse(`{"user-1": {"name": "Jane", "livesIn": {"city": "LA", "state": "CA"}, "age": 42}}`)},
		// {name: "TRANSLATE", saql: `TRANSLATE("FR", {US: "United States", UK: "United Kingdom", FR: "France"})`, wantRebuild: `TRANSLATE("FR", {US: "United States", UK: "United Kingdom", FR: "France"})`, wantValue: "France"},
		// {name: "TRANSLATE", saql: `TRANSLATE(42, {foo: "bar", bar: "baz"})`, wantRebuild: `TRANSLATE(42, {foo: "bar", bar: "baz"})`, wantValue: 42},
		// {name: "TRANSLATE", saql: `TRANSLATE(42, {foo: "bar", bar: "baz"}, "not found!")`, wantRebuild: `TRANSLATE(42, {foo: "bar", bar: "baz"}, "not found!")`, wantValue: "not found!"},
		// UNSET(doc, "_id", "_key", "foo", "bar")
		// UNSET(doc, ["_id", "_key", "foo", "bar"])
		// UNSET_RECURSIVE(doc, "_id", "_key", "foo", "bar")
		// UNSET_RECURSIVE(doc, ["_id", "_key", "foo", "bar"])
		{name: "VALUES", saql: `VALUES({"_key": "users/jane", "name": "Jane", "age": 35})`, wantRebuild: `VALUES({"_key": "users/jane", "name": "Jane", "age": 35})`, wantValue: jsonParse(`[35, "Jane", "users/jane"]`)},
		{name: "VALUES", saql: `VALUES({"_key": "users/jane", "name": "Jane", "age": 35}, true)`, wantRebuild: `VALUES({"_key": "users/jane", "name": "Jane", "age": 35}, true)`, wantValue: jsonParse(`[35, "Jane"]`)},
		// {name: "ZIP", saql: `ZIP(["name", "active", "hobbies"], ["some user", true, ["swimming", "riding"]])`, wantRebuild: `ZIP(["name", "active", "hobbies"], ["some user", true, ["swimming", "riding"]])`, wantValue: jsonParse(`{"name": "some user", "active": true, "hobbies": ["swimming", "riding"]}`)},

		// https://www.arangodb.com/docs/3.7/aql/functions-numeric.html
		{name: "ABS", saql: `ABS(-5)`, wantRebuild: `ABS(-5)`, wantValue: 5},
		{name: "ABS", saql: `ABS(+5)`, wantRebuild: `ABS(5)`, wantValue: 5},
		{name: "ABS", saql: `ABS(3.5)`, wantRebuild: `ABS(3.5)`, wantValue: 3.5},
		{name: "ACOS", saql: `ACOS(-1)`, wantRebuild: `ACOS(-1)`, wantValue: 3.141592653589793},
		{name: "ACOS", saql: `ACOS(0)`, wantRebuild: `ACOS(0)`, wantValue: 1.5707963267948966},
		{name: "ACOS", saql: `ACOS(1)`, wantRebuild: `ACOS(1)`, wantValue: 0},
		{name: "ACOS", saql: `ACOS(2)`, wantRebuild: `ACOS(2)`, wantValue: nil},
		{name: "ASIN", saql: `ASIN(1)`, wantRebuild: `ASIN(1)`, wantValue: 1.5707963267948966},
		{name: "ASIN", saql: `ASIN(0)`, wantRebuild: `ASIN(0)`, wantValue: 0},
		{name: "ASIN", saql: `ASIN(-1)`, wantRebuild: `ASIN(-1)`, wantValue: -1.5707963267948966},
		{name: "ASIN", saql: `ASIN(2)`, wantRebuild: `ASIN(2)`, wantValue: nil},
		{name: "ATAN", saql: `ATAN(-1)`, wantRebuild: `ATAN(-1)`, wantValue: -0.7853981633974483},
		{name: "ATAN", saql: `ATAN(0)`, wantRebuild: `ATAN(0)`, wantValue: 0},
		{name: "ATAN", saql: `ATAN(10)`, wantRebuild: `ATAN(10)`, wantValue: 1.4711276743037347},
		{name: "AVERAGE", saql: `AVERAGE([5, 2, 9, 2])`, wantRebuild: `AVERAGE([5, 2, 9, 2])`, wantValue: 4.5},
		{name: "AVERAGE", saql: `AVERAGE([-3, -5, 2])`, wantRebuild: `AVERAGE([-3, -5, 2])`, wantValue: -2},
		{name: "AVERAGE", saql: `AVERAGE([999, 80, 4, 4, 4, 3, 3, 3])`, wantRebuild: `AVERAGE([999, 80, 4, 4, 4, 3, 3, 3])`, wantValue: 137.5},
		{name: "CEIL", saql: `CEIL(2.49)`, wantRebuild: `CEIL(2.49)`, wantValue: 3},
		{name: "CEIL", saql: `CEIL(2.50)`, wantRebuild: `CEIL(2.50)`, wantValue: 3},
		{name: "CEIL", saql: `CEIL(-2.50)`, wantRebuild: `CEIL(-2.50)`, wantValue: -2},
		{name: "CEIL", saql: `CEIL(-2.51)`, wantRebuild: `CEIL(-2.51)`, wantValue: -2},
		{name: "COS", saql: `COS(1)`, wantRebuild: `COS(1)`, wantValue: 0.5403023058681398},
		{name: "COS", saql: `COS(0)`, wantRebuild: `COS(0)`, wantValue: 1},
		{name: "COS", saql: `COS(-3.141592653589783)`, wantRebuild: `COS(-3.141592653589783)`, wantValue: -1},
		{name: "COS", saql: `COS(RADIANS(45))`, wantRebuild: `COS(RADIANS(45))`, wantValue: 0.7071067811865476},
		{name: "DEGREES", saql: `DEGREES(0.7853981633974483)`, wantRebuild: `DEGREES(0.7853981633974483)`, wantValue: 45},
		{name: "DEGREES", saql: `DEGREES(0)`, wantRebuild: `DEGREES(0)`, wantValue: 0},
		{name: "DEGREES", saql: `DEGREES(3.141592653589793)`, wantRebuild: `DEGREES(3.141592653589793)`, wantValue: 180},
		{name: "EXP", saql: `EXP(1)`, wantRebuild: `EXP(1)`, wantValue: 2.718281828459045},
		{name: "EXP", saql: `EXP(10)`, wantRebuild: `EXP(10)`, wantValue: 22026.46579480671},
		{name: "EXP", saql: `EXP(0)`, wantRebuild: `EXP(0)`, wantValue: 1},
		{name: "EXP2", saql: `EXP2(16)`, wantRebuild: `EXP2(16)`, wantValue: 65536},
		{name: "EXP2", saql: `EXP2(1)`, wantRebuild: `EXP2(1)`, wantValue: 2},
		{name: "EXP2", saql: `EXP2(0)`, wantRebuild: `EXP2(0)`, wantValue: 1},
		{name: "FLOOR", saql: `FLOOR(2.49)`, wantRebuild: `FLOOR(2.49)`, wantValue: 2},
		{name: "FLOOR", saql: `FLOOR(2.50)`, wantRebuild: `FLOOR(2.50)`, wantValue: 2},
		{name: "FLOOR", saql: `FLOOR(-2.50)`, wantRebuild: `FLOOR(-2.50)`, wantValue: -3},
		{name: "FLOOR", saql: `FLOOR(-2.51)`, wantRebuild: `FLOOR(-2.51)`, wantValue: -3},
		{name: "LOG", saql: `LOG(2.718281828459045)`, wantRebuild: `LOG(2.718281828459045)`, wantValue: 1},
		{name: "LOG", saql: `LOG(10)`, wantRebuild: `LOG(10)`, wantValue: 2.302585092994046},
		{name: "LOG", saql: `LOG(0)`, wantRebuild: `LOG(0)`, wantValue: nil},
		{name: "LOG2", saql: `LOG2(1024)`, wantRebuild: `LOG2(1024)`, wantValue: 10},
		{name: "LOG2", saql: `LOG2(8)`, wantRebuild: `LOG2(8)`, wantValue: 3},
		{name: "LOG2", saql: `LOG2(0)`, wantRebuild: `LOG2(0)`, wantValue: nil},
		{name: "LOG10", saql: `LOG10(10000)`, wantRebuild: `LOG10(10000)`, wantValue: 4},
		{name: "LOG10", saql: `LOG10(10)`, wantRebuild: `LOG10(10)`, wantValue: 1},
		{name: "LOG10", saql: `LOG10(0)`, wantRebuild: `LOG10(0)`, wantValue: nil},
		{name: "MAX", saql: `MAX([5, 9, -2, null, 1])`, wantRebuild: `MAX([5, 9, -2, null, 1])`, wantValue: 9},
		{name: "MAX", saql: `MAX([null, null])`, wantRebuild: `MAX([null, null])`, wantValue: nil},
		{name: "MEDIAN", saql: `MEDIAN([1, 2, 3])`, wantRebuild: `MEDIAN([1, 2, 3])`, wantValue: 2},
		{name: "MEDIAN", saql: `MEDIAN([1, 2, 3, 4])`, wantRebuild: `MEDIAN([1, 2, 3, 4])`, wantValue: 2.5},
		{name: "MEDIAN", saql: `MEDIAN([4, 2, 3, 1])`, wantRebuild: `MEDIAN([4, 2, 3, 1])`, wantValue: 2.5},
		{name: "MEDIAN", saql: `MEDIAN([999, 80, 4, 4, 4, 3, 3, 3])`, wantRebuild: `MEDIAN([999, 80, 4, 4, 4, 3, 3, 3])`, wantValue: 4},
		{name: "MIN", saql: `MIN([5, 9, -2, null, 1])`, wantRebuild: `MIN([5, 9, -2, null, 1])`, wantValue: -2},
		{name: "MIN", saql: `MIN([null, null])`, wantRebuild: `MIN([null, null])`, wantValue: nil},
		// {name: "PERCENTILE", saql: `PERCENTILE([1, 2, 3, 4], 50)`, wantRebuild: `PERCENTILE([1, 2, 3, 4], 50)`, wantValue: 2},
		// {name: "PERCENTILE", saql: `PERCENTILE([1, 2, 3, 4], 50, "rank")`, wantRebuild: `PERCENTILE([1, 2, 3, 4], 50, "rank")`, wantValue: 2},
		// {name: "PERCENTILE", saql: `PERCENTILE([1, 2, 3, 4], 50, "interpolation")`, wantRebuild: `PERCENTILE([1, 2, 3, 4], 50, "interpolation")`, wantValue: 2.5},
		{name: "PI", saql: `PI()`, wantRebuild: `PI()`, wantValue: 3.141592653589793},
		{name: "POW", saql: `POW(2, 4)`, wantRebuild: `POW(2, 4)`, wantValue: 16},
		{name: "POW", saql: `POW(5, -1)`, wantRebuild: `POW(5, -1)`, wantValue: 0.2},
		{name: "POW", saql: `POW(5, 0)`, wantRebuild: `POW(5, 0)`, wantValue: 1},
		{name: "PRODUCT", saql: `PRODUCT([1, 2, 3, 4])`, wantRebuild: `PRODUCT([1, 2, 3, 4])`, wantValue: 24},
		{name: "PRODUCT", saql: `PRODUCT([null, -5, 6])`, wantRebuild: `PRODUCT([null, -5, 6])`, wantValue: -30},
		{name: "PRODUCT", saql: `PRODUCT([])`, wantRebuild: `PRODUCT([])`, wantValue: 1},
		{name: "RADIANS", saql: `RADIANS(180)`, wantRebuild: `RADIANS(180)`, wantValue: 3.141592653589793},
		{name: "RADIANS", saql: `RADIANS(90)`, wantRebuild: `RADIANS(90)`, wantValue: 1.5707963267948966},
		{name: "RADIANS", saql: `RADIANS(0)`, wantRebuild: `RADIANS(0)`, wantValue: 0},
		// {name: "RAND", saql: `RAND()`, wantRebuild: `RAND()`, wantValue: 0.3503170117504508},
		// {name: "RAND", saql: `RAND()`, wantRebuild: `RAND()`, wantValue: 0.6138226173882478},
		{name: "RANGE", saql: `RANGE(1, 4)`, wantRebuild: `RANGE(1, 4)`, wantValue: []any{float64(1), float64(2), float64(3), float64(4)}},
		{name: "RANGE", saql: `RANGE(1, 4, 2)`, wantRebuild: `RANGE(1, 4, 2)`, wantValue: []any{float64(1), float64(3)}},
		{name: "RANGE", saql: `RANGE(1, 4, 3)`, wantRebuild: `RANGE(1, 4, 3)`, wantValue: []any{float64(1), float64(4)}},
		{name: "RANGE", saql: `RANGE(1.5, 2.5)`, wantRebuild: `RANGE(1.5, 2.5)`, wantValue: []any{float64(1), float64(2)}},
		{name: "RANGE", saql: `RANGE(1.5, 2.5, 1)`, wantRebuild: `RANGE(1.5, 2.5, 1)`, wantValue: []any{1.5, 2.5}},
		{name: "RANGE", saql: `RANGE(1.5, 2.5, 0.5)`, wantRebuild: `RANGE(1.5, 2.5, 0.5)`, wantValue: []any{1.5, 2.0, 2.5}},
		{name: "RANGE", saql: `RANGE(-0.75, 1.1, 0.5)`, wantRebuild: `RANGE(-0.75, 1.1, 0.5)`, wantValue: []any{-0.75, -0.25, 0.25, 0.75}},
		{name: "ROUND", saql: `ROUND(2.49)`, wantRebuild: `ROUND(2.49)`, wantValue: 2},
		{name: "ROUND", saql: `ROUND(2.50)`, wantRebuild: `ROUND(2.50)`, wantValue: 3},
		{name: "ROUND", saql: `ROUND(-2.50)`, wantRebuild: `ROUND(-2.50)`, wantValue: -2},
		{name: "ROUND", saql: `ROUND(-2.51)`, wantRebuild: `ROUND(-2.51)`, wantValue: -3},
		{name: "SQRT", saql: `SQRT(9)`, wantRebuild: `SQRT(9)`, wantValue: 3},
		{name: "SQRT", saql: `SQRT(2)`, wantRebuild: `SQRT(2)`, wantValue: 1.4142135623730951},
		{name: "POW", saql: `POW(4096, 1/4)`, wantRebuild: `POW(4096, 1 / 4)`, wantValue: 8},
		{name: "POW", saql: `POW(27, 1/3)`, wantRebuild: `POW(27, 1 / 3)`, wantValue: 3},
		{name: "POW", saql: `POW(9, 1/2)`, wantRebuild: `POW(9, 1 / 2)`, wantValue: 3},
		// {name: "STDDEV_POPULATION", saql: `STDDEV_POPULATION([1, 3, 6, 5, 2])`, wantRebuild: `STDDEV_POPULATION([1, 3, 6, 5, 2])`, wantValue: 1.854723699099141},
		// {name: "STDDEV_SAMPLE", saql: `STDDEV_SAMPLE([1, 3, 6, 5, 2])`, wantRebuild: `STDDEV_SAMPLE([1, 3, 6, 5, 2])`, wantValue: 2.0736441353327724},
		{name: "SUM", saql: `SUM([1, 2, 3, 4])`, wantRebuild: `SUM([1, 2, 3, 4])`, wantValue: 10},
		{name: "SUM", saql: `SUM([null, -5, 6])`, wantRebuild: `SUM([null, -5, 6])`, wantValue: 1},
		{name: "SUM", saql: `SUM([])`, wantRebuild: `SUM([])`, wantValue: 0},
		{name: "TAN", saql: `TAN(10)`, wantRebuild: `TAN(10)`, wantValue: 0.6483608274590866},
		{name: "TAN", saql: `TAN(5)`, wantRebuild: `TAN(5)`, wantValue: -3.380515006246586},
		{name: "TAN", saql: `TAN(0)`, wantRebuild: `TAN(0)`, wantValue: 0},
		// {name: "VARIANCE_POPULATION", saql: `VARIANCE_POPULATION([1, 3, 6, 5, 2])`, wantRebuild: `VARIANCE_POPULATION([1, 3, 6, 5, 2])`, wantValue: 3.4400000000000004},
		// {name: "VARIANCE_SAMPLE", saql: `VARIANCE_SAMPLE([1, 3, 6, 5, 2])`, wantRebuild: `VARIANCE_SAMPLE([1, 3, 6, 5, 2])`, wantValue: 4.300000000000001},

		// Errors
		{name: "Function Error 1", saql: "UNKNOWN(value)", wantRebuild: "UNKNOWN(value)", wantRebuildErr: true, wantEvalErr: true, values: `{"value": true}`},
		{name: "Function Error 2", saql: "ABS(value, value2)", wantRebuild: "ABS(value, value2)", wantEvalErr: true, values: `{"value": true, "value2": false}`},
		{name: "Function Error 3", saql: `ABS("abs")`, wantRebuild: `ABS("abs")`, wantEvalErr: true},
	}
	for _, tt := range tests {
		tt := tt

		parser := &caql.Parser{}

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			expr, err := parser.Parse(tt.saql)
			if (err != nil) != tt.wantParseErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantParseErr)
				if expr != nil {
					t.Error(expr.String())
				}

				return
			}
			if err != nil {
				return
			}

			got, err := expr.String()
			if (err != nil) != tt.wantRebuildErr {
				t.Error(expr.String())
				t.Errorf("String() error = %v, wantErr %v", err, tt.wantParseErr)

				return
			}
			if err != nil {
				return
			}
			if got != tt.wantRebuild {
				t.Errorf("String() got = %v, want %v", got, tt.wantRebuild)
			}

			var myJSON map[string]any
			if tt.values != "" {
				err = json.Unmarshal([]byte(tt.values), &myJSON)
				if err != nil {
					t.Fatal(err)
				}
			}

			value, err := expr.Eval(myJSON)
			if (err != nil) != tt.wantEvalErr {
				t.Error(expr.String())
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantParseErr)

				return
			}
			if err != nil {
				return
			}

			wantValue := tt.wantValue
			if i, ok := wantValue.(int); ok {
				wantValue = float64(i)
			}

			valueFloat, ok := value.(float64)
			wantValueFloat, ok2 := wantValue.(float64)
			if ok && ok2 {
				if math.Abs(valueFloat-wantValueFloat) > 0.0001 {
					t.Error(expr.String())
					t.Errorf("Eval() got = %T %#v, want %T %#v", value, value, wantValue, wantValue)
				}
			} else {
				if !reflect.DeepEqual(value, wantValue) {
					t.Error(expr.String())
					t.Errorf("Eval() got = %T %#v, want %T %#v", value, value, wantValue, wantValue)
				}
			}
		})
	}
}

func jsonParse(s string) any {
	if s == "" {
		return nil
	}
	var j any
	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		panic(s + err.Error())
	}

	return j
}
