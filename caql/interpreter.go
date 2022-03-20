package caql

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SecurityBrewery/catalyst/generated/caql/parser"
)

type aqlInterpreter struct {
	*parser.BaseCAQLParserListener
	values map[string]any
	stack  []any
	errs   []error
}

// push is a helper function for pushing new node to the listener Stack.
func (s *aqlInterpreter) push(i any) {
	s.stack = append(s.stack, i)
}

// pop is a helper function for poping a node from the listener Stack.
func (s *aqlInterpreter) pop() (n any) {
	// Check that we have nodes in the stack.
	size := len(s.stack)
	if size < 1 {
		s.appendErrors(ErrStack)

		return
	}

	// Pop the last value from the Stack.
	n, s.stack = s.stack[size-1], s.stack[:size-1]

	return
}

func (s *aqlInterpreter) binaryPop() (any, any) {
	right, left := s.pop(), s.pop()

	return left, right
}

// ExitExpression is called when production expression is exited.
func (s *aqlInterpreter) ExitExpression(ctx *parser.ExpressionContext) {
	switch {
	case ctx.Value_literal() != nil:
		// pass
	case ctx.Reference() != nil:
		// pass
	case ctx.Operator_unary() != nil:
		// pass

	case ctx.T_PLUS() != nil:
		s.push(plus(s.binaryPop()))
	case ctx.T_MINUS() != nil:
		s.push(minus(s.binaryPop()))
	case ctx.T_TIMES() != nil:
		s.push(times(s.binaryPop()))
	case ctx.T_DIV() != nil:
		s.push(div(s.binaryPop()))
	case ctx.T_MOD() != nil:
		s.push(mod(s.binaryPop()))
	case ctx.T_RANGE() != nil:
		s.push(aqlrange(s.binaryPop()))
	case ctx.T_LT() != nil && ctx.GetEq_op() == nil:
		s.push(lt(s.binaryPop()))
	case ctx.T_GT() != nil && ctx.GetEq_op() == nil:
		s.push(gt(s.binaryPop()))
	case ctx.T_LE() != nil && ctx.GetEq_op() == nil:
		s.push(le(s.binaryPop()))
	case ctx.T_GE() != nil && ctx.GetEq_op() == nil:
		s.push(ge(s.binaryPop()))
	case ctx.T_IN() != nil && ctx.GetEq_op() == nil:
		s.push(maybeNot(ctx, in(s.binaryPop())))
	case ctx.T_EQ() != nil && ctx.GetEq_op() == nil:
		s.push(eq(s.binaryPop()))
	case ctx.T_NE() != nil && ctx.GetEq_op() == nil:
		s.push(ne(s.binaryPop()))
	case ctx.T_ALL() != nil && ctx.GetEq_op() != nil:
		right, left := s.pop(), s.pop()
		s.push(all(left.([]any), getOp(ctx.GetEq_op().GetTokenType()), right))
	case ctx.T_ANY() != nil && ctx.GetEq_op() != nil:
		right, left := s.pop(), s.pop()
		s.push(anyElement(left.([]any), getOp(ctx.GetEq_op().GetTokenType()), right))
	case ctx.T_NONE() != nil && ctx.GetEq_op() != nil:
		right, left := s.pop(), s.pop()
		s.push(none(left.([]any), getOp(ctx.GetEq_op().GetTokenType()), right))
	case ctx.T_ALL() != nil && ctx.T_NOT() != nil && ctx.T_IN() != nil:
		right, left := s.pop(), s.pop()
		s.push(all(left.([]any), in, right))
	case ctx.T_ANY() != nil && ctx.T_NOT() != nil && ctx.T_IN() != nil:
		right, left := s.pop(), s.pop()
		s.push(anyElement(left.([]any), in, right))
	case ctx.T_NONE() != nil && ctx.T_NOT() != nil && ctx.T_IN() != nil:
		right, left := s.pop(), s.pop()
		s.push(none(left.([]any), in, right))
	case ctx.T_LIKE() != nil:
		m, err := like(s.binaryPop())
		s.appendErrors(err)
		s.push(maybeNot(ctx, m))
	case ctx.T_REGEX_MATCH() != nil:
		m, err := regexMatch(s.binaryPop())
		s.appendErrors(err)
		s.push(maybeNot(ctx, m))
	case ctx.T_REGEX_NON_MATCH() != nil:
		m, err := regexNonMatch(s.binaryPop())
		s.appendErrors(err)
		s.push(maybeNot(ctx, m))
	case ctx.T_AND() != nil:
		s.push(and(s.binaryPop()))
	case ctx.T_OR() != nil:
		s.push(or(s.binaryPop()))
	case ctx.T_QUESTION() != nil && len(ctx.AllExpression()) == 3:
		right, middle, left := s.pop(), s.pop(), s.pop()
		s.push(ternary(left, middle, right))
	case ctx.T_QUESTION() != nil && len(ctx.AllExpression()) == 2:
		right, left := s.pop(), s.pop()
		s.push(ternary(left, nil, right))
	default:
		panic("unknown expression")
	}
}

