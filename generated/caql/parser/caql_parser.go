// Code generated from CAQLParser.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // CAQLParser

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type CAQLParser struct {
	*antlr.BaseParser
}

var caqlparserParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func caqlparserParserInit() {
	staticData := &caqlparserParserStaticData
	staticData.literalNames = []string{
		"", "'.'", "'=~'", "'!~'", "'=='", "'!='", "'<'", "'>'", "'<='", "'>='",
		"'+'", "'-'", "'*'", "'/'", "'%'", "'?'", "':'", "'::'", "'..'", "','",
		"'('", "')'", "'{'", "'}'", "'['", "']'",
	}
	staticData.symbolicNames = []string{
		"", "DOT", "T_REGEX_MATCH", "T_REGEX_NON_MATCH", "T_EQ", "T_NE", "T_LT",
		"T_GT", "T_LE", "T_GE", "T_PLUS", "T_MINUS", "T_TIMES", "T_DIV", "T_MOD",
		"T_QUESTION", "T_COLON", "T_SCOPE", "T_RANGE", "T_COMMA", "T_OPEN",
		"T_CLOSE", "T_OBJECT_OPEN", "T_OBJECT_CLOSE", "T_ARRAY_OPEN", "T_ARRAY_CLOSE",
		"T_AGGREGATE", "T_ALL", "T_AND", "T_ANY", "T_ASC", "T_COLLECT", "T_DESC",
		"T_DISTINCT", "T_FALSE", "T_FILTER", "T_FOR", "T_GRAPH", "T_IN", "T_INBOUND",
		"T_INSERT", "T_INTO", "T_K_SHORTEST_PATHS", "T_LET", "T_LIKE", "T_LIMIT",
		"T_NONE", "T_NOT", "T_NULL", "T_OR", "T_OUTBOUND", "T_REMOVE", "T_REPLACE",
		"T_RETURN", "T_SHORTEST_PATH", "T_SORT", "T_TRUE", "T_UPDATE", "T_UPSERT",
		"T_WITH", "T_KEEP", "T_COUNT", "T_OPTIONS", "T_PRUNE", "T_SEARCH", "T_TO",
		"T_CURRENT", "T_NEW", "T_OLD", "T_STRING", "T_INT", "T_FLOAT", "T_PARAMETER",
		"T_QUOTED_STRING", "SINGLE_LINE_COMMENT", "MULTILINE_COMMENT", "SPACES",
		"UNEXPECTED_CHAR", "ERROR_RECONGNIGION",
	}
	staticData.ruleNames = []string{
		"parse", "expression", "operator_unary", "reference", "compound_value",
		"function_call", "value_literal", "array", "object", "object_element",
		"object_element_name",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 78, 190, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 30, 8, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 3, 1, 46, 8, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 64, 8, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 5, 1, 84, 8, 1, 10, 1, 12, 1, 87, 9, 1, 1, 2, 1, 2,
		1, 2, 1, 2, 1, 2, 1, 2, 3, 2, 95, 8, 2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1,
		3, 1, 3, 1, 3, 3, 3, 105, 8, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1,
		3, 1, 3, 5, 3, 115, 8, 3, 10, 3, 12, 3, 118, 9, 3, 1, 4, 1, 4, 3, 4, 122,
		8, 4, 1, 5, 1, 5, 1, 5, 3, 5, 127, 8, 5, 1, 5, 1, 5, 5, 5, 131, 8, 5, 10,
		5, 12, 5, 134, 9, 5, 1, 5, 3, 5, 137, 8, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1,
		7, 1, 7, 3, 7, 145, 8, 7, 1, 7, 1, 7, 5, 7, 149, 8, 7, 10, 7, 12, 7, 152,
		9, 7, 1, 7, 3, 7, 155, 8, 7, 1, 7, 1, 7, 1, 8, 1, 8, 3, 8, 161, 8, 8, 1,
		8, 1, 8, 5, 8, 165, 8, 8, 10, 8, 12, 8, 168, 9, 8, 1, 8, 3, 8, 171, 8,
		8, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1,
		9, 1, 9, 3, 9, 186, 8, 9, 1, 10, 1, 10, 1, 10, 2, 132, 150, 2, 2, 6, 11,
		0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 0, 9, 1, 0, 10, 11, 1, 0, 12, 14,
		1, 0, 6, 9, 1, 0, 4, 5, 3, 0, 27, 27, 29, 29, 46, 46, 2, 0, 4, 9, 38, 38,
		2, 0, 2, 3, 44, 44, 5, 0, 34, 34, 48, 48, 56, 56, 70, 71, 73, 73, 2, 0,
		69, 69, 73, 73, 214, 0, 22, 1, 0, 0, 0, 2, 29, 1, 0, 0, 0, 4, 94, 1, 0,
		0, 0, 6, 104, 1, 0, 0, 0, 8, 121, 1, 0, 0, 0, 10, 123, 1, 0, 0, 0, 12,
		140, 1, 0, 0, 0, 14, 142, 1, 0, 0, 0, 16, 158, 1, 0, 0, 0, 18, 185, 1,
		0, 0, 0, 20, 187, 1, 0, 0, 0, 22, 23, 3, 2, 1, 0, 23, 24, 5, 0, 0, 1, 24,
		1, 1, 0, 0, 0, 25, 26, 6, 1, -1, 0, 26, 30, 3, 12, 6, 0, 27, 30, 3, 6,
		3, 0, 28, 30, 3, 4, 2, 0, 29, 25, 1, 0, 0, 0, 29, 27, 1, 0, 0, 0, 29, 28,
		1, 0, 0, 0, 30, 85, 1, 0, 0, 0, 31, 32, 10, 13, 0, 0, 32, 33, 7, 0, 0,
		0, 33, 84, 3, 2, 1, 14, 34, 35, 10, 12, 0, 0, 35, 36, 7, 1, 0, 0, 36, 84,
		3, 2, 1, 13, 37, 38, 10, 11, 0, 0, 38, 39, 5, 18, 0, 0, 39, 84, 3, 2, 1,
		12, 40, 41, 10, 10, 0, 0, 41, 42, 7, 2, 0, 0, 42, 84, 3, 2, 1, 11, 43,
		45, 10, 9, 0, 0, 44, 46, 5, 47, 0, 0, 45, 44, 1, 0, 0, 0, 45, 46, 1, 0,
		0, 0, 46, 47, 1, 0, 0, 0, 47, 48, 5, 38, 0, 0, 48, 84, 3, 2, 1, 10, 49,
		50, 10, 8, 0, 0, 50, 51, 7, 3, 0, 0, 51, 84, 3, 2, 1, 9, 52, 53, 10, 7,
		0, 0, 53, 54, 7, 4, 0, 0, 54, 55, 7, 5, 0, 0, 55, 84, 3, 2, 1, 8, 56, 57,
		10, 6, 0, 0, 57, 58, 7, 4, 0, 0, 58, 59, 5, 47, 0, 0, 59, 60, 5, 38, 0,
		0, 60, 84, 3, 2, 1, 7, 61, 63, 10, 5, 0, 0, 62, 64, 5, 47, 0, 0, 63, 62,
		1, 0, 0, 0, 63, 64, 1, 0, 0, 0, 64, 65, 1, 0, 0, 0, 65, 66, 7, 6, 0, 0,
		66, 84, 3, 2, 1, 6, 67, 68, 10, 4, 0, 0, 68, 69, 5, 28, 0, 0, 69, 84, 3,
		2, 1, 5, 70, 71, 10, 3, 0, 0, 71, 72, 5, 49, 0, 0, 72, 84, 3, 2, 1, 4,
		73, 74, 10, 2, 0, 0, 74, 75, 5, 15, 0, 0, 75, 76, 3, 2, 1, 0, 76, 77, 5,
		16, 0, 0, 77, 78, 3, 2, 1, 3, 78, 84, 1, 0, 0, 0, 79, 80, 10, 1, 0, 0,
		80, 81, 5, 15, 0, 0, 81, 82, 5, 16, 0, 0, 82, 84, 3, 2, 1, 2, 83, 31, 1,
		0, 0, 0, 83, 34, 1, 0, 0, 0, 83, 37, 1, 0, 0, 0, 83, 40, 1, 0, 0, 0, 83,
		43, 1, 0, 0, 0, 83, 49, 1, 0, 0, 0, 83, 52, 1, 0, 0, 0, 83, 56, 1, 0, 0,
		0, 83, 61, 1, 0, 0, 0, 83, 67, 1, 0, 0, 0, 83, 70, 1, 0, 0, 0, 83, 73,
		1, 0, 0, 0, 83, 79, 1, 0, 0, 0, 84, 87, 1, 0, 0, 0, 85, 83, 1, 0, 0, 0,
		85, 86, 1, 0, 0, 0, 86, 3, 1, 0, 0, 0, 87, 85, 1, 0, 0, 0, 88, 89, 5, 10,
		0, 0, 89, 95, 3, 2, 1, 0, 90, 91, 5, 11, 0, 0, 91, 95, 3, 2, 1, 0, 92,
		93, 5, 47, 0, 0, 93, 95, 3, 2, 1, 0, 94, 88, 1, 0, 0, 0, 94, 90, 1, 0,
		0, 0, 94, 92, 1, 0, 0, 0, 95, 5, 1, 0, 0, 0, 96, 97, 6, 3, -1, 0, 97, 105,
		5, 69, 0, 0, 98, 105, 3, 8, 4, 0, 99, 105, 3, 10, 5, 0, 100, 101, 5, 20,
		0, 0, 101, 102, 3, 2, 1, 0, 102, 103, 5, 21, 0, 0, 103, 105, 1, 0, 0, 0,
		104, 96, 1, 0, 0, 0, 104, 98, 1, 0, 0, 0, 104, 99, 1, 0, 0, 0, 104, 100,
		1, 0, 0, 0, 105, 116, 1, 0, 0, 0, 106, 107, 10, 2, 0, 0, 107, 108, 5, 1,
		0, 0, 108, 115, 5, 69, 0, 0, 109, 110, 10, 1, 0, 0, 110, 111, 5, 24, 0,
		0, 111, 112, 3, 2, 1, 0, 112, 113, 5, 25, 0, 0, 113, 115, 1, 0, 0, 0, 114,
		106, 1, 0, 0, 0, 114, 109, 1, 0, 0, 0, 115, 118, 1, 0, 0, 0, 116, 114,
		1, 0, 0, 0, 116, 117, 1, 0, 0, 0, 117, 7, 1, 0, 0, 0, 118, 116, 1, 0, 0,
		0, 119, 122, 3, 14, 7, 0, 120, 122, 3, 16, 8, 0, 121, 119, 1, 0, 0, 0,
		121, 120, 1, 0, 0, 0, 122, 9, 1, 0, 0, 0, 123, 124, 5, 69, 0, 0, 124, 126,
		5, 20, 0, 0, 125, 127, 3, 2, 1, 0, 126, 125, 1, 0, 0, 0, 126, 127, 1, 0,
		0, 0, 127, 132, 1, 0, 0, 0, 128, 129, 5, 19, 0, 0, 129, 131, 3, 2, 1, 0,
		130, 128, 1, 0, 0, 0, 131, 134, 1, 0, 0, 0, 132, 133, 1, 0, 0, 0, 132,
		130, 1, 0, 0, 0, 133, 136, 1, 0, 0, 0, 134, 132, 1, 0, 0, 0, 135, 137,
		5, 19, 0, 0, 136, 135, 1, 0, 0, 0, 136, 137, 1, 0, 0, 0, 137, 138, 1, 0,
		0, 0, 138, 139, 5, 21, 0, 0, 139, 11, 1, 0, 0, 0, 140, 141, 7, 7, 0, 0,
		141, 13, 1, 0, 0, 0, 142, 144, 5, 24, 0, 0, 143, 145, 3, 2, 1, 0, 144,
		143, 1, 0, 0, 0, 144, 145, 1, 0, 0, 0, 145, 150, 1, 0, 0, 0, 146, 147,
		5, 19, 0, 0, 147, 149, 3, 2, 1, 0, 148, 146, 1, 0, 0, 0, 149, 152, 1, 0,
		0, 0, 150, 151, 1, 0, 0, 0, 150, 148, 1, 0, 0, 0, 151, 154, 1, 0, 0, 0,
		152, 150, 1, 0, 0, 0, 153, 155, 5, 19, 0, 0, 154, 153, 1, 0, 0, 0, 154,
		155, 1, 0, 0, 0, 155, 156, 1, 0, 0, 0, 156, 157, 5, 25, 0, 0, 157, 15,
		1, 0, 0, 0, 158, 160, 5, 22, 0, 0, 159, 161, 3, 18, 9, 0, 160, 159, 1,
		0, 0, 0, 160, 161, 1, 0, 0, 0, 161, 166, 1, 0, 0, 0, 162, 163, 5, 19, 0,
		0, 163, 165, 3, 18, 9, 0, 164, 162, 1, 0, 0, 0, 165, 168, 1, 0, 0, 0, 166,
		164, 1, 0, 0, 0, 166, 167, 1, 0, 0, 0, 167, 170, 1, 0, 0, 0, 168, 166,
		1, 0, 0, 0, 169, 171, 5, 19, 0, 0, 170, 169, 1, 0, 0, 0, 170, 171, 1, 0,
		0, 0, 171, 172, 1, 0, 0, 0, 172, 173, 5, 23, 0, 0, 173, 17, 1, 0, 0, 0,
		174, 186, 5, 69, 0, 0, 175, 176, 3, 20, 10, 0, 176, 177, 5, 16, 0, 0, 177,
		178, 3, 2, 1, 0, 178, 186, 1, 0, 0, 0, 179, 180, 5, 24, 0, 0, 180, 181,
		3, 2, 1, 0, 181, 182, 5, 25, 0, 0, 182, 183, 5, 16, 0, 0, 183, 184, 3,
		2, 1, 0, 184, 186, 1, 0, 0, 0, 185, 174, 1, 0, 0, 0, 185, 175, 1, 0, 0,
		0, 185, 179, 1, 0, 0, 0, 186, 19, 1, 0, 0, 0, 187, 188, 7, 8, 0, 0, 188,
		21, 1, 0, 0, 0, 20, 29, 45, 63, 83, 85, 94, 104, 114, 116, 121, 126, 132,
		136, 144, 150, 154, 160, 166, 170, 185,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// CAQLParserInit initializes any static state used to implement CAQLParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewCAQLParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func CAQLParserInit() {
	staticData := &caqlparserParserStaticData
	staticData.once.Do(caqlparserParserInit)
}

// NewCAQLParser produces a new parser instance for the optional input antlr.TokenStream.
func NewCAQLParser(input antlr.TokenStream) *CAQLParser {
	CAQLParserInit()
	this := new(CAQLParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &caqlparserParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	this.RuleNames = staticData.ruleNames
	this.LiteralNames = staticData.literalNames
	this.SymbolicNames = staticData.symbolicNames
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
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

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
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValue_literalContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValue_literalContext)
}

func (s *ExpressionContext) Reference() IReferenceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IReferenceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IReferenceContext)
}

