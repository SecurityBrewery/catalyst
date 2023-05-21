package caql

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/SecurityBrewery/catalyst/generated/caql/parser"
)

type Parser struct {
	Searcher Searcher
	Prefix   string
}

func (p *Parser) Parse(aql string) (t *Tree, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	// Set up the input
	inputStream := antlr.NewInputStream(aql)

	errorListener := &errorListener{}

	// Create the Lexer
	lexer := parser.NewCAQLLexer(inputStream)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errorListener)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the Parser
	aqlParser := parser.NewCAQLParser(stream)

	aqlParser.RemoveErrorListeners()
	aqlParser.AddErrorListener(errorListener)
	aqlParser.SetErrorHandler(antlr.NewBailErrorStrategy())
	if errorListener.errs != nil {
		err = errorListener.errs[0]
	}

	return &Tree{aqlParser: aqlParser, parseContext: aqlParser.Parse(), searcher: p.Searcher, prefix: p.Prefix}, err
}

type Tree struct {
	parseContext parser.IParseContext
	aqlParser    *parser.CAQLParser
	searcher     Searcher
	prefix       string
}

func (t *Tree) Eval(values map[string]any) (i any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	interpreter := aqlInterpreter{values: values}

	antlr.ParseTreeWalkerDefault.Walk(&interpreter, t.parseContext)

	if interpreter.errs != nil {
		return nil, interpreter.errs[0]
	}

	return interpreter.stack[0], nil
}

func (t *Tree) String() (s string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	builder := aqlBuilder{searcher: t.searcher, prefix: t.prefix}

	antlr.ParseTreeWalkerDefault.Walk(&builder, t.parseContext)

	return builder.stack[0], err
}

func (t *Tree) BleveString() (s string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	builder := bleveBuilder{}

	antlr.ParseTreeWalkerDefault.Walk(&builder, t.parseContext)

	if builder.err != nil {
		return "", builder.err
	}

	return builder.stack[0], err
}

type errorListener struct {
	*antlr.DefaultErrorListener
	errs []error
}

func (el *errorListener) SyntaxError(_ antlr.Recognizer, _ any, line, column int, msg string, _ antlr.RecognitionException) {
	el.errs = append(el.errs, fmt.Errorf("line "+strconv.Itoa(line)+":"+strconv.Itoa(column)+" "+msg))
}

func (el *errorListener) ReportAmbiguity(_ antlr.Parser, _ *antlr.DFA, _, _ int, _ bool, _ *antlr.BitSet, _ antlr.ATNConfigSet) {
	el.errs = append(el.errs, errors.New("ReportAmbiguity"))
}

func (el *errorListener) ReportAttemptingFullContext(_ antlr.Parser, _ *antlr.DFA, _, _ int, _ *antlr.BitSet, _ antlr.ATNConfigSet) {
	el.errs = append(el.errs, errors.New("ReportAttemptingFullContext"))
}

func (el *errorListener) ReportContextSensitivity(_ antlr.Parser, _ *antlr.DFA, _, _, _ int, _ antlr.ATNConfigSet) {
	el.errs = append(el.errs, errors.New("ReportContextSensitivity"))
}
