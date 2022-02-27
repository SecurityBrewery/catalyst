// Code generated from CAQLParser.g4 by ANTLR 4.9.3. DO NOT EDIT.

package parser // CAQLParser

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 80, 192,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 3, 2,
	3, 2, 3, 2, 3, 3, 3, 3, 3, 3, 3, 3, 5, 3, 32, 10, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 5, 3,
	48, 10, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 5, 3, 66, 10, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 7, 3, 86, 10, 3, 12, 3, 14, 3, 89, 11, 3, 3, 4, 3, 4,
	3, 4, 3, 4, 3, 4, 3, 4, 5, 4, 97, 10, 4, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5,
	3, 5, 3, 5, 3, 5, 5, 5, 107, 10, 5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5,
	3, 5, 3, 5, 7, 5, 117, 10, 5, 12, 5, 14, 5, 120, 11, 5, 3, 6, 3, 6, 5,
	6, 124, 10, 6, 3, 7, 3, 7, 3, 7, 5, 7, 129, 10, 7, 3, 7, 3, 7, 7, 7, 133,
	10, 7, 12, 7, 14, 7, 136, 11, 7, 3, 7, 5, 7, 139, 10, 7, 3, 7, 3, 7, 3,
	8, 3, 8, 3, 9, 3, 9, 5, 9, 147, 10, 9, 3, 9, 3, 9, 7, 9, 151, 10, 9, 12,
	9, 14, 9, 154, 11, 9, 3, 9, 5, 9, 157, 10, 9, 3, 9, 3, 9, 3, 10, 3, 10,
	5, 10, 163, 10, 10, 3, 10, 3, 10, 7, 10, 167, 10, 10, 12, 10, 14, 10, 170,
	11, 10, 3, 10, 5, 10, 173, 10, 10, 3, 10, 3, 10, 3, 11, 3, 11, 3, 11, 3,
	11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 5, 11, 188, 10, 11,
	3, 12, 3, 12, 3, 12, 4, 134, 152, 4, 4, 8, 13, 2, 4, 6, 8, 10, 12, 14,
	16, 18, 20, 22, 2, 11, 3, 2, 12, 13, 3, 2, 14, 16, 3, 2, 8, 11, 3, 2, 6,
	7, 5, 2, 29, 29, 31, 31, 48, 48, 4, 2, 6, 11, 40, 40, 4, 2, 4, 5, 46, 46,
	7, 2, 36, 36, 50, 50, 58, 58, 72, 73, 75, 75, 4, 2, 71, 71, 75, 75, 2,
	216, 2, 24, 3, 2, 2, 2, 4, 31, 3, 2, 2, 2, 6, 96, 3, 2, 2, 2, 8, 106, 3,
	2, 2, 2, 10, 123, 3, 2, 2, 2, 12, 125, 3, 2, 2, 2, 14, 142, 3, 2, 2, 2,
	16, 144, 3, 2, 2, 2, 18, 160, 3, 2, 2, 2, 20, 187, 3, 2, 2, 2, 22, 189,
	3, 2, 2, 2, 24, 25, 5, 4, 3, 2, 25, 26, 7, 2, 2, 3, 26, 3, 3, 2, 2, 2,
	27, 28, 8, 3, 1, 2, 28, 32, 5, 14, 8, 2, 29, 32, 5, 8, 5, 2, 30, 32, 5,
	6, 4, 2, 31, 27, 3, 2, 2, 2, 31, 29, 3, 2, 2, 2, 31, 30, 3, 2, 2, 2, 32,
	87, 3, 2, 2, 2, 33, 34, 12, 15, 2, 2, 34, 35, 9, 2, 2, 2, 35, 86, 5, 4,
	3, 16, 36, 37, 12, 14, 2, 2, 37, 38, 9, 3, 2, 2, 38, 86, 5, 4, 3, 15, 39,
	40, 12, 13, 2, 2, 40, 41, 7, 20, 2, 2, 41, 86, 5, 4, 3, 14, 42, 43, 12,
	12, 2, 2, 43, 44, 9, 4, 2, 2, 44, 86, 5, 4, 3, 13, 45, 47, 12, 11, 2, 2,
	46, 48, 7, 49, 2, 2, 47, 46, 3, 2, 2, 2, 47, 48, 3, 2, 2, 2, 48, 49, 3,
	2, 2, 2, 49, 50, 7, 40, 2, 2, 50, 86, 5, 4, 3, 12, 51, 52, 12, 10, 2, 2,
	52, 53, 9, 5, 2, 2, 53, 86, 5, 4, 3, 11, 54, 55, 12, 9, 2, 2, 55, 56, 9,
	6, 2, 2, 56, 57, 9, 7, 2, 2, 57, 86, 5, 4, 3, 10, 58, 59, 12, 8, 2, 2,
	59, 60, 9, 6, 2, 2, 60, 61, 7, 49, 2, 2, 61, 62, 7, 40, 2, 2, 62, 86, 5,
	4, 3, 9, 63, 65, 12, 7, 2, 2, 64, 66, 7, 49, 2, 2, 65, 64, 3, 2, 2, 2,
	65, 66, 3, 2, 2, 2, 66, 67, 3, 2, 2, 2, 67, 68, 9, 8, 2, 2, 68, 86, 5,
	4, 3, 8, 69, 70, 12, 6, 2, 2, 70, 71, 7, 30, 2, 2, 71, 86, 5, 4, 3, 7,
	72, 73, 12, 5, 2, 2, 73, 74, 7, 51, 2, 2, 74, 86, 5, 4, 3, 6, 75, 76, 12,
	4, 2, 2, 76, 77, 7, 17, 2, 2, 77, 78, 5, 4, 3, 2, 78, 79, 7, 18, 2, 2,
	79, 80, 5, 4, 3, 5, 80, 86, 3, 2, 2, 2, 81, 82, 12, 3, 2, 2, 82, 83, 7,
	17, 2, 2, 83, 84, 7, 18, 2, 2, 84, 86, 5, 4, 3, 4, 85, 33, 3, 2, 2, 2,
	85, 36, 3, 2, 2, 2, 85, 39, 3, 2, 2, 2, 85, 42, 3, 2, 2, 2, 85, 45, 3,
	2, 2, 2, 85, 51, 3, 2, 2, 2, 85, 54, 3, 2, 2, 2, 85, 58, 3, 2, 2, 2, 85,
	63, 3, 2, 2, 2, 85, 69, 3, 2, 2, 2, 85, 72, 3, 2, 2, 2, 85, 75, 3, 2, 2,
	2, 85, 81, 3, 2, 2, 2, 86, 89, 3, 2, 2, 2, 87, 85, 3, 2, 2, 2, 87, 88,
	3, 2, 2, 2, 88, 5, 3, 2, 2, 2, 89, 87, 3, 2, 2, 2, 90, 91, 7, 12, 2, 2,
	91, 97, 5, 4, 3, 2, 92, 93, 7, 13, 2, 2, 93, 97, 5, 4, 3, 2, 94, 95, 7,
	49, 2, 2, 95, 97, 5, 4, 3, 2, 96, 90, 3, 2, 2, 2, 96, 92, 3, 2, 2, 2, 96,
	94, 3, 2, 2, 2, 97, 7, 3, 2, 2, 2, 98, 99, 8, 5, 1, 2, 99, 107, 7, 71,
	2, 2, 100, 107, 5, 10, 6, 2, 101, 107, 5, 12, 7, 2, 102, 103, 7, 22, 2,
	2, 103, 104, 5, 4, 3, 2, 104, 105, 7, 23, 2, 2, 105, 107, 3, 2, 2, 2, 106,
	98, 3, 2, 2, 2, 106, 100, 3, 2, 2, 2, 106, 101, 3, 2, 2, 2, 106, 102, 3,
	2, 2, 2, 107, 118, 3, 2, 2, 2, 108, 109, 12, 4, 2, 2, 109, 110, 7, 3, 2,
	2, 110, 117, 7, 71, 2, 2, 111, 112, 12, 3, 2, 2, 112, 113, 7, 26, 2, 2,
	113, 114, 5, 4, 3, 2, 114, 115, 7, 27, 2, 2, 115, 117, 3, 2, 2, 2, 116,
	108, 3, 2, 2, 2, 116, 111, 3, 2, 2, 2, 117, 120, 3, 2, 2, 2, 118, 116,
	3, 2, 2, 2, 118, 119, 3, 2, 2, 2, 119, 9, 3, 2, 2, 2, 120, 118, 3, 2, 2,
	2, 121, 124, 5, 16, 9, 2, 122, 124, 5, 18, 10, 2, 123, 121, 3, 2, 2, 2,
	123, 122, 3, 2, 2, 2, 124, 11, 3, 2, 2, 2, 125, 126, 7, 71, 2, 2, 126,
	128, 7, 22, 2, 2, 127, 129, 5, 4, 3, 2, 128, 127, 3, 2, 2, 2, 128, 129,
	3, 2, 2, 2, 129, 134, 3, 2, 2, 2, 130, 131, 7, 21, 2, 2, 131, 133, 5, 4,
	3, 2, 132, 130, 3, 2, 2, 2, 133, 136, 3, 2, 2, 2, 134, 135, 3, 2, 2, 2,
	134, 132, 3, 2, 2, 2, 135, 138, 3, 2, 2, 2, 136, 134, 3, 2, 2, 2, 137,
	139, 7, 21, 2, 2, 138, 137, 3, 2, 2, 2, 138, 139, 3, 2, 2, 2, 139, 140,
	3, 2, 2, 2, 140, 141, 7, 23, 2, 2, 141, 13, 3, 2, 2, 2, 142, 143, 9, 9,
	2, 2, 143, 15, 3, 2, 2, 2, 144, 146, 7, 26, 2, 2, 145, 147, 5, 4, 3, 2,
	146, 145, 3, 2, 2, 2, 146, 147, 3, 2, 2, 2, 147, 152, 3, 2, 2, 2, 148,
	149, 7, 21, 2, 2, 149, 151, 5, 4, 3, 2, 150, 148, 3, 2, 2, 2, 151, 154,
	3, 2, 2, 2, 152, 153, 3, 2, 2, 2, 152, 150, 3, 2, 2, 2, 153, 156, 3, 2,
	2, 2, 154, 152, 3, 2, 2, 2, 155, 157, 7, 21, 2, 2, 156, 155, 3, 2, 2, 2,
	156, 157, 3, 2, 2, 2, 157, 158, 3, 2, 2, 2, 158, 159, 7, 27, 2, 2, 159,
	17, 3, 2, 2, 2, 160, 162, 7, 24, 2, 2, 161, 163, 5, 20, 11, 2, 162, 161,
	3, 2, 2, 2, 162, 163, 3, 2, 2, 2, 163, 168, 3, 2, 2, 2, 164, 165, 7, 21,
	2, 2, 165, 167, 5, 20, 11, 2, 166, 164, 3, 2, 2, 2, 167, 170, 3, 2, 2,
	2, 168, 166, 3, 2, 2, 2, 168, 169, 3, 2, 2, 2, 169, 172, 3, 2, 2, 2, 170,
	168, 3, 2, 2, 2, 171, 173, 7, 21, 2, 2, 172, 171, 3, 2, 2, 2, 172, 173,
	3, 2, 2, 2, 173, 174, 3, 2, 2, 2, 174, 175, 7, 25, 2, 2, 175, 19, 3, 2,
	2, 2, 176, 188, 7, 71, 2, 2, 177, 178, 5, 22, 12, 2, 178, 179, 7, 18, 2,
	2, 179, 180, 5, 4, 3, 2, 180, 188, 3, 2, 2, 2, 181, 182, 7, 26, 2, 2, 182,
	183, 5, 4, 3, 2, 183, 184, 7, 27, 2, 2, 184, 185, 7, 18, 2, 2, 185, 186,
	5, 4, 3, 2, 186, 188, 3, 2, 2, 2, 187, 176, 3, 2, 2, 2, 187, 177, 3, 2,
	2, 2, 187, 181, 3, 2, 2, 2, 188, 21, 3, 2, 2, 2, 189, 190, 9, 10, 2, 2,
	190, 23, 3, 2, 2, 2, 22, 31, 47, 65, 85, 87, 96, 106, 116, 118, 123, 128,
	134, 138, 146, 152, 156, 162, 168, 172, 187,
}
var literalNames = []string{
	"", "'.'", "'=~'", "'!~'", "'=='", "'!='", "'<'", "'>'", "'<='", "'>='",
	"'+'", "'-'", "'*'", "'/'", "'%'", "'?'", "':'", "'::'", "'..'", "','",
	"'('", "')'", "'{'", "'}'", "'['", "']'",
}
var symbolicNames = []string{
	"", "DOT", "T_REGEX_MATCH", "T_REGEX_NON_MATCH", "T_EQ", "T_NE", "T_LT",
	"T_GT", "T_LE", "T_GE", "T_PLUS", "T_MINUS", "T_TIMES", "T_DIV", "T_MOD",
	"T_QUESTION", "T_COLON", "T_SCOPE", "T_RANGE", "T_COMMA", "T_OPEN", "T_CLOSE",
	"T_OBJECT_OPEN", "T_OBJECT_CLOSE", "T_ARRAY_OPEN", "T_ARRAY_CLOSE", "T_AGGREGATE",
	"T_ALL", "T_AND", "T_ANY", "T_ASC", "T_COLLECT", "T_DESC", "T_DISTINCT",
	"T_FALSE", "T_FILTER", "T_FOR", "T_GRAPH", "T_IN", "T_INBOUND", "T_INSERT",
	"T_INTO", "T_K_SHORTEST_PATHS", "T_LET", "T_LIKE", "T_LIMIT", "T_NONE",
	"T_NOT", "T_NULL", "T_OR", "T_OUTBOUND", "T_REMOVE", "T_REPLACE", "T_RETURN",
	"T_SHORTEST_PATH", "T_SORT", "T_TRUE", "T_UPDATE", "T_UPSERT", "T_WITH",
	"T_KEEP", "T_COUNT", "T_OPTIONS", "T_PRUNE", "T_SEARCH", "T_TO", "T_CURRENT",
	"T_NEW", "T_OLD", "T_STRING", "T_INT", "T_FLOAT", "T_PARAMETER", "T_QUOTED_STRING",
	"SINGLE_LINE_COMMENT", "MULTILINE_COMMENT", "SPACES", "UNEXPECTED_CHAR",
	"ERROR_RECONGNIGION",
}