func (s *ExpressionContext) Operator_unary() IOperator_unaryContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOperator_unaryContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOperator_unaryContext)
}

func (s *ExpressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *ExpressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

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
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

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
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICompound_valueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICompound_valueContext)
}

func (s *ReferenceContext) Function_call() IFunction_callContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunction_callContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunction_callContext)
}

func (s *ReferenceContext) T_OPEN() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_OPEN, 0)
}

func (s *ReferenceContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ReferenceContext) T_CLOSE() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_CLOSE, 0)
}

func (s *ReferenceContext) Reference() IReferenceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IReferenceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

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
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArrayContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArrayContext)
}

func (s *Compound_valueContext) Object() IObjectContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

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
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *Function_callContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

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
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *ArrayContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

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
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IObject_elementContext); ok {
			len++
		}
	}

	tst := make([]IObject_elementContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IObject_elementContext); ok {
			tst[i] = t.(IObject_elementContext)
			i++
		}
	}

	return tst
}

func (s *ObjectContext) Object_element(i int) IObject_elementContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObject_elementContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

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
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObject_element_nameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObject_element_nameContext)
}

func (s *Object_elementContext) T_COLON() antlr.TerminalNode {
	return s.GetToken(CAQLParserT_COLON, 0)
}

func (s *Object_elementContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *Object_elementContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

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