func (s *aqlInterpreter) appendErrors(err error) {
	if err != nil {
		s.errs = append(s.errs, err)
	}
}

// ExitOperator_unary is called when production operator_unary is exited.
func (s *aqlInterpreter) ExitOperator_unary(ctx *parser.Operator_unaryContext) {
	value := s.pop()
	switch {
	case ctx.T_PLUS() != nil:
		s.push(value.(float64))
	case ctx.T_MINUS() != nil:
		s.push(-value.(float64))
	case ctx.T_NOT() != nil:
		s.push(!toBool(value))
	default:
		panic(fmt.Sprintf("unexpected operation: %s", ctx.GetText()))
	}
}

// ExitReference is called when production reference is exited.
func (s *aqlInterpreter) ExitReference(ctx *parser.ReferenceContext) {
	switch {
	case ctx.DOT() != nil:
		reference := s.pop()

		s.push(reference.(map[string]any)[ctx.T_STRING().GetText()])
	case ctx.T_STRING() != nil:
		s.push(s.getVar(ctx.T_STRING().GetText()))
	case ctx.Compound_value() != nil:
		// pass
	case ctx.Function_call() != nil:
		// pass
	case ctx.T_OPEN() != nil:
		// pass
	case ctx.T_ARRAY_OPEN() != nil:
		key := s.pop()
		reference := s.pop()

		if f, ok := key.(float64); ok {
			index := int(f)
			if index < 0 {
				index = len(reference.([]any)) + index
			}

			s.push(reference.([]any)[index])

			return
		}

		s.push(reference.(map[string]any)[key.(string)])
	default:
		panic(fmt.Sprintf("unexpected value: %s", ctx.GetText()))
	}
}

// ExitCompound_value is called when production compound_value is exited.
func (s *aqlInterpreter) ExitCompound_value(ctx *parser.Compound_valueContext) {
	// pass
}

// ExitFunction_call is called when production function_call is exited.
func (s *aqlInterpreter) ExitFunction_call(ctx *parser.Function_callContext) {
	s.function(ctx)
}

