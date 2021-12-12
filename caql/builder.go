package caql

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SecurityBrewery/catalyst/generated/caql/parser"
)

type Searcher interface {
	Search(term string) (ids []string, err error)
}

type aqlBuilder struct {
	*parser.BaseCAQLParserListener
	searcher Searcher
	stack    []string
	prefix   string
}

// push is a helper function for pushing new node to the listener Stack.
func (s *aqlBuilder) push(i string) {
	s.stack = append(s.stack, i)
}

// pop is a helper function for poping a node from the listener Stack.
func (s *aqlBuilder) pop() (n string) {
	// Check that we have nodes in the stack.
	size := len(s.stack)
	if size < 1 {
		panic(ErrStack)
	}

	// Pop the last value from the Stack.
	n, s.stack = s.stack[size-1], s.stack[:size-1]

	return
}

func (s *aqlBuilder) binaryPop() (string, string) {
	right, left := s.pop(), s.pop()
	return left, right
}

// ExitExpression is called when production expression is exited.
func (s *aqlBuilder) ExitExpression(ctx *parser.ExpressionContext) {
	switch {
	case ctx.Value_literal() != nil:
		if ctx.GetParent().GetParent() == nil {
			s.push(s.toBoolString(s.pop()))
		}
	case ctx.Reference() != nil:
		ref := s.pop()
		if ref == "d.id" {
			s.push("d._key")
		} else {
			s.push(ref)
		}
		// pass
	case ctx.Operator_unary() != nil:
		s.push(s.toBoolString(s.pop()))

	case ctx.T_PLUS() != nil:
		left, right := s.binaryPop()
		s.push(fmt.Sprintf("%s + %s", left, right))
	case ctx.T_MINUS() != nil:
		left, right := s.binaryPop()
		s.push(fmt.Sprintf("%s - %s", left, right))
	case ctx.T_TIMES() != nil:
		left, right := s.binaryPop()
		s.push(fmt.Sprintf("%s * %s", left, right))
	case ctx.T_DIV() != nil:
		left, right := s.binaryPop()
		s.push(fmt.Sprintf("%s / %s", left, right))
	case ctx.T_MOD() != nil:
		left, right := s.binaryPop()
		s.push(fmt.Sprintf("%s %% %s", left, right))

	case ctx.T_RANGE() != nil:
		left, right := s.binaryPop()
		s.push(fmt.Sprintf("%s..%s", left, right))

	case ctx.T_LT() != nil && ctx.GetEq_op() == nil:
		left, right := s.binaryPop()
		s.push(fmt.Sprintf("%s < %s", left, right))
	case ctx.T_GT() != nil && ctx.GetEq_op() == nil:
		left, right := s.binaryPop()
		s.push(fmt.Sprintf("%s > %s", left, right))
	case ctx.T_LE() != nil && ctx.GetEq_op() == nil:
		left, right := s.binaryPop()
		s.push(fmt.Sprintf("%s <= %s", left, right))
	case ctx.T_GE() != nil && ctx.GetEq_op() == nil:
		left, right := s.binaryPop()
		s.push(fmt.Sprintf("%s >= %s", left, right))

	case ctx.T_IN() != nil && ctx.GetEq_op() == nil:
		left, right := s.binaryPop()
		if ctx.T_NOT() != nil {
			s.push(fmt.Sprintf("%s NOT IN %s", left, right))
		} else {
			s.push(fmt.Sprintf("%s IN %s", left, right))
		}

	case ctx.T_EQ() != nil && ctx.GetEq_op() == nil:
		left, right := s.binaryPop()
		s.push(fmt.Sprintf("%s == %s", left, right))
	case ctx.T_NE() != nil && ctx.GetEq_op() == nil:
		left, right := s.binaryPop()
		s.push(fmt.Sprintf("%s != %s", left, right))

	case ctx.T_ALL() != nil && ctx.GetEq_op() != nil:
		right, left := s.pop(), s.pop()
		s.push(fmt.Sprintf("%s ALL %s %s", left, ctx.GetEq_op().GetText(), right))
	case ctx.T_ANY() != nil && ctx.GetEq_op() != nil:
		right, left := s.pop(), s.pop()
		s.push(fmt.Sprintf("%s ANY %s %s", left, ctx.GetEq_op().GetText(), right))
	case ctx.T_NONE() != nil && ctx.GetEq_op() != nil:
		right, left := s.pop(), s.pop()
		s.push(fmt.Sprintf("%s NONE %s %s", left, ctx.GetEq_op().GetText(), right))

	case ctx.T_ALL() != nil && ctx.T_NOT() != nil && ctx.T_IN() != nil:
		right, left := s.pop(), s.pop()
		s.push(fmt.Sprintf("%s ALL IN %s", left, right))
	case ctx.T_ANY() != nil && ctx.T_NOT() != nil && ctx.T_IN() != nil:
		right, left := s.pop(), s.pop()
		s.push(fmt.Sprintf("%s ANY IN %s", left, right))
	case ctx.T_NONE() != nil && ctx.T_NOT() != nil && ctx.T_IN() != nil:
		right, left := s.pop(), s.pop()
		s.push(fmt.Sprintf("%s NONE IN %s", left, right))

	case ctx.T_LIKE() != nil:
		left, right := s.binaryPop()
		if ctx.T_NOT() != nil {
			s.push(fmt.Sprintf("%s NOT LIKE %s", left, right))
		} else {
			s.push(fmt.Sprintf("%s LIKE %s", left, right))
		}
	case ctx.T_REGEX_MATCH() != nil:
		left, right := s.binaryPop()
		if ctx.T_NOT() != nil {
			s.push(fmt.Sprintf("%s NOT =~ %s", left, right))
		} else {
			s.push(fmt.Sprintf("%s =~ %s", left, right))
		}
	case ctx.T_REGEX_NON_MATCH() != nil:
		left, right := s.binaryPop()
		if ctx.T_NOT() != nil {
			s.push(fmt.Sprintf("%s NOT !~ %s", left, right))
		} else {
			s.push(fmt.Sprintf("%s !~ %s", left, right))
		}

	case ctx.T_AND() != nil:
		left, right := s.binaryPop()
		left = s.toBoolString(left)
		right = s.toBoolString(right)
		s.push(fmt.Sprintf("%s AND %s", left, right))
	case ctx.T_OR() != nil:
		left, right := s.binaryPop()
		left = s.toBoolString(left)
		right = s.toBoolString(right)
		s.push(fmt.Sprintf("%s OR %s", left, right))

	case ctx.T_QUESTION() != nil && len(ctx.AllExpression()) == 3:
		right, middle, left := s.pop(), s.pop(), s.pop()
		s.push(fmt.Sprintf("%s ? %s : %s", left, middle, right))
	case ctx.T_QUESTION() != nil && len(ctx.AllExpression()) == 2:
		right, left := s.pop(), s.pop()
		s.push(fmt.Sprintf("%s ? : %s", left, right))

	default:
		panic("unknown expression")
	}
}