var ruleNames = []string{
	"parse", "expression", "operator_unary", "reference", "compound_value",
	"function_call", "value_literal", "array", "object", "object_element",
	"object_element_name",
}

type CAQLParser struct {
	*antlr.BaseParser
}

// NewCAQLParser produces a new parser instance for the optional input antlr.TokenStream.
//
// The *CAQLParser instance produced may be reused by calling the SetInputStream method.
// The initial parser configuration is expensive to construct, and the object is not thread-safe;
// however, if used within a Golang sync.Pool, the construction cost amortizes well and the
// objects can be used in a thread-safe manner.
func NewCAQLParser(input antlr.TokenStream) *CAQLParser {
	this := new(CAQLParser)
	deserializer := antlr.NewATNDeserializer(nil)
	deserializedATN := deserializer.DeserializeFromUInt16(parserATN)
	decisionToDFA := make([]*antlr.DFA, len(deserializedATN.DecisionToState))
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "CAQLParser.g4"

	return this
}

// CAQLParser tokens.
const (
	CAQLParserEOF                 = antlr.TokenEOF
	CAQLParserDOT                 = 1
	CAQLParserT_REGEX_MATCH       = 2
	CAQLParserT_REGEX_NON_MATCH   = 3
	CAQLParserT_EQ                = 4
	CAQLParserT_NE                = 5
	CAQLParserT_LT                = 6
	CAQLParserT_GT                = 7
	CAQLParserT_LE                = 8
	CAQLParserT_GE                = 9
	CAQLParserT_PLUS              = 10
	CAQLParserT_MINUS             = 11
	CAQLParserT_TIMES             = 12
	CAQLParserT_DIV               = 13
	CAQLParserT_MOD               = 14
	CAQLParserT_QUESTION          = 15
	CAQLParserT_COLON             = 16
	CAQLParserT_SCOPE             = 17
	CAQLParserT_RANGE             = 18
	CAQLParserT_COMMA             = 19
	CAQLParserT_OPEN              = 20
	CAQLParserT_CLOSE             = 21
	CAQLParserT_OBJECT_OPEN       = 22
	CAQLParserT_OBJECT_CLOSE      = 23
	CAQLParserT_ARRAY_OPEN        = 24
	CAQLParserT_ARRAY_CLOSE       = 25
	CAQLParserT_AGGREGATE         = 26
	CAQLParserT_ALL               = 27
	CAQLParserT_AND               = 28
	CAQLParserT_ANY               = 29
	CAQLParserT_ASC               = 30
	CAQLParserT_COLLECT           = 31
	CAQLParserT_DESC              = 32
	CAQLParserT_DISTINCT          = 33
	CAQLParserT_FALSE             = 34
	CAQLParserT_FILTER            = 35
	CAQLParserT_FOR               = 36
	CAQLParserT_GRAPH             = 37
	CAQLParserT_IN                = 38
	CAQLParserT_INBOUND           = 39
	CAQLParserT_INSERT            = 40
	CAQLParserT_INTO              = 41
	CAQLParserT_K_SHORTEST_PATHS  = 42
	CAQLParserT_LET               = 43
	CAQLParserT_LIKE              = 44
	CAQLParserT_LIMIT             = 45
	CAQLParserT_NONE              = 46
	CAQLParserT_NOT               = 47
	CAQLParserT_NULL              = 48
	CAQLParserT_OR                = 49
	CAQLParserT_OUTBOUND          = 50
	CAQLParserT_REMOVE            = 51
	CAQLParserT_REPLACE           = 52
	CAQLParserT_RETURN            = 53
	CAQLParserT_SHORTEST_PATH     = 54
	CAQLParserT_SORT              = 55
	CAQLParserT_TRUE              = 56
	CAQLParserT_UPDATE            = 57
	CAQLParserT_UPSERT            = 58
	CAQLParserT_WITH              = 59
	CAQLParserT_KEEP              = 60
	CAQLParserT_COUNT             = 61
	CAQLParserT_OPTIONS           = 62
	CAQLParserT_PRUNE             = 63
	CAQLParserT_SEARCH            = 64
	CAQLParserT_TO                = 65
	CAQLParserT_CURRENT           = 66
	CAQLParserT_NEW               = 67
	CAQLParserT_OLD               = 68
	CAQLParserT_STRING            = 69
	CAQLParserT_INT               = 70
	CAQLParserT_FLOAT             = 71
	CAQLParserT_PARAMETER         = 72
	CAQLParserT_QUOTED_STRING     = 73
	CAQLParserSINGLE_LINE_COMMENT = 74
	CAQLParserMULTILINE_COMMENT   = 75
	CAQLParserSPACES              = 76
	CAQLParserUNEXPECTED_CHAR     = 77
	CAQLParserERROR_RECONGNIGION  = 78
)

