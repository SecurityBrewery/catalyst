package caql

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/SecurityBrewery/catalyst/generated/caql/parser"
)

var TooComplexError = errors.New("unsupported features for index queries, use advanced search instead")

type bleveBuilder struct {
	*parser.BaseCAQLParserListener
	stack []string
	err   error
}

// push is a helper function for pushing new node to the listener Stack.
func (s *bleveBuilder) push(i string) {
	s.stack = append(s.stack, i)
}

// pop is a helper function for poping a node from the listener Stack.
func (s *bleveBuilder) pop() (n string) {
	// Check that we have nodes in the stack.
	size := len(s.stack)
	if size < 1 {
		panic(ErrStack)
	}

	// Pop the last value from the Stack.
	n, s.stack = s.stack[size-1], s.stack[:size-1]

	return
}

func (s *bleveBuilder) binaryPop() (interface{}, interface{}) {
	right, left := s.pop(), s.pop()
	return left, right
}

// ExitExpression is called when production expression is exited.
func (s *bleveBuilder) ExitExpression(ctx *parser.ExpressionContext) {
	switch {
	case ctx.Value_literal() != nil:
		// pass
	case ctx.Reference() != nil:
		// pass
	case ctx.Operator_unary() != nil:
		s.err = TooComplexError
		return

	case ctx.T_PLUS() != nil:
		fallthrough
	case ctx.T_MINUS() != nil:
		fallthrough
	case ctx.T_TIMES() != nil:
		fallthrough
	case ctx.T_DIV() != nil:
		fallthrough
	case ctx.T_MOD() != nil:
		s.err = TooComplexError
		return

	case ctx.T_RANGE() != nil:
		s.err = TooComplexError
		return

	case ctx.T_LT() != nil && ctx.GetEq_op() == nil:
		left, right := s.binaryPop()
		s.push(fmt.Sprintf("%s:<%s", left, right))
	case ctx.T_GT() != nil && ctx.GetEq_op() == nil:
		left, right := s.binaryPop()
		s.push(fmt.Sprintf("%s:>%s", left, right))
	case ctx.T_LE() != nil && ctx.GetEq_op() == nil:
		left, right := s.binaryPop()
		s.push(fmt.Sprintf("%s:<=%s", left, right))
	case ctx.T_GE() != nil && ctx.GetEq_op() == nil:
		left, right := s.binaryPop()
		s.push(fmt.Sprintf("%s:>=%s", left, right))

	case ctx.T_IN() != nil && ctx.GetEq_op() == nil:
		s.err = TooComplexError
		return

	case ctx.T_EQ() != nil && ctx.GetEq_op() == nil:
		left, right := s.binaryPop()
		s.push(fmt.Sprintf("%s:%s", left, right))
	case ctx.T_NE() != nil && ctx.GetEq_op() == nil:
		left, right := s.binaryPop()
		s.push(fmt.Sprintf("-%s:%s", left, right))

	case ctx.T_ALL() != nil && ctx.GetEq_op() != nil:
		fallthrough
	case ctx.T_ANY() != nil && ctx.GetEq_op() != nil:
		fallthrough
	case ctx.T_NONE() != nil && ctx.GetEq_op() != nil:
		s.err = TooComplexError
		return

	case ctx.T_ALL() != nil && ctx.T_NOT() != nil && ctx.T_IN() != nil:
		fallthrough
	case ctx.T_ANY() != nil && ctx.T_NOT() != nil && ctx.T_IN() != nil:
		fallthrough
	case ctx.T_NONE() != nil && ctx.T_NOT() != nil && ctx.T_IN() != nil:
		s.err = TooComplexError
		return

	case ctx.T_LIKE() != nil:
		s.err = errors.New("index queries are like queries by default")
		return

	case ctx.T_REGEX_MATCH() != nil:
		left, right := s.binaryPop()
		if ctx.T_NOT() != nil {
			s.err = TooComplexError
			return
		} else {
			s.push(fmt.Sprintf("%s:/%s/", left, right))
		}
	case ctx.T_REGEX_NON_MATCH() != nil:
		s.err = errors.New("index query cannot contain regex non matches, use advanced search instead")
		return

	case ctx.T_AND() != nil:
		left, right := s.binaryPop()
		s.push(fmt.Sprintf("%s %s", left, right))
	case ctx.T_OR() != nil:
		s.err = errors.New("index query cannot contain OR, use advanced search instead")
		return

	case ctx.T_QUESTION() != nil && len(ctx.AllExpression()) == 3:
		s.err = errors.New("index query cannot contain ternary operations, use advanced search instead")
		return
	case ctx.T_QUESTION() != nil && len(ctx.AllExpression()) == 2:
		s.err = errors.New("index query cannot contain ternary operations, use advanced search instead")
		return

	default:
		panic("unknown expression")
	}
}

// ExitReference is called when production reference is exited.
func (s *bleveBuilder) ExitReference(ctx *parser.ReferenceContext) {
	switch {
	case ctx.DOT() != nil:
		reference := s.pop()

		s.push(fmt.Sprintf("%s.%s", reference, ctx.T_STRING().GetText()))
	case ctx.T_STRING() != nil:
		s.push(ctx.T_STRING().GetText())
	case ctx.Compound_value() != nil:
		s.err = TooComplexError
		return
	case ctx.Function_call() != nil:
		s.err = TooComplexError
		return
	case ctx.T_OPEN() != nil:
		s.err = TooComplexError
		return
	case ctx.T_ARRAY_OPEN() != nil:
		s.err = TooComplexError
		return
	default:
		panic(fmt.Sprintf("unexpected value: %s", ctx.GetText()))
	}
}

// ExitValue_literal is called when production value_literal is exited.
func (s *bleveBuilder) ExitValue_literal(ctx *parser.Value_literalContext) {
	if ctx.T_QUOTED_STRING() != nil {
		st, err := unquote(ctx.GetText())
		if err != nil {
			panic(err)
		}
		s.push(strconv.Quote(st))
	} else {
		s.push(ctx.GetText())
	}
}
