import lexerModule from "./grammar/CAQLLexer.js";
import parserModule from "./grammar/CAQLParser.js";

import antlr4 from "antlr4";

class ErrorListener extends antlr4.error.ErrorListener {
 /**
  * Checks syntax error
  *
  * @param {object} recognizer The parsing support code essentially. Most of it is error recovery stuff
  * @param {object} symbol Offending symbol
  * @param {int} line Line of offending symbol
  * @param {int} column Position in line of offending symbol
  * @param {string} message Error message
  * @param {string} payload Stack trace
  */
 syntaxError(
     recognizer: any,
     symbol: any,
     line: number,
     column: number,
     message: string,
     payload: string
 ) {
  throw {symbol, line, column, message, payload};
 }
}

export function validateCAQL(term: string): any {
 const chars = new antlr4.InputStream(term);
 const lexer = new lexerModule(chars);

 const tokens = new antlr4.CommonTokenStream(lexer);
 const parser = new parserModule(tokens);

 const listener = new ErrorListener();

 parser.removeErrorListeners();
 parser.addErrorListener(listener);
 try {
  parser.parse()
  return null;
 } catch (error) {
  return error;
 }
}