// CAQLParser rules.
const (
	CAQLParserRULE_parse               = 0
	CAQLParserRULE_expression          = 1
	CAQLParserRULE_operator_unary      = 2
	CAQLParserRULE_reference           = 3
	CAQLParserRULE_compound_value      = 4
	CAQLParserRULE_function_call       = 5
	CAQLParserRULE_value_literal       = 6
	CAQLParserRULE_array               = 7
	CAQLParserRULE_object              = 8
	CAQLParserRULE_object_element      = 9
	CAQLParserRULE_object_element_name = 10
)

// IParseContext is an interface to support dynamic dispatch.
type IParseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsParseContext differentiates from other interfaces.
	IsParseContext()
}

type ParseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyParseContext() *ParseContext {
	var p = new(ParseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CAQLParserRULE_parse
	return p
}

func (*ParseContext) IsParseContext() {}

func NewParseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParseContext {
	var p = new(ParseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CAQLParserRULE_parse

	return p
}

func (s *ParseContext) GetParser() antlr.Parser { return s.parser }

func (s *ParseContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ParseContext) EOF() antlr.TerminalNode {
	return s.GetToken(CAQLParserEOF, 0)
}

func (s *ParseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CAQLParserListener); ok {
		listenerT.EnterParse(s)
	}
}

func (s *ParseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CAQLParserListener); ok {
		listenerT.ExitParse(s)
	}
}

