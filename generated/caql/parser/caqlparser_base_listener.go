// Code generated from CAQLParser.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // CAQLParser

import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseCAQLParserListener is a complete listener for a parse tree produced by CAQLParser.
type BaseCAQLParserListener struct{}

var _ CAQLParserListener = &BaseCAQLParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseCAQLParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseCAQLParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseCAQLParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseCAQLParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterParse is called when production parse is entered.
func (s *BaseCAQLParserListener) EnterParse(ctx *ParseContext) {}

// ExitParse is called when production parse is exited.
func (s *BaseCAQLParserListener) ExitParse(ctx *ParseContext) {}

// EnterExpression is called when production expression is entered.
func (s *BaseCAQLParserListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BaseCAQLParserListener) ExitExpression(ctx *ExpressionContext) {}

// EnterOperator_unary is called when production operator_unary is entered.
func (s *BaseCAQLParserListener) EnterOperator_unary(ctx *Operator_unaryContext) {}

// ExitOperator_unary is called when production operator_unary is exited.
func (s *BaseCAQLParserListener) ExitOperator_unary(ctx *Operator_unaryContext) {}

// EnterReference is called when production reference is entered.
func (s *BaseCAQLParserListener) EnterReference(ctx *ReferenceContext) {}

// ExitReference is called when production reference is exited.
func (s *BaseCAQLParserListener) ExitReference(ctx *ReferenceContext) {}

// EnterCompound_value is called when production compound_value is entered.
func (s *BaseCAQLParserListener) EnterCompound_value(ctx *Compound_valueContext) {}

// ExitCompound_value is called when production compound_value is exited.
func (s *BaseCAQLParserListener) ExitCompound_value(ctx *Compound_valueContext) {}

// EnterFunction_call is called when production function_call is entered.
func (s *BaseCAQLParserListener) EnterFunction_call(ctx *Function_callContext) {}

// ExitFunction_call is called when production function_call is exited.
func (s *BaseCAQLParserListener) ExitFunction_call(ctx *Function_callContext) {}

// EnterValue_literal is called when production value_literal is entered.
func (s *BaseCAQLParserListener) EnterValue_literal(ctx *Value_literalContext) {}

// ExitValue_literal is called when production value_literal is exited.
func (s *BaseCAQLParserListener) ExitValue_literal(ctx *Value_literalContext) {}

// EnterArray is called when production array is entered.
func (s *BaseCAQLParserListener) EnterArray(ctx *ArrayContext) {}

// ExitArray is called when production array is exited.
func (s *BaseCAQLParserListener) ExitArray(ctx *ArrayContext) {}

// EnterObject is called when production object is entered.
func (s *BaseCAQLParserListener) EnterObject(ctx *ObjectContext) {}

// ExitObject is called when production object is exited.
func (s *BaseCAQLParserListener) ExitObject(ctx *ObjectContext) {}

// EnterObject_element is called when production object_element is entered.
func (s *BaseCAQLParserListener) EnterObject_element(ctx *Object_elementContext) {}

// ExitObject_element is called when production object_element is exited.
func (s *BaseCAQLParserListener) ExitObject_element(ctx *Object_elementContext) {}

// EnterObject_element_name is called when production object_element_name is entered.
func (s *BaseCAQLParserListener) EnterObject_element_name(ctx *Object_element_nameContext) {}

// ExitObject_element_name is called when production object_element_name is exited.
func (s *BaseCAQLParserListener) ExitObject_element_name(ctx *Object_element_nameContext) {}
