// Code generated from CAQLParser.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // CAQLParser

import "github.com/antlr/antlr4/runtime/Go/antlr"

// CAQLParserListener is a complete listener for a parse tree produced by CAQLParser.
type CAQLParserListener interface {
	antlr.ParseTreeListener

	// EnterParse is called when entering the parse production.
	EnterParse(c *ParseContext)

	// EnterExpression is called when entering the expression production.
	EnterExpression(c *ExpressionContext)

	// EnterOperator_unary is called when entering the operator_unary production.
	EnterOperator_unary(c *Operator_unaryContext)

	// EnterReference is called when entering the reference production.
	EnterReference(c *ReferenceContext)

	// EnterCompound_value is called when entering the compound_value production.
	EnterCompound_value(c *Compound_valueContext)

	// EnterFunction_call is called when entering the function_call production.
	EnterFunction_call(c *Function_callContext)

	// EnterValue_literal is called when entering the value_literal production.
	EnterValue_literal(c *Value_literalContext)

	// EnterArray is called when entering the array production.
	EnterArray(c *ArrayContext)

	// EnterObject is called when entering the object production.
	EnterObject(c *ObjectContext)

	// EnterObject_element is called when entering the object_element production.
	EnterObject_element(c *Object_elementContext)

	// EnterObject_element_name is called when entering the object_element_name production.
	EnterObject_element_name(c *Object_element_nameContext)

	// ExitParse is called when exiting the parse production.
	ExitParse(c *ParseContext)

	// ExitExpression is called when exiting the expression production.
	ExitExpression(c *ExpressionContext)

	// ExitOperator_unary is called when exiting the operator_unary production.
	ExitOperator_unary(c *Operator_unaryContext)

	// ExitReference is called when exiting the reference production.
	ExitReference(c *ReferenceContext)

	// ExitCompound_value is called when exiting the compound_value production.
	ExitCompound_value(c *Compound_valueContext)

	// ExitFunction_call is called when exiting the function_call production.
	ExitFunction_call(c *Function_callContext)

	// ExitValue_literal is called when exiting the value_literal production.
	ExitValue_literal(c *Value_literalContext)

	// ExitArray is called when exiting the array production.
	ExitArray(c *ArrayContext)

	// ExitObject is called when exiting the object production.
	ExitObject(c *ObjectContext)

	// ExitObject_element is called when exiting the object_element production.
	ExitObject_element(c *Object_elementContext)

	// ExitObject_element_name is called when exiting the object_element_name production.
	ExitObject_element_name(c *Object_element_nameContext)
}