func (p *CAQLParser) Parse() (localctx IParseContext) {
	this := p
	_ = this

	localctx = NewParseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, CAQLParserRULE_parse)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(22)
		p.expression(0)
	}
	{
		p.SetState(23)
		p.Match(CAQLParserEOF)
	}

	return localctx
}

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetEq_op returns the eq_op token.
	GetEq_op() antlr.Token

	// SetEq_op sets the eq_op token.
	SetEq_op(antlr.Token)

	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}

type ExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	eq_op  antlr.Token
}

func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CAQLParserRULE_expression
	return p
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CAQLParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) GetEq_op() antlr.Token { return s.eq_op }

func (s *ExpressionContext) SetEq_op(v antlr.Token) { s.eq_op = v }

func (s *ExpressionContext) Value_literal() IValue_literalContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IValue_literalContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IValue_literalContext)
}

func (s *ExpressionContext) Reference() IReferenceContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IReferenceContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IReferenceContext)
}

func (s *ExpressionContext) Operator_unary() IOperator_unaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IOperator_unaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IOperator_unaryContext)
}

func (s *ExpressionContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *ExpressionContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ExpressionContext) T_PLUS() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_PLUS, 0)
}

func (s *ExpressionContext) T_MINUS() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_MINUS, 0)
}

func (s *ExpressionContext) T_TIMES() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_TIMES, 0)
}

func (s *ExpressionContext) T_DIV() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_DIV, 0)
}

func (s *ExpressionContext) T_MOD() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_MOD, 0)
}

func (s *ExpressionContext) T_RANGE() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_RANGE, 0)
}

func (s *ExpressionContext) T_LT() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_LT, 0)
}

func (s *ExpressionContext) T_GT() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_GT, 0)
}

func (s *ExpressionContext) T_LE() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_LE, 0)
}

func (s *ExpressionContext) T_GE() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_GE, 0)
}

func (s *ExpressionContext) T_IN() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_IN, 0)
}

func (s *ExpressionContext) T_NOT() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_NOT, 0)
}

func (s *ExpressionContext) T_EQ() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_EQ, 0)
}

func (s *ExpressionContext) T_NE() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_NE, 0)
}

func (s *ExpressionContext) T_ALL() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_ALL, 0)
}

func (s *ExpressionContext) T_ANY() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_ANY, 0)
}

func (s *ExpressionContext) T_NONE() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_NONE, 0)
}

func (s *ExpressionContext) T_LIKE() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_LIKE, 0)
}

func (s *ExpressionContext) T_REGEX_MATCH() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_REGEX_MATCH, 0)
}

func (s *ExpressionContext) T_REGEX_NON_MATCH() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_REGEX_NON_MATCH, 0)
}

func (s *ExpressionContext) T_AND() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_AND, 0)
}

func (s *ExpressionContext) T_OR() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_OR, 0)
}

func (s *ExpressionContext) T_QUESTION() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_QUESTION, 0)
}

func (s *ExpressionContext) T_COLON() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_COLON, 0)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CAQLParserListener); ok {
		listenerT.EnterExpression(s)
	}
}

func (s *ExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CAQLParserListener); ok {
		listenerT.ExitExpression(s)
	}
}

func (p *CAQLParser) Expression() (localctx IExpressionContext) {
	return p.expression(0)
}