// ExitValue_literal is called when production value_literal is exited.
func (s *aqlInterpreter) ExitValue_literal(ctx *parser.Value_literalContext) {
	switch {
	case ctx.T_QUOTED_STRING() != nil:
		st, err := unquote(ctx.GetText())
		s.appendErrors(err)
		s.push(st)
	case ctx.T_INT() != nil:
		t := ctx.GetText()

		switch {
		case strings.HasPrefix(strings.ToLower(t), "0b"):
			i64, err := strconv.ParseInt(t[2:], 2, 64)
			s.appendErrors(err)
			s.push(float64(i64))
		case strings.HasPrefix(strings.ToLower(t), "0x"):
			i64, err := strconv.ParseInt(t[2:], 16, 64)
			s.appendErrors(err)
			s.push(float64(i64))
		default:
			i, err := strconv.Atoi(t)
			s.appendErrors(err)
			s.push(float64(i))
		}
	case ctx.T_FLOAT() != nil:
		i, err := strconv.ParseFloat(ctx.GetText(), 64)
		s.appendErrors(err)
		s.push(i)
	case ctx.T_NULL() != nil:
		s.push(nil)
	case ctx.T_TRUE() != nil:
		s.push(true)
	case ctx.T_FALSE() != nil:
		s.push(false)
	default:
		panic(fmt.Sprintf("unexpected value: %s", ctx.GetText()))
	}
}

// ExitArray is called when production array is exited.
func (s *aqlInterpreter) ExitArray(ctx *parser.ArrayContext) {
	array := []any{}
	for range ctx.AllExpression() {
		// prepend element
		array = append([]any{s.pop()}, array...)
	}
	s.push(array)
}

// ExitObject is called when production object is exited.
func (s *aqlInterpreter) ExitObject(ctx *parser.ObjectContext) {
	object := map[string]any{}
	for range ctx.AllObject_element() {
		key, value := s.pop(), s.pop()

		object[key.(string)] = value
	}
	s.push(object)
}

// ExitObject_element is called when production object_element is exited.
func (s *aqlInterpreter) ExitObject_element(ctx *parser.Object_elementContext) {
	switch {
	case ctx.T_STRING() != nil:
		s.push(ctx.GetText())
		s.push(s.getVar(ctx.GetText()))
	case ctx.Object_element_name() != nil, ctx.T_ARRAY_OPEN() != nil:
		key, value := s.pop(), s.pop()

		s.push(key)
		s.push(value)
	default:
		panic(fmt.Sprintf("unexpected value: %s", ctx.GetText()))
	}
}

// ExitObject_element_name is called when production object_element_name is exited.
func (s *aqlInterpreter) ExitObject_element_name(ctx *parser.Object_element_nameContext) {
	switch {
	case ctx.T_STRING() != nil:
		s.push(ctx.T_STRING().GetText())
	case ctx.T_QUOTED_STRING() != nil:
		st, err := unquote(ctx.T_QUOTED_STRING().GetText())
		if err != nil {
			s.appendErrors(fmt.Errorf("%w: %s", err, ctx.GetText()))
		}
		s.push(st)
	default:
		panic(fmt.Sprintf("unexpected value: %s", ctx.GetText()))
	}
}

func (s *aqlInterpreter) getVar(identifier string) any {
	v, ok := s.values[identifier]
	if !ok {
		s.appendErrors(ErrUndefined)
	}

	return v
}

func maybeNot(ctx *parser.ExpressionContext, m bool) bool {
	if ctx.T_NOT() != nil {
		return !m
	}

	return m
}

func getOp(tokenType int) func(left, right any) bool {
	switch tokenType {
	case parser.CAQLLexerT_EQ:
		return eq
	case parser.CAQLLexerT_NE:
		return ne
	case parser.CAQLLexerT_LT:
		return lt
	case parser.CAQLLexerT_GT:
		return gt
	case parser.CAQLLexerT_LE:
		return le
	case parser.CAQLLexerT_GE:
		return ge
	case parser.CAQLLexerT_IN:
		return in
	default:
		panic("unknown token type")
	}
}

func all(slice []any, op func(any, any) bool, expr any) bool {
	for _, e := range slice {
		if !op(e, expr) {
			return false
		}
	}

	return true
}

func anyElement(slice []any, op func(any, any) bool, expr any) bool {
	for _, e := range slice {
		if op(e, expr) {
			return true
		}
	}

	return false
}

func none(slice []any, op func(any, any) bool, expr any) bool {
	for _, e := range slice {
		if op(e, expr) {
			return false
		}
	}

	return true
}