func (s *aqlBuilder) toBoolString(v string) string {
	_, err := unquote(v)
	if err == nil {
		ids, err := s.searcher.Search(v)
		if err != nil {
			panic("invalid search " + err.Error())
		}
		return fmt.Sprintf(`d._key IN ["%s"]`, strings.Join(ids, `","`))
	}
	return v
}

// ExitOperator_unary is called when production operator_unary is exited.
func (s *aqlBuilder) ExitOperator_unary(ctx *parser.Operator_unaryContext) {
	value := s.pop()
	switch {
	case ctx.T_PLUS() != nil:
		s.push(value)
	case ctx.T_MINUS() != nil:
		s.push(fmt.Sprintf("-%s", value))
	case ctx.T_NOT() != nil:
		s.push(fmt.Sprintf("NOT %s", value))
	default:
		panic(fmt.Sprintf("unexpected operation: %s", ctx.GetText()))
	}
}

// ExitReference is called when production reference is exited.
func (s *aqlBuilder) ExitReference(ctx *parser.ReferenceContext) {
	switch {
	case ctx.DOT() != nil:
		reference := s.pop()
		if s.prefix != "" && !strings.HasPrefix(reference, s.prefix) {
			reference = s.prefix + reference
		}
		s.push(fmt.Sprintf("%s.%s", reference, ctx.T_STRING().GetText()))
	case ctx.T_STRING() != nil:
		reference := ctx.T_STRING().GetText()
		if s.prefix != "" && !strings.HasPrefix(reference, s.prefix) {
			reference = s.prefix + reference
		}
		s.push(reference)
	case ctx.Compound_value() != nil:
		// pass
	case ctx.Function_call() != nil:
		// pass
	case ctx.T_OPEN() != nil:
		s.push(fmt.Sprintf("(%s)", s.pop()))
	case ctx.T_ARRAY_OPEN() != nil:
		key := s.pop()
		reference := s.pop()

		s.push(fmt.Sprintf("%s[%s]", reference, key))
	default:
		panic(fmt.Sprintf("unexpected value: %s", ctx.GetText()))
	}
}