func (p *CAQLParser) expression(_p int) (localctx IExpressionContext) {
	this := p
	_ = this

	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExpressionContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 2
	p.EnterRecursionRule(localctx, 2, CAQLParserRULE_expression, _p)
	var _la int

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(29)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CAQLParserT_FALSE, CAQLParserT_NULL, CAQLParserT_TRUE, CAQLParserT_INT, CAQLParserT_FLOAT, CAQLParserT_QUOTED_STRING:
		{
			p.SetState(26)
			p.Value_literal()
		}

	case CAQLParserT_OPEN, CAQLParserT_OBJECT_OPEN, CAQLParserT_ARRAY_OPEN, CAQLParserT_STRING:
		{
			p.SetState(27)
			p.reference(0)
		}

	case CAQLParserT_PLUS, CAQLParserT_MINUS, CAQLParserT_NOT:
		{
			p.SetState(28)
			p.Operator_unary()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(85)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 4, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(83)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext()) {
			case 1:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, CAQLParserRULE_expression)
				p.SetState(31)

				if !(p.Precpred(p.GetParserRuleContext(), 13)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 13)", ""))
				}
				{
					p.SetState(32)
					_la = p.GetTokenStream().LA(1)

					if !(_la == CAQLParserT_PLUS || _la == CAQLParserT_MINUS) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(33)
					p.expression(14)
				}

			case 2:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, CAQLParserRULE_expression)
				p.SetState(34)

				if !(p.Precpred(p.GetParserRuleContext(), 12)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 12)", ""))
				}
				{
					p.SetState(35)
					_la = p.GetTokenStream().LA(1)

					if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<CAQLParserT_TIMES)|(1<<CAQLParserT_DIV)|(1<<CAQLParserT_MOD))) != 0) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(36)
					p.expression(13)
				}

			case 3:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, CAQLParserRULE_expression)
				p.SetState(37)

				if !(p.Precpred(p.GetParserRuleContext(), 11)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 11)", ""))
				}
				{
					p.SetState(38)
					p.Match(CAQLParserT_RANGE)
				}
				{
					p.SetState(39)
					p.expression(12)
				}

			case 4:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, CAQLParserRULE_expression)
				p.SetState(40)

				if !(p.Precpred(p.GetParserRuleContext(), 10)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 10)", ""))
				}
				{
					p.SetState(41)
					_la = p.GetTokenStream().LA(1)

					if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<CAQLParserT_LT)|(1<<CAQLParserT_GT)|(1<<CAQLParserT_LE)|(1<<CAQLParserT_GE))) != 0) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(42)
					p.expression(11)
				}

			case 5:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, CAQLParserRULE_expression)
				p.SetState(43)

				if !(p.Precpred(p.GetParserRuleContext(), 9)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 9)", ""))
				}
				p.SetState(45)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)

				if _la == CAQLParserT_NOT {
					{
						p.SetState(44)
						p.Match(CAQLParserT_NOT)
					}

				}
				{
					p.SetState(47)
					p.Match(CAQLParserT_IN)
				}
				{
					p.SetState(48)
					p.expression(10)
				}

			case 6:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, CAQLParserRULE_expression)
				p.SetState(49)

				if !(p.Precpred(p.GetParserRuleContext(), 8)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 8)", ""))
				}
				{
					p.SetState(50)
					_la = p.GetTokenStream().LA(1)

					if !(_la == CAQLParserT_EQ || _la == CAQLParserT_NE) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(51)
					p.expression(9)
				}

			case 7:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, CAQLParserRULE_expression)
				p.SetState(52)

				if !(p.Precpred(p.GetParserRuleContext(), 7)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 7)", ""))
				}
				{
					p.SetState(53)
					_la = p.GetTokenStream().LA(1)

					if !(((_la-27)&-(0x1f+1)) == 0 && ((1<<uint((_la-27)))&((1<<(CAQLParserT_ALL-27))|(1<<(CAQLParserT_ANY-27))|(1<<(CAQLParserT_NONE-27)))) != 0) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(54)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*ExpressionContext).eq_op = _lt

					_la = p.GetTokenStream().LA(1)

					if !((((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<CAQLParserT_EQ)|(1<<CAQLParserT_NE)|(1<<CAQLParserT_LT)|(1<<CAQLParserT_GT)|(1<<CAQLParserT_LE)|(1<<CAQLParserT_GE))) != 0) || _la == CAQLParserT_IN) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*ExpressionContext).eq_op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(55)
					p.expression(8)
				}

			case 8:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, CAQLParserRULE_expression)
				p.SetState(56)

				if !(p.Precpred(p.GetParserRuleContext(), 6)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 6)", ""))
				}
				{
					p.SetState(57)
					_la = p.GetTokenStream().LA(1)

					if !(((_la-27)&-(0x1f+1)) == 0 && ((1<<uint((_la-27)))&((1<<(CAQLParserT_ALL-27))|(1<<(CAQLParserT_ANY-27))|(1<<(CAQLParserT_NONE-27)))) != 0) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(58)
					p.Match(CAQLParserT_NOT)
				}
				{
					p.SetState(59)
					p.Match(CAQLParserT_IN)
				}
				{
					p.SetState(60)
					p.expression(7)
				}

			case 9:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, CAQLParserRULE_expression)
				p.SetState(61)

				if !(p.Precpred(p.GetParserRuleContext(), 5)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 5)", ""))
				}
				p.SetState(63)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)

				if _la == CAQLParserT_NOT {
					{
						p.SetState(62)
						p.Match(CAQLParserT_NOT)
					}

				}
				{
					p.SetState(65)
					_la = p.GetTokenStream().LA(1)

					if !(_la == CAQLParserT_REGEX_MATCH || _la == CAQLParserT_REGEX_NON_MATCH || _la == CAQLParserT_LIKE) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(66)
					p.expression(6)
				}

			case 10:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, CAQLParserRULE_expression)
				p.SetState(67)

				if !(p.Precpred(p.GetParserRuleContext(), 4)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 4)", ""))
				}
				{
					p.SetState(68)
					p.Match(CAQLParserT_AND)
				}
				{
					p.SetState(69)
					p.expression(5)
				}

			case 11:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, CAQLParserRULE_expression)
				p.SetState(70)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
				}
				{
					p.SetState(71)
					p.Match(CAQLParserT_OR)
				}
				{
					p.SetState(72)
					p.expression(4)
				}

			case 12:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, CAQLParserRULE_expression)
				p.SetState(73)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
				}
				{
					p.SetState(74)
					p.Match(CAQLParserT_QUESTION)
				}
				{
					p.SetState(75)
					p.expression(0)
				}
				{
					p.SetState(76)
					p.Match(CAQLParserT_COLON)
				}
				{
					p.SetState(77)
					p.expression(3)
				}

			case 13:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, CAQLParserRULE_expression)
				p.SetState(79)

				if !(p.Precpred(p.GetParserRuleContext(), 1)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
				}
				{
					p.SetState(80)
					p.Match(CAQLParserT_QUESTION)
				}
				{
					p.SetState(81)
					p.Match(CAQLParserT_COLON)
				}
				{
					p.SetState(82)
					p.expression(2)
				}

			}

		}
		p.SetState(87)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 4, p.GetParserRuleContext())
	}

	return localctx
}

// IOperator_unaryContext is an interface to support dynamic dispatch.
type IOperator_unaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsOperator_unaryContext differentiates from other interfaces.
	IsOperator_unaryContext()
}

type Operator_unaryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOperator_unaryContext() *Operator_unaryContext {
	var p = new(Operator_unaryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CAQLParserRULE_operator_unary
	return p
}

func (*Operator_unaryContext) IsOperator_unaryContext() {}

func NewOperator_unaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Operator_unaryContext {
	var p = new(Operator_unaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CAQLParserRULE_operator_unary

	return p
}

func (s *Operator_unaryContext) GetParser() antlr.Parser { return s.parser }

func (s *Operator_unaryContext) T_PLUS() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_PLUS, 0)
}

func (s *Operator_unaryContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *Operator_unaryContext) T_MINUS() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_MINUS, 0)
}

func (s *Operator_unaryContext) T_NOT() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_NOT, 0)
}

func (s *Operator_unaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Operator_unaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Operator_unaryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CAQLParserListener); ok {
		listenerT.EnterOperator_unary(s)
	}
}

func (s *Operator_unaryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CAQLParserListener); ok {
		listenerT.ExitOperator_unary(s)
	}
}

func (p *CAQLParser) Operator_unary() (localctx IOperator_unaryContext) {
	this := p
	_ = this

	localctx = NewOperator_unaryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, CAQLParserRULE_operator_unary)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(94)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CAQLParserT_PLUS:
		{
			p.SetState(88)
			p.Match(CAQLParserT_PLUS)
		}
		{
			p.SetState(89)
			p.expression(0)
		}

	case CAQLParserT_MINUS:
		{
			p.SetState(90)
			p.Match(CAQLParserT_MINUS)
		}
		{
			p.SetState(91)
			p.expression(0)
		}

	case CAQLParserT_NOT:
		{
			p.SetState(92)
			p.Match(CAQLParserT_NOT)
		}
		{
			p.SetState(93)
			p.expression(0)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IReferenceContext is an interface to support dynamic dispatch.
type IReferenceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsReferenceContext differentiates from other interfaces.
	IsReferenceContext()
}

type ReferenceContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyReferenceContext() *ReferenceContext {
	var p = new(ReferenceContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CAQLParserRULE_reference
	return p
}

func (*ReferenceContext) IsReferenceContext() {}

func NewReferenceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ReferenceContext {
	var p = new(ReferenceContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CAQLParserRULE_reference

	return p
}

func (s *ReferenceContext) GetParser() antlr.Parser { return s.parser }

func (s *ReferenceContext) T_STRING() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_STRING, 0)
}

func (s *ReferenceContext) Compound_value() ICompound_valueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICompound_valueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICompound_valueContext)
}

