package caql_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/SecurityBrewery/catalyst/caql"
)

type MockSearcher struct{}

func (m MockSearcher) Search(_ string) (ids []string, err error) {
	return []string{"1", "2", "3"}, nil
}

func TestParseSAQLEval(t *testing.T) {
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
		// Custom
		{name: "Compare 1", saql: "1 <= 2", wantRebuild: "1 <= 2", wantValue: true},
		{name: "Compare 2", saql: "1 >= 2", wantRebuild: "1 >= 2", wantValue: false},
		{name: "Compare 3", saql: "1 == 2", wantRebuild: "1 == 2", wantValue: false},
		{name: "Compare 4", saql: "1 > 2", wantRebuild: "1 > 2", wantValue: false},
		{name: "Compare 5", saql: "1 < 2", wantRebuild: "1 < 2", wantValue: true},
		{name: "Compare 6", saql: "1 != 2", wantRebuild: "1 != 2", wantValue: true},

		{name: "SymbolRef 1", saql: "name", wantRebuild: "name", wantValue: false, values: `{"name": false}`},
		{name: "SymbolRef 2", saql: "d.name", wantRebuild: "d.name", wantValue: false, values: `{"d": {"name": false}}`},
		{name: "SymbolRef 3", saql: "name == false", wantRebuild: "name == false", wantValue: true, values: `{"name": false}`},
		{name: "SymbolRef Error 1", saql: "name, title", wantParseErr: true},
		{name: "SymbolRef Error 2", saql: "unknown", wantRebuild: "unknown", wantValue: false, wantEvalErr: true, values: `{}`},

		{name: "Misc 1", saql: `active == true && age < 39`, wantRebuild: `active == true AND age < 39`, wantValue: true, values: `{"active": true, "age": 2}`},
		{name: "Misc 2", saql: `(attr == 10) AND foo == 'bar' OR NOT baz`, wantRebuild: `(attr == 10) AND foo == "bar" OR NOT baz`, wantValue: false, values: `{"attr": 2, "foo": "bar", "baz": true}`},
		{name: "Misc 3", saql: `attr == 10 AND (foo == 'bar' OR foo == 'baz')`, wantRebuild: `attr == 10 AND (foo == "bar" OR foo == "baz")`, wantValue: false, values: `{"attr": 2, "foo": "bar", "baz": true}`},
		{name: "Misc 4", saql: `5 > 1 AND "a" != "b"`, wantRebuild: `5 > 1 AND "a" != "b"`, wantValue: true},

		{name: "LIKE 1", saql: `"foo" LIKE "%f%"`, wantRebuild: `"foo" LIKE "%f%"`, wantValue: true},
		{name: "LIKE 2", saql: `"foo" NOT LIKE "%f%"`, wantRebuild: `"foo" NOT LIKE "%f%"`, wantValue: false},
		{name: "LIKE 3", saql: `NOT "foo" LIKE "%f%"`, wantRebuild: `NOT "foo" LIKE "%f%"`, wantValue: false},

		{name: "Summand 1", saql: "1 + 2", wantRebuild: "1 + 2", wantValue: 3},
		{name: "Summand 2", saql: "1 - 2", wantRebuild: "1 - 2", wantValue: -1},

		{name: "Factor 1", saql: "1 * 2", wantRebuild: "1 * 2", wantValue: 2},
		{name: "Factor 2", saql: "1 / 2", wantRebuild: "1 / 2", wantValue: 0.5},
		{name: "Factor 3", saql: "1.0 / 2.0", wantRebuild: "1.0 / 2.0", wantValue: 0.5},
		{name: "Factor 4", saql: "1 % 2", wantRebuild: "1 % 2", wantValue: 1},

		{name: "Term 1", saql: "(1 + 2) * 2", wantRebuild: "(1 + 2) * 2", wantValue: 6},
		{name: "Term 2", saql: "2 * (1 + 2)", wantRebuild: "2 * (1 + 2)", wantValue: 6},

		// https://www.arangodb.com/docs/3.7/aql/fundamentals-data-types.html
		{name: "Null 1", saql: `null`, wantRebuild: "null"},
		{name: "Bool 1", saql: `true`, wantRebuild: "true", wantValue: true},
		{name: "Bool 2", saql: `false`, wantRebuild: "false", wantValue: false},
		{name: "Numeric 1", saql: "1", wantRebuild: "1", wantValue: 1},
		{name: "Numeric 2", saql: "+1", wantRebuild: "1", wantValue: 1},
		{name: "Numeric 3", saql: "42", wantRebuild: "42", wantValue: 42},
		{name: "Numeric 4", saql: "-1", wantRebuild: "-1", wantValue: -1},
		{name: "Numeric 5", saql: "-42", wantRebuild: "-42", wantValue: -42},
		{name: "Numeric 6", saql: "1.23", wantRebuild: "1.23", wantValue: 1.23},
		{name: "Numeric 7", saql: "-99.99", wantRebuild: "-99.99", wantValue: -99.99},
		{name: "Numeric 8", saql: "0.5", wantRebuild: "0.5", wantValue: 0.5},
		{name: "Numeric 9", saql: ".5", wantRebuild: ".5", wantValue: 0.5},
		{name: "Numeric 10", saql: "-4.87e103", wantRebuild: "-4.87e103", wantValue: -4.87e+103},
		{name: "Numeric 11", saql: "0b10", wantRebuild: "0b10", wantValue: 2},
		{name: "Numeric 12", saql: "0x10", wantRebuild: "0x10", wantValue: 16},
		{name: "Numeric Error 1", saql: "1.", wantParseErr: true},
		{name: "Numeric Error 2", saql: "01.23", wantParseErr: true},
		{name: "Numeric Error 3", saql: "00.23", wantParseErr: true},
		{name: "Numeric Error 4", saql: "00", wantParseErr: true},

		// {name: "String 1", saql: `"yikes!"`, wantRebuild: `"yikes!"`, wantValue: "yikes!"},
		// {name: "String 2", saql: `"don't know"`, wantRebuild: `"don't know"`, wantValue: "don't know"},
		// {name: "String 3", saql: `"this is a \"quoted\" word"`, wantRebuild: `"this is a \"quoted\" word"`, wantValue: "this is a \"quoted\" word"},
		// {name: "String 4", saql: `"this is a longer string."`, wantRebuild: `"this is a longer string."`, wantValue: "this is a longer string."},
		// {name: "String 5", saql: `"the path separator on Windows is \\"`, wantRebuild: `"the path separator on Windows is \\"`, wantValue: "the path separator on Windows is \\"},
		// {name: "String 6", saql: `'yikes!'`, wantRebuild: `"yikes!"`, wantValue: "yikes!"},
		// {name: "String 7", saql: `'don\'t know'`, wantRebuild: `"don't know"`, wantValue: "don't know"},
		// {name: "String 8", saql: `'this is a "quoted" word'`, wantRebuild: `"this is a \"quoted\" word"`, wantValue: "this is a \"quoted\" word"},
		// {name: "String 9", saql: `'this is a longer string.'`, wantRebuild: `"this is a longer string."`, wantValue: "this is a longer string."},
		// {name: "String 10", saql: `'the path separator on Windows is \\'`, wantRebuild: `"the path separator on Windows is \\"`, wantValue: `the path separator on Windows is \`},

		{name: "Array 1", saql: "[]", wantRebuild: "[]", wantValue: []any{}},
		{name: "Array 2", saql: `[true]`, wantRebuild: `[true]`, wantValue: []any{true}},
		{name: "Array 3", saql: `[1, 2, 3]`, wantRebuild: `[1, 2, 3]`, wantValue: []any{float64(1), float64(2), float64(3)}},
		{
			name: "Array 4", saql: `[-99, "yikes!", [false, ["no"], []], 1]`, wantRebuild: `[-99, "yikes!", [false, ["no"], []], 1]`,
			wantValue: []any{-99.0, "yikes!", []any{false, []any{"no"}, []any{}}, float64(1)},
		},
		{name: "Array 5", saql: `[["fox", "marshal"]]`, wantRebuild: `[["fox", "marshal"]]`, wantValue: []any{[]any{"fox", "marshal"}}},
		{name: "Array 6", saql: `[1, 2, 3,]`, wantRebuild: `[1, 2, 3]`, wantValue: []any{float64(1), float64(2), float64(3)}},

		{name: "Array Error 1", saql: "(1,2,3)", wantParseErr: true},
		{name: "Array Access 1", saql: "u.friends[0]", wantRebuild: "u.friends[0]", wantValue: 7, values: `{"u": {"friends": [7,8,9]}}`},
		{name: "Array Access 2", saql: "u.friends[2]", wantRebuild: "u.friends[2]", wantValue: 9, values: `{"u": {"friends": [7,8,9]}}`},
		{name: "Array Access 3", saql: "u.friends[-1]", wantRebuild: "u.friends[-1]", wantValue: 9, values: `{"u": {"friends": [7,8,9]}}`},
		{name: "Array Access 4", saql: "u.friends[-2]", wantRebuild: "u.friends[-2]", wantValue: 8, values: `{"u": {"friends": [7,8,9]}}`},

		{name: "Object 1", saql: "{}", wantRebuild: "{}", wantValue: map[string]any{}},
		{name: "Object 2", saql: `{a: 1}`, wantRebuild: "{a: 1}", wantValue: map[string]any{"a": float64(1)}},
		{name: "Object 3", saql: `{'a': 1}`, wantRebuild: `{'a': 1}`, wantValue: map[string]any{"a": float64(1)}},
		{name: "Object 4", saql: `{"a": 1}`, wantRebuild: `{"a": 1}`, wantValue: map[string]any{"a": float64(1)}},
		{name: "Object 5", saql: `{'return': 1}`, wantRebuild: `{'return': 1}`, wantValue: map[string]any{"return": float64(1)}},
		{name: "Object 6", saql: `{"return": 1}`, wantRebuild: `{"return": 1}`, wantValue: map[string]any{"return": float64(1)}},
		{name: "Object 9", saql: `{a: 1,}`, wantRebuild: "{a: 1}", wantValue: map[string]any{"a": float64(1)}},
		{name: "Object 10", saql: `{"a": 1,}`, wantRebuild: `{"a": 1}`, wantValue: map[string]any{"a": float64(1)}},
		// {"Object 8", "{`return`: 1}", `{"return": 1}`, true},
		// {"Object 7", "{´return´: 1}", `{"return": 1}`, true},
		{name: "Object Error 1: return is a keyword", saql: `{like: 1}`, wantParseErr: true},

		{name: "Object Access 1", saql: "u.address.city.name", wantRebuild: "u.address.city.name", wantValue: "Munich", values: `{"u": {"address": {"city": {"name": "Munich"}}}}`},
		{name: "Object Access 2", saql: "u.friends[0].name.first", wantRebuild: "u.friends[0].name.first", wantValue: "Kevin", values: `{"u": {"friends": [{"name": {"first": "Kevin"}}]}}`},
		{name: "Object Access 3", saql: `u["address"]["city"]["name"]`, wantRebuild: `u["address"]["city"]["name"]`, wantValue: "Munich", values: `{"u": {"address": {"city": {"name": "Munich"}}}}`},
		{name: "Object Access 4", saql: `u["friends"][0]["name"]["first"]`, wantRebuild: `u["friends"][0]["name"]["first"]`, wantValue: "Kevin", values: `{"u": {"friends": [{"name": {"first": "Kevin"}}]}}`},
		{name: "Object Access 5", saql: "u._key", wantRebuild: "u._key", wantValue: false, values: `{"u": {"_key": false}}`},

		// This query language does not support binds
		// https://www.arangodb.com/docs/3.7/aql/fundamentals-bind-parameters.html
		// {name: "Bind 1", saql: "u.id == @id && u.name == @name", wantRebuild: `u.id == @id AND u.name == @name`, wantValue: true},
		// {name: "Bind 2", saql: "u.id == CONCAT('prefix', @id, 'suffix') && u.name == @name", wantRebuild: `u.id == CONCAT('prefix', @id, 'suffix') AND u.name == @name`, wantValue: false},
		// {name: "Bind 3", saql: "doc.@attr.@subattr", wantRebuild: `doc.@attr.@subattr`, wantValue: true, values: `{"doc": {"@attr": {"@subattr": true}}}`},
		// {name: "Bind 4", saql: "doc[@attr][@subattr]", wantRebuild: `doc[@attr][@subattr]`, wantValue: true, values: `{"doc": {"@attr": {"@subattr": true}}}`},

		// https://www.arangodb.com/docs/3.7/aql/fundamentals-type-value-order.html
		{name: "Compare 7", saql: `null < false`, wantRebuild: `null < false`, wantValue: true},
		{name: "Compare 8", saql: `null < true`, wantRebuild: `null < true`, wantValue: true},
		{name: "Compare 9", saql: `null < 1`, wantRebuild: `null < 1`, wantValue: true},
		{name: "Compare 10", saql: `null < ''`, wantRebuild: `null < ""`, wantValue: true},
		{name: "Compare 11", saql: `null < ' '`, wantRebuild: `null < " "`, wantValue: true},
		{name: "Compare 12", saql: `null < '3'`, wantRebuild: `null < "3"`, wantValue: true},
		{name: "Compare 13", saql: `null < 'abc'`, wantRebuild: `null < "abc"`, wantValue: true},
		{name: "Compare 14", saql: `null < []`, wantRebuild: `null < []`, wantValue: true},
		{name: "Compare 15", saql: `null < {}`, wantRebuild: `null < {}`, wantValue: true},
		{name: "Compare 16", saql: `false < true`, wantRebuild: `false < true`, wantValue: true},
		{name: "Compare 17", saql: `false < 5`, wantRebuild: `false < 5`, wantValue: true},
		{name: "Compare 18", saql: `false < ''`, wantRebuild: `false < ""`, wantValue: true},
		{name: "Compare 19", saql: `false < ' '`, wantRebuild: `false < " "`, wantValue: true},
		{name: "Compare 20", saql: `false < '7'`, wantRebuild: `false < "7"`, wantValue: true},
		{name: "Compare 21", saql: `false < 'abc'`, wantRebuild: `false < "abc"`, wantValue: true},
		{name: "Compare 22", saql: `false < []`, wantRebuild: `false < []`, wantValue: true},
		{name: "Compare 23", saql: `false < {}`, wantRebuild: `false < {}`, wantValue: true},
		{name: "Compare 24", saql: `true < 9`, wantRebuild: `true < 9`, wantValue: true},
		{name: "Compare 25", saql: `true < ''`, wantRebuild: `true < ""`, wantValue: true},
		{name: "Compare 26", saql: `true < ' '`, wantRebuild: `true < " "`, wantValue: true},
		{name: "Compare 27", saql: `true < '11'`, wantRebuild: `true < "11"`, wantValue: true},
		{name: "Compare 28", saql: `true < 'abc'`, wantRebuild: `true < "abc"`, wantValue: true},
		{name: "Compare 29", saql: `true < []`, wantRebuild: `true < []`, wantValue: true},
		{name: "Compare 30", saql: `true < {}`, wantRebuild: `true < {}`, wantValue: true},
		{name: "Compare 31", saql: `13 < ''`, wantRebuild: `13 < ""`, wantValue: true},
		{name: "Compare 32", saql: `15 < ' '`, wantRebuild: `15 < " "`, wantValue: true},
		{name: "Compare 33", saql: `17 < '18'`, wantRebuild: `17 < "18"`, wantValue: true},
		{name: "Compare 34", saql: `21 < 'abc'`, wantRebuild: `21 < "abc"`, wantValue: true},
		{name: "Compare 35", saql: `23 < []`, wantRebuild: `23 < []`, wantValue: true},
		{name: "Compare 36", saql: `25 < {}`, wantRebuild: `25 < {}`, wantValue: true},
		{name: "Compare 37", saql: `'' < ' '`, wantRebuild: `"" < " "`, wantValue: true},
		{name: "Compare 38", saql: `'' < '27'`, wantRebuild: `"" < "27"`, wantValue: true},
		{name: "Compare 39", saql: `'' < 'abc'`, wantRebuild: `"" < "abc"`, wantValue: true},
		{name: "Compare 40", saql: `'' < []`, wantRebuild: `"" < []`, wantValue: true},
		{name: "Compare 41", saql: `'' < {}`, wantRebuild: `"" < {}`, wantValue: true},
		{name: "Compare 42", saql: `[] < {}`, wantRebuild: `[] < {}`, wantValue: true},
		{name: "Compare 43", saql: `[] < [29]`, wantRebuild: `[] < [29]`, wantValue: true},
		{name: "Compare 44", saql: `[1] < [2]`, wantRebuild: `[1] < [2]`, wantValue: true},
		{name: "Compare 45", saql: `[1, 2] < [2]`, wantRebuild: `[1, 2] < [2]`, wantValue: true},
		{name: "Compare 46", saql: `[99, 99] < [100]`, wantRebuild: `[99, 99] < [100]`, wantValue: true},
		{name: "Compare 47", saql: `[false] < [true]`, wantRebuild: `[false] < [true]`, wantValue: true},
		{name: "Compare 48", saql: `[false, 1] < [false, '']`, wantRebuild: `[false, 1] < [false, ""]`, wantValue: true},
		{name: "Compare 49", saql: `{} < {"a": 1}`, wantRebuild: `{} < {"a": 1}`, wantValue: true},
		{name: "Compare 50", saql: `{} == {"a": null}`, wantRebuild: `{} == {"a": null}`, wantValue: true},
		{name: "Compare 51", saql: `{"a": 1} < {"a": 2}`, wantRebuild: `{"a": 1} < {"a": 2}`, wantValue: true},
		{name: "Compare 52", saql: `{"b": 1} < {"a": 0}`, wantRebuild: `{"b": 1} < {"a": 0}`, wantValue: true},
		{name: "Compare 53", saql: `{"a": {"c": true}} < {"a": {"c": 0}}`, wantRebuild: `{"a": {"c": true}} < {"a": {"c": 0}}`, wantValue: true},
		{name: "Compare 54", saql: `{"a": {"c": true, "a": 0}} < {"a": {"c": false, "a": 1}}`, wantRebuild: `{"a": {"c": true, "a": 0}} < {"a": {"c": false, "a": 1}}`, wantValue: true},
		{name: "Compare 55", saql: `{"a": 1, "b": 2} == {"b": 2, "a": 1}`, wantRebuild: `{"a": 1, "b": 2} == {"b": 2, "a": 1}`, wantValue: true},

		// https://www.arangodb.com/docs/3.7/aql/operators.html
		{name: "Compare 56", saql: `0 == null`, wantRebuild: `0 == null`, wantValue: false},
		{name: "Compare 57", saql: `1 > 0`, wantRebuild: `1 > 0`, wantValue: true},
		{name: "Compare 58", saql: `true != null`, wantRebuild: `true != null`, wantValue: true},
		{name: "Compare 59", saql: `45 <= "yikes!"`, wantRebuild: `45 <= "yikes!"`, wantValue: true},
		{name: "Compare 60", saql: `65 != "65"`, wantRebuild: `65 != "65"`, wantValue: true},
		{name: "Compare 61", saql: `65 == 65`, wantRebuild: `65 == 65`, wantValue: true},
		{name: "Compare 62", saql: `1.23 > 1.32`, wantRebuild: `1.23 > 1.32`, wantValue: false},
		{name: "Compare 63", saql: `1.5 IN [2, 3, 1.5]`, wantRebuild: `1.5 IN [2, 3, 1.5]`, wantValue: true},
		{name: "Compare 64", saql: `"foo" IN null`, wantRebuild: `"foo" IN null`, wantValue: false},
		{name: "Compare 65", saql: `42 NOT IN [17, 40, 50]`, wantRebuild: `42 NOT IN [17, 40, 50]`, wantValue: true},
		{name: "Compare 66", saql: `"abc" == "abc"`, wantRebuild: `"abc" == "abc"`, wantValue: true},
		{name: "Compare 67", saql: `"abc" == "ABC"`, wantRebuild: `"abc" == "ABC"`, wantValue: false},
		{name: "Compare 68", saql: `"foo" LIKE "f%"`, wantRebuild: `"foo" LIKE "f%"`, wantValue: true},
		{name: "Compare 69", saql: `"foo" NOT LIKE "f%"`, wantRebuild: `"foo" NOT LIKE "f%"`, wantValue: false},
		{name: "Compare 70", saql: `"foo" =~ "^f[o].$"`, wantRebuild: `"foo" =~ "^f[o].$"`, wantValue: true},
		{name: "Compare 71", saql: `"foo" !~ "[a-z]+bar$"`, wantRebuild: `"foo" !~ "[a-z]+bar$"`, wantValue: true},

		{name: "Compare 72", saql: `"abc" LIKE "a%"`, wantRebuild: `"abc" LIKE "a%"`, wantValue: true},
		{name: "Compare 73", saql: `"abc" LIKE "_bc"`, wantRebuild: `"abc" LIKE "_bc"`, wantValue: true},
		{name: "Compare 74", saql: `"a_b_foo" LIKE "a\\_b\\_foo"`, wantRebuild: `"a_b_foo" LIKE "a\\_b\\_foo"`, wantValue: true},

		// https://www.arangodb.com/docs/3.7/aql/operators.html#array-comparison-operators
		{name: "Compare Array 1", saql: `[1, 2, 3] ALL IN [2, 3, 4]`, wantRebuild: `[1, 2, 3] ALL IN [2, 3, 4]`, wantValue: false},
		{name: "Compare Array 2", saql: `[1, 2, 3] ALL IN [1, 2, 3]`, wantRebuild: `[1, 2, 3] ALL IN [1, 2, 3]`, wantValue: true},
		{name: "Compare Array 3", saql: `[1, 2, 3] NONE IN [3]`, wantRebuild: `[1, 2, 3] NONE IN [3]`, wantValue: false},
		{name: "Compare Array 4", saql: `[1, 2, 3] NONE IN [23, 42]`, wantRebuild: `[1, 2, 3] NONE IN [23, 42]`, wantValue: true},
		{name: "Compare Array 5", saql: `[1, 2, 3] ANY IN [4, 5, 6]`, wantRebuild: `[1, 2, 3] ANY IN [4, 5, 6]`, wantValue: false},
		{name: "Compare Array 6", saql: `[1, 2, 3] ANY IN [1, 42]`, wantRebuild: `[1, 2, 3] ANY IN [1, 42]`, wantValue: true},
		{name: "Compare Array 7", saql: `[1, 2, 3] ANY == 2`, wantRebuild: `[1, 2, 3] ANY == 2`, wantValue: true},
		{name: "Compare Array 8", saql: `[1, 2, 3] ANY == 4`, wantRebuild: `[1, 2, 3] ANY == 4`, wantValue: false},
		{name: "Compare Array 9", saql: `[1, 2, 3] ANY > 0`, wantRebuild: `[1, 2, 3] ANY > 0`, wantValue: true},
		{name: "Compare Array 10", saql: `[1, 2, 3] ANY <= 1`, wantRebuild: `[1, 2, 3] ANY <= 1`, wantValue: true},
		{name: "Compare Array 11", saql: `[1, 2, 3] NONE < 99`, wantRebuild: `[1, 2, 3] NONE < 99`, wantValue: false},
		{name: "Compare Array 12", saql: `[1, 2, 3] NONE > 10`, wantRebuild: `[1, 2, 3] NONE > 10`, wantValue: true},
		{name: "Compare Array 13", saql: `[1, 2, 3] ALL > 2`, wantRebuild: `[1, 2, 3] ALL > 2`, wantValue: false},
		{name: "Compare Array 14", saql: `[1, 2, 3] ALL > 0`, wantRebuild: `[1, 2, 3] ALL > 0`, wantValue: true},
		{name: "Compare Array 15", saql: `[1, 2, 3] ALL >= 3`, wantRebuild: `[1, 2, 3] ALL >= 3`, wantValue: false},
		{name: "Compare Array 16", saql: `["foo", "bar"] ALL != "moo"`, wantRebuild: `["foo", "bar"] ALL != "moo"`, wantValue: true},
		{name: "Compare Array 17", saql: `["foo", "bar"] NONE == "bar"`, wantRebuild: `["foo", "bar"] NONE == "bar"`, wantValue: false},
		{name: "Compare Array 18", saql: `["foo", "bar"] ANY == "foo"`, wantRebuild: `["foo", "bar"] ANY == "foo"`, wantValue: true},

		// https://www.arangodb.com/docs/3.7/aql/operators.html#logical-operators
		{name: "Logical 1", saql: "active == true OR age < 39", wantRebuild: "active == true OR age < 39", wantValue: true, values: `{"active": true, "age": 4}`},
		{name: "Logical 2", saql: "active == true || age < 39", wantRebuild: "active == true OR age < 39", wantValue: true, values: `{"active": true, "age": 4}`},
		{name: "Logical 3", saql: "active == true AND age < 39", wantRebuild: "active == true AND age < 39", wantValue: true, values: `{"active": true, "age": 4}`},
		{name: "Logical 4", saql: "active == true && age < 39", wantRebuild: "active == true AND age < 39", wantValue: true, values: `{"active": true, "age": 4}`},
		{name: "Logical 5", saql: "!active", wantRebuild: "NOT active", wantValue: false, values: `{"active": true}`},
		{name: "Logical 6", saql: "NOT active", wantRebuild: "NOT active", wantValue: false, values: `{"active": true}`},
		{name: "Logical 7", saql: "not active", wantRebuild: "NOT active", wantValue: false, values: `{"active": true}`},
		{name: "Logical 8", saql: "NOT NOT active", wantRebuild: "NOT NOT active", wantValue: true, values: `{"active": true}`},

		{name: "Logical 9", saql: `u.age > 15 && u.address.city != ""`, wantRebuild: `u.age > 15 AND u.address.city != ""`, wantValue: false, values: `{"u": {"age": 2, "address": {"city": "Munich"}}}`},
		{name: "Logical 10", saql: `true || false`, wantRebuild: `true OR false`, wantValue: true},
		{name: "Logical 11", saql: `NOT u.isInvalid`, wantRebuild: `NOT u.isInvalid`, wantValue: false, values: `{"u": {"isInvalid": true}}`},
		{name: "Logical 12", saql: `1 || ! 0`, wantRebuild: `1 OR NOT 0`, wantValue: 1},

		{name: "Logical 13", saql: `25 > 1  &&  42 != 7`, wantRebuild: `25 > 1 AND 42 != 7`, wantValue: true},
		{name: "Logical 14", saql: `22 IN [23, 42]  ||  23 NOT IN [22, 7]`, wantRebuild: `22 IN [23, 42] OR 23 NOT IN [22, 7]`, wantValue: true},
		{name: "Logical 15", saql: `25 != 25`, wantRebuild: `25 != 25`, wantValue: false},

		{name: "Logical 16", saql: `1 || 7`, wantRebuild: `1 OR 7`, wantValue: 1},
		// {name: "Logical 17", saql: `null || "foo"`, wantRebuild: `null OR "foo"`, wantValue: "foo"},
		{name: "Logical 17", saql: `null || "foo"`, wantRebuild: `null OR d._key IN ["1","2","3"]`, wantValue: "foo", values: `{"d": {"_key": "1"}}`}, // eval != rebuild
		{name: "Logical 18", saql: `null && true`, wantRebuild: `null AND true`, wantValue: nil},
		{name: "Logical 19", saql: `true && 23`, wantRebuild: `true AND 23`, wantValue: 23},

		{name: "Logical 20", saql: "true == (6 < 8)", wantRebuild: "true == (6 < 8)", wantValue: true},
		{name: "Logical 21", saql: "true == 6 < 8", wantRebuild: "true == 6 < 8", wantValue: true}, // does not work in go

		// https://www.arangodb.com/docs/3.7/aql/operators.html#arithmetic-operators
		{name: "Arithmetic 1", saql: `1 + 1`, wantRebuild: `1 + 1`, wantValue: 2},
		{name: "Arithmetic 2", saql: `33 - 99`, wantRebuild: `33 - 99`, wantValue: -66},
		{name: "Arithmetic 3", saql: `12.4 * 4.5`, wantRebuild: `12.4 * 4.5`, wantValue: 55.8},
		{name: "Arithmetic 4", saql: `13.0 / 0.1`, wantRebuild: `13.0 / 0.1`, wantValue: 130.0},
		{name: "Arithmetic 5", saql: `23 % 7`, wantRebuild: `23 % 7`, wantValue: 2},
		{name: "Arithmetic 6", saql: `-15`, wantRebuild: `-15`, wantValue: -15},
		{name: "Arithmetic 7", saql: `+9.99`, wantRebuild: `9.99`, wantValue: 9.99},

		{name: "Arithmetic 8", saql: `1 + "a"`, wantRebuild: `1 + "a"`, wantValue: 1},
		{name: "Arithmetic 9", saql: `1 + "99"`, wantRebuild: `1 + "99"`, wantValue: 100},
		{name: "Arithmetic 10", saql: `1 + null`, wantRebuild: `1 + null`, wantValue: 1},
		{name: "Arithmetic 11", saql: `null + 1`, wantRebuild: `null + 1`, wantValue: 1},
		{name: "Arithmetic 12", saql: `3 + []`, wantRebuild: `3 + []`, wantValue: 3},
		{name: "Arithmetic 13", saql: `24 + [2]`, wantRebuild: `24 + [2]`, wantValue: 26},
		{name: "Arithmetic 14", saql: `24 + [2, 4]`, wantRebuild: `24 + [2, 4]`, wantValue: 24},
		{name: "Arithmetic 15", saql: `25 - null`, wantRebuild: `25 - null`, wantValue: 25},
		{name: "Arithmetic 16", saql: `17 - true`, wantRebuild: `17 - true`, wantValue: 16},
		{name: "Arithmetic 17", saql: `23 * {}`, wantRebuild: `23 * {}`, wantValue: 0},
		{name: "Arithmetic 18", saql: `5 * [7]`, wantRebuild: `5 * [7]`, wantValue: 35},
		{name: "Arithmetic 19", saql: `24 / "12"`, wantRebuild: `24 / "12"`, wantValue: 2},
		{name: "Arithmetic Error 1: Division by zero", saql: `1 / 0`, wantRebuild: `1 / 0`, wantValue: 0},

		// https://www.arangodb.com/docs/3.7/aql/operators.html#ternary-operator
		{name: "Ternary 1", saql: `u.age > 15 || u.active == true ? u.userId : null`, wantRebuild: `u.age > 15 OR u.active == true ? u.userId : null`, wantValue: 45, values: `{"u": {"active": true, "age": 2, "userId": 45}}`},
		{name: "Ternary 2", saql: `u.value ? : 'value is null, 0 or not present'`, wantRebuild: `u.value ? : "value is null, 0 or not present"`, wantValue: "value is null, 0 or not present", values: `{"u": {"value": 0}}`},

		// https://www.arangodb.com/docs/3.7/aql/operators.html#range-operator
		{name: "Range 1", saql: `2010..2013`, wantRebuild: `2010..2013`, wantValue: []float64{2010, 2011, 2012, 2013}},
		// {"Array operators 1", `u.friends[*].name`, `u.friends[*].name`, false},

		// Security
		{name: "Security 1", saql: `doc.value == 1 || true REMOVE doc IN collection //`, wantParseErr: true},
		{name: "Security 2", saql: `doc.value == 1 || true INSERT {foo: "bar"} IN collection //`, wantParseErr: true},

		// https://www.arangodb.com/docs/3.7/aql/operators.html#operator-precedence
		{name: "Precedence", saql: `2 > 15 && "a" != ""`, wantRebuild: `2 > 15 AND "a" != ""`, wantValue: false},
	}
	for _, tt := range tests {
		tt := tt
		parser := &caql.Parser{
			Searcher: &MockSearcher{},
		}

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

			if !reflect.DeepEqual(value, wantValue) {
				t.Error(expr.String())
				t.Errorf("Eval() got = %T %#v, want %T %#v", value, value, wantValue, wantValue)
			}
		})
	}
}
