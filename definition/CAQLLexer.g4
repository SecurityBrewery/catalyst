/*
 * The MIT License (MIT)
 *
 * Copyright (c) 2020 by Martin Mirchev
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and
 * associated documentation files (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge, publish, distribute,
 * sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or
 * substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT
 * NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 * Project : sqlite-parser; an ANTLR4 grammar for SQLite https://github.com/bkiers/sqlite-parser
 * Developed by : Bart Kiers, bart@big-o.nl
 */

// $antlr-format alignTrailingComments on, columnLimit 150, maxEmptyLinesToKeep 1, reflowComments off, useTab off
// $antlr-format allowShortRulesOnASingleLine on, alignSemicolons ownLine

lexer grammar CAQLLexer;

channels { ERRORCHANNEL }


DOT:       '.';

// https://github.com/arangodb/arangodb/blob/devel/arangod/Aql/grammar.y
T_REGEX_MATCH:     '=~'; // "~= operator"
T_REGEX_NON_MATCH: '!~'; // "~! operator"

T_EQ: '=='; // "== operator";
T_NE: '!='; // "!= operator";
T_LT: '<';  // "< operator";
T_GT: '>';  // "> operator";
T_LE: '<='; // "<= operator";
T_GE: '>='; // ">= operator";

T_PLUS:  '+'; // "+ operator"
T_MINUS: '-'; // "- operator"
T_TIMES: '*'; // "* operator"
T_DIV:   '/'; // "/ operator"
T_MOD:   '%'; // "% operator"

T_QUESTION: '?';  // "?"
T_COLON:    ':';  // ":"
T_SCOPE:    '::'; // "::"
T_RANGE:    '..'; // ".."

T_COMMA:        ','; // ","
T_OPEN:         '('; // "("
T_CLOSE:        ')'; // ")"
T_OBJECT_OPEN:  '{'; // "{"
T_OBJECT_CLOSE: '}'; // "}"
T_ARRAY_OPEN:   '['; // "["
T_ARRAY_CLOSE:  ']'; // "]"


// https://www.arangodb.com/docs/stable/aql/fundamentals-syntax.html#keywords
T_AGGREGATE:         A G G R E G A T E;
T_ALL:               A L L;
T_AND:               (A N D | '&&');
T_ANY:               A N Y;
T_ASC:               A S C;
T_COLLECT:           C O L L E C T;
T_DESC:              D E S C;
T_DISTINCT:          D I S T I N C T;
T_FALSE:             F A L S E;
T_FILTER:            F I L T E R;
T_FOR:               F O R;
T_GRAPH:             G R A P H;
T_IN:                I N;
T_INBOUND:           I N B O U N D;
T_INSERT:            I N S E R T;
T_INTO:              I N T O;
T_K_SHORTEST_PATHS:  K '_' S H O R T E S T '_' P A T H S;
T_LET:               L E T;
T_LIKE:              L I K E;
T_LIMIT:             L I M I T;
T_NONE:              N O N E;
T_NOT:               (N O T | '!');
T_NULL:              N U L L;
T_OR:                (O R | '||');
T_OUTBOUND:          O U T B O U N D;
T_REMOVE:            R E M O V E;
T_REPLACE:           R E P L A C E;
T_RETURN:            R E T U R N;
T_SHORTEST_PATH:     S H O R T E S T '_' P A T H;
T_SORT:              S O R T;
T_TRUE:              T R U E;
T_UPDATE:            U P D A T E;
T_UPSERT:            U P S E R T;
T_WITH:              W I T H;

T_KEEP:              K E E P;
T_COUNT:             C O U N T;
T_OPTIONS:           O P T I O N S;
T_PRUNE:             P R U N E;
T_SEARCH:            S E A R C H;
T_TO:                T O;

T_CURRENT:           C U R R E N T;
T_NEW:               N E W;
T_OLD:               O L D;

T_STRING: [a-zA-Z_] [a-zA-Z_0-9]*;

T_INT: [1-9] DIGIT* | '0' | '0x' HEX_DIGIT+ | '0b' [0-1]+;
T_FLOAT: ( [1-9] DIGIT* | '0' )? '.' DIGIT+ (E [-+]? DIGIT+)?;

T_PARAMETER: '@' T_STRING;

T_QUOTED_STRING: ('\'' ('\\'. | '\'\'' | ~('\'' | '\\'))* '\'' | '"' ( '\\'. | '""' | ~('"'| '\\') )* '"');

SINGLE_LINE_COMMENT: '//' ~[\r\n]* (('\r'? '\n') | EOF) -> channel(HIDDEN);

MULTILINE_COMMENT: '/*' .*? '*/' -> channel(HIDDEN);

SPACES: [ \u000B\t\r\n] -> channel(HIDDEN);

UNEXPECTED_CHAR: .;

fragment HEX_DIGIT: [0-9a-fA-F];
fragment DIGIT:     [0-9];

fragment A: [aA];
fragment B: [bB];
fragment C: [cC];
fragment D: [dD];
fragment E: [eE];
fragment F: [fF];
fragment G: [gG];
fragment H: [hH];
fragment I: [iI];
fragment J: [jJ];
fragment K: [kK];
fragment L: [lL];
fragment M: [mM];
fragment N: [nN];
fragment O: [oO];
fragment P: [pP];
fragment Q: [qQ];
fragment R: [rR];
fragment S: [sS];
fragment T: [tT];
fragment U: [uU];
fragment V: [vV];
fragment W: [wW];
fragment X: [xX];
fragment Y: [yY];
fragment Z: [zZ];

ERROR_RECONGNIGION:                  .    -> channel(ERRORCHANNEL);