func (s *ReferenceContext) Function_call() IFunction_callContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunction_callContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFunction_callContext)
}

func (s *ReferenceContext) T_OPEN() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_OPEN, 0)
}

func (s *ReferenceContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ReferenceContext) T_CLOSE() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_CLOSE, 0)
}

func (s *ReferenceContext) Reference() IReferenceContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IReferenceContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IReferenceContext)
}

func (s *ReferenceContext) DOT() antlr.TerminalNode {
	return s.GetToken(CAQLParserDOT, 0)
}

func (s *ReferenceContext) T_ARRAY_OPEN() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_ARRAY_OPEN, 0)
}

func (s *ReferenceContext) T_ARRAY_CLOSE() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_ARRAY_CLOSE, 0)
}

func (s *ReferenceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ReferenceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ReferenceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CAQLParserListener); ok {
		listenerT.EnterReference(s)
	}
}

func (s *ReferenceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CAQLParserListener); ok {
		listenerT.ExitReference(s)
	}
}

func (p *CAQLParser) Reference() (localctx IReferenceContext) {
	return p.reference(0)
}

func (p *CAQLParser) reference(_p int) (localctx IReferenceContext) {
	this := p
	_ = this

	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewReferenceContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IReferenceContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 6
	p.EnterRecursionRule(localctx, 6, CAQLParserRULE_reference, _p)

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(104)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 6, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(97)
			p.Match(CAQLParserT_STRING)
		}

	case 2:
		{
			p.SetState(98)
			p.Compound_value()
		}

	case 3:
		{
			p.SetState(99)
			p.Function_call()
		}

	case 4:
		{
			p.SetState(100)
			p.Match(CAQLParserT_OPEN)
		}
		{
			p.SetState(101)
			p.expression(0)
		}
		{
			p.SetState(102)
			p.Match(CAQLParserT_CLOSE)
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(116)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(114)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 7, p.GetParserRuleContext()) {
			case 1:
				localctx = NewReferenceContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, CAQLParserRULE_reference)
				p.SetState(106)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
				}
				{
					p.SetState(107)
					p.Match(CAQLParserDOT)
				}
				{
					p.SetState(108)
					p.Match(CAQLParserT_STRING)
				}

			case 2:
				localctx = NewReferenceContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, CAQLParserRULE_reference)
				p.SetState(109)

				if !(p.Precpred(p.GetParserRuleContext(), 1)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
				}
				{
					p.SetState(110)
					p.Match(CAQLParserT_ARRAY_OPEN)
				}
				{
					p.SetState(111)
					p.expression(0)
				}
				{
					p.SetState(112)
					p.Match(CAQLParserT_ARRAY_CLOSE)
				}

			}

		}
		p.SetState(118)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext())
	}

	return localctx
}

// ICompound_valueContext is an interface to support dynamic dispatch.
type ICompound_valueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCompound_valueContext differentiates from other interfaces.
	IsCompound_valueContext()
}

type Compound_valueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCompound_valueContext() *Compound_valueContext {
	var p = new(Compound_valueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CAQLParserRULE_compound_value
	return p
}

func (*Compound_valueContext) IsCompound_valueContext() {}

func NewCompound_valueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Compound_valueContext {
	var p = new(Compound_valueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CAQLParserRULE_compound_value

	return p
}

func (s *Compound_valueContext) GetParser() antlr.Parser { return s.parser }

func (s *Compound_valueContext) Array() IArrayContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArrayContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArrayContext)
}

func (s *Compound_valueContext) Object() IObjectContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IObjectContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IObjectContext)
}

func (s *Compound_valueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Compound_valueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Compound_valueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CAQLParserListener); ok {
		listenerT.EnterCompound_value(s)
	}
}

func (s *Compound_valueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CAQLParserListener); ok {
		listenerT.ExitCompound_value(s)
	}
}

func (p *CAQLParser) Compound_value() (localctx ICompound_valueContext) {
	this := p
	_ = this

	localctx = NewCompound_valueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, CAQLParserRULE_compound_value)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(121)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CAQLParserT_ARRAY_OPEN:
		{
			p.SetState(119)
			p.Array()
		}

	case CAQLParserT_OBJECT_OPEN:
		{
			p.SetState(120)
			p.Object()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IFunction_callContext is an interface to support dynamic dispatch.
type IFunction_callContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFunction_callContext differentiates from other interfaces.
	IsFunction_callContext()
}

type Function_callContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunction_callContext() *Function_callContext {
	var p = new(Function_callContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CAQLParserRULE_function_call
	return p
}

func (*Function_callContext) IsFunction_callContext() {}

func NewFunction_callContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Function_callContext {
	var p = new(Function_callContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CAQLParserRULE_function_call

	return p
}

func (s *Function_callContext) GetParser() antlr.Parser { return s.parser }

func (s *Function_callContext) T_STRING() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_STRING, 0)
}

func (s *Function_callContext) T_OPEN() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_OPEN, 0)
}

func (s *Function_callContext) T_CLOSE() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_CLOSE, 0)
}

func (s *Function_callContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *Function_callContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *Function_callContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(CAQLParserT_COMMA)
}

func (s *Function_callContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(CAQLParserT_COMMA, i)
}

func (s *Function_callContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Function_callContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Function_callContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CAQLParserListener); ok {
		listenerT.EnterFunction_call(s)
	}
}

func (s *Function_callContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CAQLParserListener); ok {
		listenerT.ExitFunction_call(s)
	}
}

func (p *CAQLParser) Function_call() (localctx IFunction_callContext) {
	this := p
	_ = this

	localctx = NewFunction_callContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, CAQLParserRULE_function_call)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(123)
		p.Match(CAQLParserT_STRING)
	}
	{
		p.SetState(124)
		p.Match(CAQLParserT_OPEN)
	}
	p.SetState(126)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if (((_la-10)&-(0x1f+1)) == 0 && ((1<<uint((_la-10)))&((1<<(CAQLParserT_PLUS-10))|(1<<(CAQLParserT_MINUS-10))|(1<<(CAQLParserT_OPEN-10))|(1<<(CAQLParserT_OBJECT_OPEN-10))|(1<<(CAQLParserT_ARRAY_OPEN-10))|(1<<(CAQLParserT_FALSE-10)))) != 0) || (((_la-47)&-(0x1f+1)) == 0 && ((1<<uint((_la-47)))&((1<<(CAQLParserT_NOT-47))|(1<<(CAQLParserT_NULL-47))|(1<<(CAQLParserT_TRUE-47))|(1<<(CAQLParserT_STRING-47))|(1<<(CAQLParserT_INT-47))|(1<<(CAQLParserT_FLOAT-47))|(1<<(CAQLParserT_QUOTED_STRING-47)))) != 0) {
		{
			p.SetState(125)
			p.expression(0)
		}

	}
	p.SetState(132)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 11, p.GetParserRuleContext())

	for _alt != 1 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1+1 {
			{
				p.SetState(128)
				p.Match(CAQLParserT_COMMA)
			}
			{
				p.SetState(129)
				p.expression(0)
			}

		}
		p.SetState(134)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 11, p.GetParserRuleContext())
	}
	p.SetState(136)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == CAQLParserT_COMMA {
		{
			p.SetState(135)
			p.Match(CAQLParserT_COMMA)
		}

	}
	{
		p.SetState(138)
		p.Match(CAQLParserT_CLOSE)
	}

	return localctx
}

