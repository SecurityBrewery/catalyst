/*
 * The MIT License (MIT)
 *
 * Copyright (c) 2014 by Bart Kiers
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
 * Developed by:
 *     Bart Kiers, bart@big-o.nl
 *     Martin Mirchev, marti_2203@abv.bg
 *     Mike Lische, mike@lischke-online.de
 */

// $antlr-format alignTrailingComments on, columnLimit 130, minEmptyLines 1, maxEmptyLinesToKeep 1, reflowComments off
// $antlr-format useTab off, allowShortRulesOnASingleLine off, allowShortBlocksOnASingleLine on, alignSemicolons ownLine

parser grammar CAQLParser;

options {
    tokenVocab = CAQLLexer;
}

parse: expression EOF
;

expression:
   value_literal
  | reference
  | operator_unary
  | expression (T_PLUS|T_MINUS) expression
  | expression (T_TIMES|T_DIV|T_MOD) expression
  | expression T_RANGE expression
  | expression (T_LT|T_GT|T_LE|T_GE) expression
  | expression T_NOT? T_IN expression
  | expression (T_EQ|T_NE) expression
  | expression (T_ALL|T_ANY|T_NONE) eq_op=(T_EQ|T_NE|T_LT|T_GT|T_LE|T_GE|T_IN) expression
  | expression (T_ALL|T_ANY|T_NONE) T_NOT T_IN expression
  | expression T_NOT? (T_LIKE|T_REGEX_MATCH|T_REGEX_NON_MATCH) expression
  | expression T_AND expression
  | expression T_OR expression
  | expression T_QUESTION expression T_COLON expression
  | expression T_QUESTION T_COLON expression
;

operator_unary: (
    T_PLUS expression
  | T_MINUS expression
  | T_NOT expression
);

reference:
    T_STRING
  | compound_value
  | function_call
  | T_OPEN expression T_CLOSE
  | reference DOT T_STRING
  | reference T_ARRAY_OPEN expression T_ARRAY_CLOSE
;

compound_value: (
    array
  | object
);

function_call: (
    T_STRING T_OPEN expression? (T_COMMA expression)*? T_COMMA? T_CLOSE
);

value_literal: (
    T_QUOTED_STRING
  | T_INT
  | T_FLOAT
  | T_NULL
  | T_TRUE
  | T_FALSE
);

array:(
    T_ARRAY_OPEN expression? (T_COMMA expression)*? T_COMMA? T_ARRAY_CLOSE
);

object:
    T_OBJECT_OPEN object_element? (T_COMMA object_element)* T_COMMA? T_OBJECT_CLOSE
;

object_element:(
    T_STRING
  | object_element_name T_COLON expression
  | T_ARRAY_OPEN expression T_ARRAY_CLOSE T_COLON expression
);

object_element_name:(
    T_STRING
  | T_QUOTED_STRING
);