// ExitCompound_value is called when production compound_value is exited.
func (s *aqlBuilder) ExitCompound_value(ctx *parser.Compound_valueContext) {
	// pass
}

// ExitFunction_call is called when production function_call is exited.
func (s *aqlBuilder) ExitFunction_call(ctx *parser.Function_callContext) {
	var array []string
	for range ctx.AllExpression() {
		// prepend element
		array = append([]string{s.pop()}, array...)
	}
	parameter := strings.Join(array, ", ")

	if !stringSliceContains(functionNames, strings.ToUpper(ctx.T_STRING().GetText())) {
		panic("unknown function")
	}

	s.push(fmt.Sprintf("%s(%s)", strings.ToUpper(ctx.T_STRING().GetText()), parameter))
}

// ExitValue_literal is called when production value_literal is exited.
func (s *aqlBuilder) ExitValue_literal(ctx *parser.Value_literalContext) {
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

// ExitArray is called when production array is exited.
func (s *aqlBuilder) ExitArray(ctx *parser.ArrayContext) {
	var elements []string
	for range ctx.AllExpression() {
		// elements = append(elements, s.pop())
		elements = append([]string{s.pop()}, elements...)
	}
	s.push("[" + strings.Join(elements, ", ") + "]")
}

// ExitObject is called when production object is exited.
func (s *aqlBuilder) ExitObject(ctx *parser.ObjectContext) {
	var elements []string
	for range ctx.AllObject_element() {
		key, value := s.pop(), s.pop()

		elements = append([]string{fmt.Sprintf("%s: %v", key, value)}, elements...)
	}
	// s.push(object)
	s.push("{" + strings.Join(elements, ", ") + "}")
}

// ExitObject_element is called when production object_element is exited.
func (s *aqlBuilder) ExitObject_element(ctx *parser.Object_elementContext) {
	switch {
	case ctx.T_STRING() != nil:
		s.push(ctx.GetText())
		s.push(ctx.GetText())
	case ctx.Object_element_name() != nil, ctx.T_ARRAY_OPEN() != nil:
		key, value := s.pop(), s.pop()

		s.push(key)
		s.push(value)
	default:
		panic(fmt.Sprintf("unexpected value: %s", ctx.GetText()))
	}
}

// ExitObject_element_name is called when production object_element_name is exited.
func (s *aqlBuilder) ExitObject_element_name(ctx *parser.Object_element_nameContext) {
	switch {
	case ctx.T_STRING() != nil:
		s.push(ctx.T_STRING().GetText())
	case ctx.T_QUOTED_STRING() != nil:
		s.push(ctx.T_QUOTED_STRING().GetText())
	default:
		panic(fmt.Sprintf("unexpected value: %s", ctx.GetText()))
	}
}