// IValue_literalContext is an interface to support dynamic dispatch.
type IValue_literalContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsValue_literalContext differentiates from other interfaces.
	IsValue_literalContext()
}

type Value_literalContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyValue_literalContext() *Value_literalContext {
	var p = new(Value_literalContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CAQLParserRULE_value_literal
	return p
}

func (*Value_literalContext) IsValue_literalContext() {}

func NewValue_literalContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Value_literalContext {
	var p = new(Value_literalContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CAQLParserRULE_value_literal

	return p
}

func (s *Value_literalContext) GetParser() antlr.Parser { return s.parser }

func (s *Value_literalContext) T_QUOTED_STRING() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_QUOTED_STRING, 0)
}

func (s *Value_literalContext) T_INT() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_INT, 0)
}

func (s *Value_literalContext) T_FLOAT() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_FLOAT, 0)
}

func (s *Value_literalContext) T_NULL() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_NULL, 0)
}

func (s *Value_literalContext) T_TRUE() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_TRUE, 0)
}

func (s *Value_literalContext) T_FALSE() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_FALSE, 0)
}

func (s *Value_literalContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Value_literalContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Value_literalContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CAQLParserListener); ok {
		listenerT.EnterValue_literal(s)
	}
}

func (s *Value_literalContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CAQLParserListener); ok {
		listenerT.ExitValue_literal(s)
	}
}

func (p *CAQLParser) Value_literal() (localctx IValue_literalContext) {
	this := p
	_ = this

	localctx = NewValue_literalContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, CAQLParserRULE_value_literal)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(140)
		_la = p.GetTokenStream().LA(1)

		if !((((_la-34)&-(0x1f+1)) == 0 && ((1<<uint((_la-34)))&((1<<(CAQLParserT_FALSE-34))|(1<<(CAQLParserT_NULL-34))|(1<<(CAQLParserT_TRUE-34)))) != 0) || (((_la-70)&-(0x1f+1)) == 0 && ((1<<uint((_la-70)))&((1<<(CAQLParserT_INT-70))|(1<<(CAQLParserT_FLOAT-70))|(1<<(CAQLParserT_QUOTED_STRING-70)))) != 0)) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IArrayContext is an interface to support dynamic dispatch.
type IArrayContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsArrayContext differentiates from other interfaces.
	IsArrayContext()
}

type ArrayContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArrayContext() *ArrayContext {
	var p = new(ArrayContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CAQLParserRULE_array
	return p
}

func (*ArrayContext) IsArrayContext() {}

func NewArrayContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrayContext {
	var p = new(ArrayContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CAQLParserRULE_array

	return p
}

func (s *ArrayContext) GetParser() antlr.Parser { return s.parser }

func (s *ArrayContext) T_ARRAY_OPEN() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_ARRAY_OPEN, 0)
}

func (s *ArrayContext) T_ARRAY_CLOSE() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_ARRAY_CLOSE, 0)
}

func (s *ArrayContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *ArrayContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ArrayContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(CAQLParserT_COMMA)
}

func (s *ArrayContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(CAQLParserT_COMMA, i)
}

func (s *ArrayContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrayContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArrayContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CAQLParserListener); ok {
		listenerT.EnterArray(s)
	}
}

func (s *ArrayContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CAQLParserListener); ok {
		listenerT.ExitArray(s)
	}
}

func (p *CAQLParser) Array() (localctx IArrayContext) {
	this := p
	_ = this

	localctx = NewArrayContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, CAQLParserRULE_array)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(142)
		p.Match(CAQLParserT_ARRAY_OPEN)
	}
	p.SetState(144)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if (((_la-10)&-(0x1f+1)) == 0 && ((1<<uint((_la-10)))&((1<<(CAQLParserT_PLUS-10))|(1<<(CAQLParserT_MINUS-10))|(1<<(CAQLParserT_OPEN-10))|(1<<(CAQLParserT_OBJECT_OPEN-10))|(1<<(CAQLParserT_ARRAY_OPEN-10))|(1<<(CAQLParserT_FALSE-10)))) != 0) || (((_la-47)&-(0x1f+1)) == 0 && ((1<<uint((_la-47)))&((1<<(CAQLParserT_NOT-47))|(1<<(CAQLParserT_NULL-47))|(1<<(CAQLParserT_TRUE-47))|(1<<(CAQLParserT_STRING-47))|(1<<(CAQLParserT_INT-47))|(1<<(CAQLParserT_FLOAT-47))|(1<<(CAQLParserT_QUOTED_STRING-47)))) != 0) {
		{
			p.SetState(143)
			p.expression(0)
		}

	}
	p.SetState(150)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 14, p.GetParserRuleContext())

	for _alt != 1 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1+1 {
			{
				p.SetState(146)
				p.Match(CAQLParserT_COMMA)
			}
			{
				p.SetState(147)
				p.expression(0)
			}

		}
		p.SetState(152)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 14, p.GetParserRuleContext())
	}
	p.SetState(154)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == CAQLParserT_COMMA {
		{
			p.SetState(153)
			p.Match(CAQLParserT_COMMA)
		}

	}
	{
		p.SetState(156)
		p.Match(CAQLParserT_ARRAY_CLOSE)
	}

	return localctx
}

// IObjectContext is an interface to support dynamic dispatch.
type IObjectContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsObjectContext differentiates from other interfaces.
	IsObjectContext()
}

type ObjectContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyObjectContext() *ObjectContext {
	var p = new(ObjectContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CAQLParserRULE_object
	return p
}

func (*ObjectContext) IsObjectContext() {}

func NewObjectContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ObjectContext {
	var p = new(ObjectContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CAQLParserRULE_object

	return p
}

func (s *ObjectContext) GetParser() antlr.Parser { return s.parser }

func (s *ObjectContext) T_OBJECT_OPEN() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_OBJECT_OPEN, 0)
}

func (s *ObjectContext) T_OBJECT_CLOSE() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_OBJECT_CLOSE, 0)
}

func (s *ObjectContext) AllObject_element() []IObject_elementContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IObject_elementContext)(nil)).Elem())
	var tst = make([]IObject_elementContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IObject_elementContext)
		}
	}

	return tst
}

func (s *ObjectContext) Object_element(i int) IObject_elementContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IObject_elementContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IObject_elementContext)
}

func (s *ObjectContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(CAQLParserT_COMMA)
}

func (s *ObjectContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(CAQLParserT_COMMA, i)
}

func (s *ObjectContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ObjectContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ObjectContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CAQLParserListener); ok {
		listenerT.EnterObject(s)
	}
}

func (s *ObjectContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CAQLParserListener); ok {
		listenerT.ExitObject(s)
	}
}

func (p *CAQLParser) Object() (localctx IObjectContext) {
	this := p
	_ = this

	localctx = NewObjectContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, CAQLParserRULE_object)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(158)
		p.Match(CAQLParserT_OBJECT_OPEN)
	}
	p.SetState(160)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == CAQLParserT_ARRAY_OPEN || _la == CAQLParserT_STRING || _la == CAQLParserT_QUOTED_STRING {
		{
			p.SetState(159)
			p.Object_element()
		}

	}
	p.SetState(166)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 17, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(162)
				p.Match(CAQLParserT_COMMA)
			}
			{
				p.SetState(163)
				p.Object_element()
			}

		}
		p.SetState(168)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 17, p.GetParserRuleContext())
	}
	p.SetState(170)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == CAQLParserT_COMMA {
		{
			p.SetState(169)
			p.Match(CAQLParserT_COMMA)
		}

	}
	{
		p.SetState(172)
		p.Match(CAQLParserT_OBJECT_CLOSE)
	}

	return localctx
}

// IObject_elementContext is an interface to support dynamic dispatch.
type IObject_elementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsObject_elementContext differentiates from other interfaces.
	IsObject_elementContext()
}

type Object_elementContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyObject_elementContext() *Object_elementContext {
	var p = new(Object_elementContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CAQLParserRULE_object_element
	return p
}

func (*Object_elementContext) IsObject_elementContext() {}

func NewObject_elementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Object_elementContext {
	var p = new(Object_elementContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CAQLParserRULE_object_element

	return p
}

func (s *Object_elementContext) GetParser() antlr.Parser { return s.parser }

func (s *Object_elementContext) T_STRING() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_STRING, 0)
}

func (s *Object_elementContext) Object_element_name() IObject_element_nameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IObject_element_nameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IObject_element_nameContext)
}

func (s *Object_elementContext) T_COLON() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_COLON, 0)
}

func (s *Object_elementContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *Object_elementContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *Object_elementContext) T_ARRAY_OPEN() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_ARRAY_OPEN, 0)
}

func (s *Object_elementContext) T_ARRAY_CLOSE() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_ARRAY_CLOSE, 0)
}

func (s *Object_elementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Object_elementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Object_elementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CAQLParserListener); ok {
		listenerT.EnterObject_element(s)
	}
}

func (s *Object_elementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CAQLParserListener); ok {
		listenerT.ExitObject_element(s)
	}
}

func (p *CAQLParser) Object_element() (localctx IObject_elementContext) {
	this := p
	_ = this

	localctx = NewObject_elementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, CAQLParserRULE_object_element)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(185)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 19, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(174)
			p.Match(CAQLParserT_STRING)
		}

	case 2:
		{
			p.SetState(175)
			p.Object_element_name()
		}
		{
			p.SetState(176)
			p.Match(CAQLParserT_COLON)
		}
		{
			p.SetState(177)
			p.expression(0)
		}

	case 3:
		{
			p.SetState(179)
			p.Match(CAQLParserT_ARRAY_OPEN)
		}
		{
			p.SetState(180)
			p.expression(0)
		}
		{
			p.SetState(181)
			p.Match(CAQLParserT_ARRAY_CLOSE)
		}
		{
			p.SetState(182)
			p.Match(CAQLParserT_COLON)
		}
		{
			p.SetState(183)
			p.expression(0)
		}

	}

	return localctx
}

// IObject_element_nameContext is an interface to support dynamic dispatch.
type IObject_element_nameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsObject_element_nameContext differentiates from other interfaces.
	IsObject_element_nameContext()
}

type Object_element_nameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyObject_element_nameContext() *Object_element_nameContext {
	var p = new(Object_element_nameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CAQLParserRULE_object_element_name
	return p
}

func (*Object_element_nameContext) IsObject_element_nameContext() {}

func NewObject_element_nameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Object_element_nameContext {
	var p = new(Object_element_nameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CAQLParserRULE_object_element_name

	return p
}

func (s *Object_element_nameContext) GetParser() antlr.Parser { return s.parser }

func (s *Object_element_nameContext) T_STRING() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_STRING, 0)
}

func (s *Object_element_nameContext) T_QUOTED_STRING() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_QUOTED_STRING, 0)
}

func (s *Object_element_nameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Object_element_nameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Object_element_nameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CAQLParserListener); ok {
		listenerT.EnterObject_element_name(s)
	}
}

func (s *Object_element_nameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CAQLParserListener); ok {
		listenerT.ExitObject_element_name(s)
	}
}

func (p *CAQLParser) Object_element_name() (localctx IObject_element_nameContext) {
	this := p
	_ = this

	localctx = NewObject_element_nameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, CAQLParserRULE_object_element_name)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(187)
		_la = p.GetTokenStream().LA(1)

		if !(_la == CAQLParserT_STRING || _la == CAQLParserT_QUOTED_STRING) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

func (p *CAQLParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 1:
		var t *ExpressionContext = nil
		if localctx != nil {
			t = localctx.(*ExpressionContext)
		}
		return p.Expression_Sempred(t, predIndex)

	case 3:
		var t *ReferenceContext = nil
		if localctx != nil {
			t = localctx.(*ReferenceContext)
		}
		return p.Reference_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *CAQLParser) Expression_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	this := p
	_ = this

	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 13)

	case 1:
		return p.Precpred(p.GetParserRuleContext(), 12)

	case 2:
		return p.Precpred(p.GetParserRuleContext(), 11)

	case 3:
		return p.Precpred(p.GetParserRuleContext(), 10)

	case 4:
		return p.Precpred(p.GetParserRuleContext(), 9)

	case 5:
		return p.Precpred(p.GetParserRuleContext(), 8)

	case 6:
		return p.Precpred(p.GetParserRuleContext(), 7)

	case 7:
		return p.Precpred(p.GetParserRuleContext(), 6)

	case 8:
		return p.Precpred(p.GetParserRuleContext(), 5)

	case 9:
		return p.Precpred(p.GetParserRuleContext(), 4)

	case 10:
		return p.Precpred(p.GetParserRuleContext(), 3)

	case 11:
		return p.Precpred(p.GetParserRuleContext(), 2)

	case 12:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *CAQLParser) Reference_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	this := p
	_ = this

	switch predIndex {
	case 13:
		return p.Precpred(p.GetParserRuleContext(), 2)

	case 14:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
