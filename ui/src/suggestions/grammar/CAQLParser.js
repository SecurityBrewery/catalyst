// Generated from CAQLParser.g4 by ANTLR 4.10.1
// jshint ignore: start
import antlr4 from 'antlr4';
import CAQLParserListener from './CAQLParserListener.js';
const serializedATN = [4,1,78,190,2,0,7,0,2,1,7,1,2,2,7,2,2,3,7,3,2,4,7,
4,2,5,7,5,2,6,7,6,2,7,7,7,2,8,7,8,2,9,7,9,2,10,7,10,1,0,1,0,1,0,1,1,1,1,
1,1,1,1,3,1,30,8,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,
1,3,1,46,8,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,
1,1,3,1,64,8,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,
1,1,1,1,1,1,1,5,1,84,8,1,10,1,12,1,87,9,1,1,2,1,2,1,2,1,2,1,2,1,2,3,2,95,
8,2,1,3,1,3,1,3,1,3,1,3,1,3,1,3,1,3,3,3,105,8,3,1,3,1,3,1,3,1,3,1,3,1,3,
1,3,1,3,5,3,115,8,3,10,3,12,3,118,9,3,1,4,1,4,3,4,122,8,4,1,5,1,5,1,5,3,
5,127,8,5,1,5,1,5,5,5,131,8,5,10,5,12,5,134,9,5,1,5,3,5,137,8,5,1,5,1,5,
1,6,1,6,1,7,1,7,3,7,145,8,7,1,7,1,7,5,7,149,8,7,10,7,12,7,152,9,7,1,7,3,
7,155,8,7,1,7,1,7,1,8,1,8,3,8,161,8,8,1,8,1,8,5,8,165,8,8,10,8,12,8,168,
9,8,1,8,3,8,171,8,8,1,8,1,8,1,9,1,9,1,9,1,9,1,9,1,9,1,9,1,9,1,9,1,9,1,9,
3,9,186,8,9,1,10,1,10,1,10,2,132,150,2,2,6,11,0,2,4,6,8,10,12,14,16,18,20,
0,9,1,0,10,11,1,0,12,14,1,0,6,9,1,0,4,5,3,0,27,27,29,29,46,46,2,0,4,9,38,
38,2,0,2,3,44,44,5,0,34,34,48,48,56,56,70,71,73,73,2,0,69,69,73,73,214,0,
22,1,0,0,0,2,29,1,0,0,0,4,94,1,0,0,0,6,104,1,0,0,0,8,121,1,0,0,0,10,123,
1,0,0,0,12,140,1,0,0,0,14,142,1,0,0,0,16,158,1,0,0,0,18,185,1,0,0,0,20,187,
1,0,0,0,22,23,3,2,1,0,23,24,5,0,0,1,24,1,1,0,0,0,25,26,6,1,-1,0,26,30,3,
12,6,0,27,30,3,6,3,0,28,30,3,4,2,0,29,25,1,0,0,0,29,27,1,0,0,0,29,28,1,0,
0,0,30,85,1,0,0,0,31,32,10,13,0,0,32,33,7,0,0,0,33,84,3,2,1,14,34,35,10,
12,0,0,35,36,7,1,0,0,36,84,3,2,1,13,37,38,10,11,0,0,38,39,5,18,0,0,39,84,
3,2,1,12,40,41,10,10,0,0,41,42,7,2,0,0,42,84,3,2,1,11,43,45,10,9,0,0,44,
46,5,47,0,0,45,44,1,0,0,0,45,46,1,0,0,0,46,47,1,0,0,0,47,48,5,38,0,0,48,
84,3,2,1,10,49,50,10,8,0,0,50,51,7,3,0,0,51,84,3,2,1,9,52,53,10,7,0,0,53,
54,7,4,0,0,54,55,7,5,0,0,55,84,3,2,1,8,56,57,10,6,0,0,57,58,7,4,0,0,58,59,
5,47,0,0,59,60,5,38,0,0,60,84,3,2,1,7,61,63,10,5,0,0,62,64,5,47,0,0,63,62,
1,0,0,0,63,64,1,0,0,0,64,65,1,0,0,0,65,66,7,6,0,0,66,84,3,2,1,6,67,68,10,
4,0,0,68,69,5,28,0,0,69,84,3,2,1,5,70,71,10,3,0,0,71,72,5,49,0,0,72,84,3,
2,1,4,73,74,10,2,0,0,74,75,5,15,0,0,75,76,3,2,1,0,76,77,5,16,0,0,77,78,3,
2,1,3,78,84,1,0,0,0,79,80,10,1,0,0,80,81,5,15,0,0,81,82,5,16,0,0,82,84,3,
2,1,2,83,31,1,0,0,0,83,34,1,0,0,0,83,37,1,0,0,0,83,40,1,0,0,0,83,43,1,0,
0,0,83,49,1,0,0,0,83,52,1,0,0,0,83,56,1,0,0,0,83,61,1,0,0,0,83,67,1,0,0,
0,83,70,1,0,0,0,83,73,1,0,0,0,83,79,1,0,0,0,84,87,1,0,0,0,85,83,1,0,0,0,
85,86,1,0,0,0,86,3,1,0,0,0,87,85,1,0,0,0,88,89,5,10,0,0,89,95,3,2,1,0,90,
91,5,11,0,0,91,95,3,2,1,0,92,93,5,47,0,0,93,95,3,2,1,0,94,88,1,0,0,0,94,
90,1,0,0,0,94,92,1,0,0,0,95,5,1,0,0,0,96,97,6,3,-1,0,97,105,5,69,0,0,98,
105,3,8,4,0,99,105,3,10,5,0,100,101,5,20,0,0,101,102,3,2,1,0,102,103,5,21,
0,0,103,105,1,0,0,0,104,96,1,0,0,0,104,98,1,0,0,0,104,99,1,0,0,0,104,100,
1,0,0,0,105,116,1,0,0,0,106,107,10,2,0,0,107,108,5,1,0,0,108,115,5,69,0,
0,109,110,10,1,0,0,110,111,5,24,0,0,111,112,3,2,1,0,112,113,5,25,0,0,113,
115,1,0,0,0,114,106,1,0,0,0,114,109,1,0,0,0,115,118,1,0,0,0,116,114,1,0,
0,0,116,117,1,0,0,0,117,7,1,0,0,0,118,116,1,0,0,0,119,122,3,14,7,0,120,122,
3,16,8,0,121,119,1,0,0,0,121,120,1,0,0,0,122,9,1,0,0,0,123,124,5,69,0,0,
124,126,5,20,0,0,125,127,3,2,1,0,126,125,1,0,0,0,126,127,1,0,0,0,127,132,
1,0,0,0,128,129,5,19,0,0,129,131,3,2,1,0,130,128,1,0,0,0,131,134,1,0,0,0,
132,133,1,0,0,0,132,130,1,0,0,0,133,136,1,0,0,0,134,132,1,0,0,0,135,137,
5,19,0,0,136,135,1,0,0,0,136,137,1,0,0,0,137,138,1,0,0,0,138,139,5,21,0,
0,139,11,1,0,0,0,140,141,7,7,0,0,141,13,1,0,0,0,142,144,5,24,0,0,143,145,
3,2,1,0,144,143,1,0,0,0,144,145,1,0,0,0,145,150,1,0,0,0,146,147,5,19,0,0,
147,149,3,2,1,0,148,146,1,0,0,0,149,152,1,0,0,0,150,151,1,0,0,0,150,148,
1,0,0,0,151,154,1,0,0,0,152,150,1,0,0,0,153,155,5,19,0,0,154,153,1,0,0,0,
154,155,1,0,0,0,155,156,1,0,0,0,156,157,5,25,0,0,157,15,1,0,0,0,158,160,
5,22,0,0,159,161,3,18,9,0,160,159,1,0,0,0,160,161,1,0,0,0,161,166,1,0,0,
0,162,163,5,19,0,0,163,165,3,18,9,0,164,162,1,0,0,0,165,168,1,0,0,0,166,
164,1,0,0,0,166,167,1,0,0,0,167,170,1,0,0,0,168,166,1,0,0,0,169,171,5,19,
0,0,170,169,1,0,0,0,170,171,1,0,0,0,171,172,1,0,0,0,172,173,5,23,0,0,173,
17,1,0,0,0,174,186,5,69,0,0,175,176,3,20,10,0,176,177,5,16,0,0,177,178,3,
2,1,0,178,186,1,0,0,0,179,180,5,24,0,0,180,181,3,2,1,0,181,182,5,25,0,0,
182,183,5,16,0,0,183,184,3,2,1,0,184,186,1,0,0,0,185,174,1,0,0,0,185,175,
1,0,0,0,185,179,1,0,0,0,186,19,1,0,0,0,187,188,7,8,0,0,188,21,1,0,0,0,20,
29,45,63,83,85,94,104,114,116,121,126,132,136,144,150,154,160,166,170,185];


const atn = new antlr4.atn.ATNDeserializer().deserialize(serializedATN);

const decisionsToDFA = atn.decisionToState.map( (ds, index) => new antlr4.dfa.DFA(ds, index) );

const sharedContextCache = new antlr4.PredictionContextCache();

export default class CAQLParser extends antlr4.Parser {

    static grammarFileName = "CAQLParser.g4";
    static literalNames = [ null, "'.'", "'=~'", "'!~'", "'=='", "'!='", 
                            "'<'", "'>'", "'<='", "'>='", "'+'", "'-'", 
                            "'*'", "'/'", "'%'", "'?'", "':'", "'::'", "'..'", 
                            "','", "'('", "')'", "'{'", "'}'", "'['", "']'" ];
    static symbolicNames = [ null, "DOT", "T_REGEX_MATCH", "T_REGEX_NON_MATCH", 
                             "T_EQ", "T_NE", "T_LT", "T_GT", "T_LE", "T_GE", 
                             "T_PLUS", "T_MINUS", "T_TIMES", "T_DIV", "T_MOD", 
                             "T_QUESTION", "T_COLON", "T_SCOPE", "T_RANGE", 
                             "T_COMMA", "T_OPEN", "T_CLOSE", "T_OBJECT_OPEN", 
                             "T_OBJECT_CLOSE", "T_ARRAY_OPEN", "T_ARRAY_CLOSE", 
                             "T_AGGREGATE", "T_ALL", "T_AND", "T_ANY", "T_ASC", 
                             "T_COLLECT", "T_DESC", "T_DISTINCT", "T_FALSE", 
                             "T_FILTER", "T_FOR", "T_GRAPH", "T_IN", "T_INBOUND", 
                             "T_INSERT", "T_INTO", "T_K_SHORTEST_PATHS", 
                             "T_LET", "T_LIKE", "T_LIMIT", "T_NONE", "T_NOT", 
                             "T_NULL", "T_OR", "T_OUTBOUND", "T_REMOVE", 
                             "T_REPLACE", "T_RETURN", "T_SHORTEST_PATH", 
                             "T_SORT", "T_TRUE", "T_UPDATE", "T_UPSERT", 
                             "T_WITH", "T_KEEP", "T_COUNT", "T_OPTIONS", 
                             "T_PRUNE", "T_SEARCH", "T_TO", "T_CURRENT", 
                             "T_NEW", "T_OLD", "T_STRING", "T_INT", "T_FLOAT", 
                             "T_PARAMETER", "T_QUOTED_STRING", "SINGLE_LINE_COMMENT", 
                             "MULTILINE_COMMENT", "SPACES", "UNEXPECTED_CHAR", 
                             "ERROR_RECONGNIGION" ];
    static ruleNames = [ "parse", "expression", "operator_unary", "reference", 
                         "compound_value", "function_call", "value_literal", 
                         "array", "object", "object_element", "object_element_name" ];

    constructor(input) {
        super(input);
        this._interp = new antlr4.atn.ParserATNSimulator(this, atn, decisionsToDFA, sharedContextCache);
        this.ruleNames = CAQLParser.ruleNames;
        this.literalNames = CAQLParser.literalNames;
        this.symbolicNames = CAQLParser.symbolicNames;
    }

    get atn() {
        return atn;
    }

    sempred(localctx, ruleIndex, predIndex) {
    	switch(ruleIndex) {
    	case 1:
    	    		return this.expression_sempred(localctx, predIndex);
    	case 3:
    	    		return this.reference_sempred(localctx, predIndex);
        default:
            throw "No predicate with index:" + ruleIndex;
       }
    }

    expression_sempred(localctx, predIndex) {
    	switch(predIndex) {
    		case 0:
    			return this.precpred(this._ctx, 13);
    		case 1:
    			return this.precpred(this._ctx, 12);
    		case 2:
    			return this.precpred(this._ctx, 11);
    		case 3:
    			return this.precpred(this._ctx, 10);
    		case 4:
    			return this.precpred(this._ctx, 9);
    		case 5:
    			return this.precpred(this._ctx, 8);
    		case 6:
    			return this.precpred(this._ctx, 7);
    		case 7:
    			return this.precpred(this._ctx, 6);
    		case 8:
    			return this.precpred(this._ctx, 5);
    		case 9:
    			return this.precpred(this._ctx, 4);
    		case 10:
    			return this.precpred(this._ctx, 3);
    		case 11:
    			return this.precpred(this._ctx, 2);
    		case 12:
    			return this.precpred(this._ctx, 1);
    		default:
    			throw "No predicate with index:" + predIndex;
    	}
    };

    reference_sempred(localctx, predIndex) {
    	switch(predIndex) {
    		case 13:
    			return this.precpred(this._ctx, 2);
    		case 14:
    			return this.precpred(this._ctx, 1);
    		default:
    			throw "No predicate with index:" + predIndex;
    	}
    };




	parse() {
	    let localctx = new ParseContext(this, this._ctx, this.state);
	    this.enterRule(localctx, 0, CAQLParser.RULE_parse);
	    try {
	        this.enterOuterAlt(localctx, 1);
	        this.state = 22;
	        this.expression(0);
	        this.state = 23;
	        this.match(CAQLParser.EOF);
	    } catch (re) {
	    	if(re instanceof antlr4.error.RecognitionException) {
		        localctx.exception = re;
		        this._errHandler.reportError(this, re);
		        this._errHandler.recover(this, re);
		    } else {
		    	throw re;
		    }
	    } finally {
	        this.exitRule();
	    }
	    return localctx;
	}


	expression(_p) {
		if(_p===undefined) {
		    _p = 0;
		}
	    const _parentctx = this._ctx;
	    const _parentState = this.state;
	    let localctx = new ExpressionContext(this, this._ctx, _parentState);
	    let _prevctx = localctx;
	    const _startState = 2;
	    this.enterRecursionRule(localctx, 2, CAQLParser.RULE_expression, _p);
	    var _la = 0; // Token type
	    try {
	        this.enterOuterAlt(localctx, 1);
	        this.state = 29;
	        this._errHandler.sync(this);
	        switch(this._input.LA(1)) {
	        case CAQLParser.T_FALSE:
	        case CAQLParser.T_NULL:
	        case CAQLParser.T_TRUE:
	        case CAQLParser.T_INT:
	        case CAQLParser.T_FLOAT:
	        case CAQLParser.T_QUOTED_STRING:
	            this.state = 26;
	            this.value_literal();
	            break;
	        case CAQLParser.T_OPEN:
	        case CAQLParser.T_OBJECT_OPEN:
	        case CAQLParser.T_ARRAY_OPEN:
	        case CAQLParser.T_STRING:
	            this.state = 27;
	            this.reference(0);
	            break;
	        case CAQLParser.T_PLUS:
	        case CAQLParser.T_MINUS:
	        case CAQLParser.T_NOT:
	            this.state = 28;
	            this.operator_unary();
	            break;
	        default:
	            throw new antlr4.error.NoViableAltException(this);
	        }
	        this._ctx.stop = this._input.LT(-1);
	        this.state = 85;
	        this._errHandler.sync(this);
	        var _alt = this._interp.adaptivePredict(this._input,4,this._ctx)
	        while(_alt!=2 && _alt!=antlr4.atn.ATN.INVALID_ALT_NUMBER) {
	            if(_alt===1) {
	                if(this._parseListeners!==null) {
	                    this.triggerExitRuleEvent();
	                }
	                _prevctx = localctx;
	                this.state = 83;
	                this._errHandler.sync(this);
	                var la_ = this._interp.adaptivePredict(this._input,3,this._ctx);
	                switch(la_) {
	                case 1:
	                    localctx = new ExpressionContext(this, _parentctx, _parentState);
	                    this.pushNewRecursionContext(localctx, _startState, CAQLParser.RULE_expression);
	                    this.state = 31;
	                    if (!( this.precpred(this._ctx, 13))) {
	                        throw new antlr4.error.FailedPredicateException(this, "this.precpred(this._ctx, 13)");
	                    }
	                    this.state = 32;
	                    _la = this._input.LA(1);
	                    if(!(_la===CAQLParser.T_PLUS || _la===CAQLParser.T_MINUS)) {
	                    this._errHandler.recoverInline(this);
	                    }
	                    else {
	                    	this._errHandler.reportMatch(this);
	                        this.consume();
	                    }
	                    this.state = 33;
	                    this.expression(14);
	                    break;

	                case 2:
	                    localctx = new ExpressionContext(this, _parentctx, _parentState);
	                    this.pushNewRecursionContext(localctx, _startState, CAQLParser.RULE_expression);
	                    this.state = 34;
	                    if (!( this.precpred(this._ctx, 12))) {
	                        throw new antlr4.error.FailedPredicateException(this, "this.precpred(this._ctx, 12)");
	                    }
	                    this.state = 35;
	                    _la = this._input.LA(1);
	                    if(!((((_la) & ~0x1f) == 0 && ((1 << _la) & ((1 << CAQLParser.T_TIMES) | (1 << CAQLParser.T_DIV) | (1 << CAQLParser.T_MOD))) !== 0))) {
	                    this._errHandler.recoverInline(this);
	                    }
	                    else {
	                    	this._errHandler.reportMatch(this);
	                        this.consume();
	                    }
	                    this.state = 36;
	                    this.expression(13);
	                    break;

	                case 3:
	                    localctx = new ExpressionContext(this, _parentctx, _parentState);
	                    this.pushNewRecursionContext(localctx, _startState, CAQLParser.RULE_expression);
	                    this.state = 37;
	                    if (!( this.precpred(this._ctx, 11))) {
	                        throw new antlr4.error.FailedPredicateException(this, "this.precpred(this._ctx, 11)");
	                    }
	                    this.state = 38;
	                    this.match(CAQLParser.T_RANGE);
	                    this.state = 39;
	                    this.expression(12);
	                    break;

	                case 4:
	                    localctx = new ExpressionContext(this, _parentctx, _parentState);
	                    this.pushNewRecursionContext(localctx, _startState, CAQLParser.RULE_expression);
	                    this.state = 40;
	                    if (!( this.precpred(this._ctx, 10))) {
	                        throw new antlr4.error.FailedPredicateException(this, "this.precpred(this._ctx, 10)");
	                    }
	                    this.state = 41;
	                    _la = this._input.LA(1);
	                    if(!((((_la) & ~0x1f) == 0 && ((1 << _la) & ((1 << CAQLParser.T_LT) | (1 << CAQLParser.T_GT) | (1 << CAQLParser.T_LE) | (1 << CAQLParser.T_GE))) !== 0))) {
	                    this._errHandler.recoverInline(this);
	                    }
	                    else {
	                    	this._errHandler.reportMatch(this);
	                        this.consume();
	                    }
	                    this.state = 42;
	                    this.expression(11);
	                    break;

	                case 5:
	                    localctx = new ExpressionContext(this, _parentctx, _parentState);
	                    this.pushNewRecursionContext(localctx, _startState, CAQLParser.RULE_expression);
	                    this.state = 43;
	                    if (!( this.precpred(this._ctx, 9))) {
	                        throw new antlr4.error.FailedPredicateException(this, "this.precpred(this._ctx, 9)");
	                    }
	                    this.state = 45;
	                    this._errHandler.sync(this);
	                    _la = this._input.LA(1);
	                    if(_la===CAQLParser.T_NOT) {
	                        this.state = 44;
	                        this.match(CAQLParser.T_NOT);
	                    }

	                    this.state = 47;
	                    this.match(CAQLParser.T_IN);
	                    this.state = 48;
	                    this.expression(10);
	                    break;

	                case 6:
	                    localctx = new ExpressionContext(this, _parentctx, _parentState);
	                    this.pushNewRecursionContext(localctx, _startState, CAQLParser.RULE_expression);
	                    this.state = 49;
	                    if (!( this.precpred(this._ctx, 8))) {
	                        throw new antlr4.error.FailedPredicateException(this, "this.precpred(this._ctx, 8)");
	                    }
	                    this.state = 50;
	                    _la = this._input.LA(1);
	                    if(!(_la===CAQLParser.T_EQ || _la===CAQLParser.T_NE)) {
	                    this._errHandler.recoverInline(this);
	                    }
	                    else {
	                    	this._errHandler.reportMatch(this);
	                        this.consume();
	                    }
	                    this.state = 51;
	                    this.expression(9);
	                    break;

	                case 7:
	                    localctx = new ExpressionContext(this, _parentctx, _parentState);
	                    this.pushNewRecursionContext(localctx, _startState, CAQLParser.RULE_expression);
	                    this.state = 52;
	                    if (!( this.precpred(this._ctx, 7))) {
	                        throw new antlr4.error.FailedPredicateException(this, "this.precpred(this._ctx, 7)");
	                    }
	                    this.state = 53;
	                    _la = this._input.LA(1);
	                    if(!(((((_la - 27)) & ~0x1f) == 0 && ((1 << (_la - 27)) & ((1 << (CAQLParser.T_ALL - 27)) | (1 << (CAQLParser.T_ANY - 27)) | (1 << (CAQLParser.T_NONE - 27)))) !== 0))) {
	                    this._errHandler.recoverInline(this);
	                    }
	                    else {
	                    	this._errHandler.reportMatch(this);
	                        this.consume();
	                    }
	                    this.state = 54;
	                    localctx.eq_op = this._input.LT(1);
	                    _la = this._input.LA(1);
	                    if(!((((_la) & ~0x1f) == 0 && ((1 << _la) & ((1 << CAQLParser.T_EQ) | (1 << CAQLParser.T_NE) | (1 << CAQLParser.T_LT) | (1 << CAQLParser.T_GT) | (1 << CAQLParser.T_LE) | (1 << CAQLParser.T_GE))) !== 0) || _la===CAQLParser.T_IN)) {
	                        localctx.eq_op = this._errHandler.recoverInline(this);
	                    }
	                    else {
	                    	this._errHandler.reportMatch(this);
	                        this.consume();
	                    }
	                    this.state = 55;
	                    this.expression(8);
	                    break;

	                case 8:
	                    localctx = new ExpressionContext(this, _parentctx, _parentState);
	                    this.pushNewRecursionContext(localctx, _startState, CAQLParser.RULE_expression);
	                    this.state = 56;
	                    if (!( this.precpred(this._ctx, 6))) {
	                        throw new antlr4.error.FailedPredicateException(this, "this.precpred(this._ctx, 6)");
	                    }
	                    this.state = 57;
	                    _la = this._input.LA(1);
	                    if(!(((((_la - 27)) & ~0x1f) == 0 && ((1 << (_la - 27)) & ((1 << (CAQLParser.T_ALL - 27)) | (1 << (CAQLParser.T_ANY - 27)) | (1 << (CAQLParser.T_NONE - 27)))) !== 0))) {
	                    this._errHandler.recoverInline(this);
	                    }
	                    else {
	                    	this._errHandler.reportMatch(this);
	                        this.consume();
	                    }
	                    this.state = 58;
	                    this.match(CAQLParser.T_NOT);
	                    this.state = 59;
	                    this.match(CAQLParser.T_IN);
	                    this.state = 60;
	                    this.expression(7);
	                    break;

	                case 9:
	                    localctx = new ExpressionContext(this, _parentctx, _parentState);
	                    this.pushNewRecursionContext(localctx, _startState, CAQLParser.RULE_expression);
	                    this.state = 61;
	                    if (!( this.precpred(this._ctx, 5))) {
	                        throw new antlr4.error.FailedPredicateException(this, "this.precpred(this._ctx, 5)");
	                    }
	                    this.state = 63;
	                    this._errHandler.sync(this);
	                    _la = this._input.LA(1);
	                    if(_la===CAQLParser.T_NOT) {
	                        this.state = 62;
	                        this.match(CAQLParser.T_NOT);
	                    }

	                    this.state = 65;
	                    _la = this._input.LA(1);
	                    if(!(_la===CAQLParser.T_REGEX_MATCH || _la===CAQLParser.T_REGEX_NON_MATCH || _la===CAQLParser.T_LIKE)) {
	                    this._errHandler.recoverInline(this);
	                    }
	                    else {
	                    	this._errHandler.reportMatch(this);
	                        this.consume();
	                    }
	                    this.state = 66;
	                    this.expression(6);
	                    break;

	                case 10:
	                    localctx = new ExpressionContext(this, _parentctx, _parentState);
	                    this.pushNewRecursionContext(localctx, _startState, CAQLParser.RULE_expression);
	                    this.state = 67;
	                    if (!( this.precpred(this._ctx, 4))) {
	                        throw new antlr4.error.FailedPredicateException(this, "this.precpred(this._ctx, 4)");
	                    }
	                    this.state = 68;
	                    this.match(CAQLParser.T_AND);
	                    this.state = 69;
	                    this.expression(5);
	                    break;

	                case 11:
	                    localctx = new ExpressionContext(this, _parentctx, _parentState);
	                    this.pushNewRecursionContext(localctx, _startState, CAQLParser.RULE_expression);
	                    this.state = 70;
	                    if (!( this.precpred(this._ctx, 3))) {
	                        throw new antlr4.error.FailedPredicateException(this, "this.precpred(this._ctx, 3)");
	                    }
	                    this.state = 71;
	                    this.match(CAQLParser.T_OR);
	                    this.state = 72;
	                    this.expression(4);
	                    break;

	                case 12:
	                    localctx = new ExpressionContext(this, _parentctx, _parentState);
	                    this.pushNewRecursionContext(localctx, _startState, CAQLParser.RULE_expression);
	                    this.state = 73;
	                    if (!( this.precpred(this._ctx, 2))) {
	                        throw new antlr4.error.FailedPredicateException(this, "this.precpred(this._ctx, 2)");
	                    }
	                    this.state = 74;
	                    this.match(CAQLParser.T_QUESTION);
	                    this.state = 75;
	                    this.expression(0);
	                    this.state = 76;
	                    this.match(CAQLParser.T_COLON);
	                    this.state = 77;
	                    this.expression(3);
	                    break;

	                case 13:
	                    localctx = new ExpressionContext(this, _parentctx, _parentState);
	                    this.pushNewRecursionContext(localctx, _startState, CAQLParser.RULE_expression);
	                    this.state = 79;
	                    if (!( this.precpred(this._ctx, 1))) {
	                        throw new antlr4.error.FailedPredicateException(this, "this.precpred(this._ctx, 1)");
	                    }
	                    this.state = 80;
	                    this.match(CAQLParser.T_QUESTION);
	                    this.state = 81;
	                    this.match(CAQLParser.T_COLON);
	                    this.state = 82;
	                    this.expression(2);
	                    break;

	                } 
	            }
	            this.state = 87;
	            this._errHandler.sync(this);
	            _alt = this._interp.adaptivePredict(this._input,4,this._ctx);
	        }

	    } catch( error) {
	        if(error instanceof antlr4.error.RecognitionException) {
		        localctx.exception = error;
		        this._errHandler.reportError(this, error);
		        this._errHandler.recover(this, error);
		    } else {
		    	throw error;
		    }
	    } finally {
	        this.unrollRecursionContexts(_parentctx)
	    }
	    return localctx;
	}



	operator_unary() {
	    let localctx = new Operator_unaryContext(this, this._ctx, this.state);
	    this.enterRule(localctx, 4, CAQLParser.RULE_operator_unary);
	    try {
	        this.enterOuterAlt(localctx, 1);
	        this.state = 94;
	        this._errHandler.sync(this);
	        switch(this._input.LA(1)) {
	        case CAQLParser.T_PLUS:
	            this.state = 88;
	            this.match(CAQLParser.T_PLUS);
	            this.state = 89;
	            this.expression(0);
	            break;
	        case CAQLParser.T_MINUS:
	            this.state = 90;
	            this.match(CAQLParser.T_MINUS);
	            this.state = 91;
	            this.expression(0);
	            break;
	        case CAQLParser.T_NOT:
	            this.state = 92;
	            this.match(CAQLParser.T_NOT);
	            this.state = 93;
	            this.expression(0);
	            break;
	        default:
	            throw new antlr4.error.NoViableAltException(this);
	        }
	    } catch (re) {
	    	if(re instanceof antlr4.error.RecognitionException) {
		        localctx.exception = re;
		        this._errHandler.reportError(this, re);
		        this._errHandler.recover(this, re);
		    } else {
		    	throw re;
		    }
	    } finally {
	        this.exitRule();
	    }
	    return localctx;
	}


	reference(_p) {
		if(_p===undefined) {
		    _p = 0;
		}
	    const _parentctx = this._ctx;
	    const _parentState = this.state;
	    let localctx = new ReferenceContext(this, this._ctx, _parentState);
	    let _prevctx = localctx;
	    const _startState = 6;
	    this.enterRecursionRule(localctx, 6, CAQLParser.RULE_reference, _p);
	    try {
	        this.enterOuterAlt(localctx, 1);
	        this.state = 104;
	        this._errHandler.sync(this);
	        var la_ = this._interp.adaptivePredict(this._input,6,this._ctx);
	        switch(la_) {
	        case 1:
	            this.state = 97;
	            this.match(CAQLParser.T_STRING);
	            break;

	        case 2:
	            this.state = 98;
	            this.compound_value();
	            break;

	        case 3:
	            this.state = 99;
	            this.function_call();
	            break;

	        case 4:
	            this.state = 100;
	            this.match(CAQLParser.T_OPEN);
	            this.state = 101;
	            this.expression(0);
	            this.state = 102;
	            this.match(CAQLParser.T_CLOSE);
	            break;

	        }
	        this._ctx.stop = this._input.LT(-1);
	        this.state = 116;
	        this._errHandler.sync(this);
	        var _alt = this._interp.adaptivePredict(this._input,8,this._ctx)
	        while(_alt!=2 && _alt!=antlr4.atn.ATN.INVALID_ALT_NUMBER) {
	            if(_alt===1) {
	                if(this._parseListeners!==null) {
	                    this.triggerExitRuleEvent();
	                }
	                _prevctx = localctx;
	                this.state = 114;
	                this._errHandler.sync(this);
	                var la_ = this._interp.adaptivePredict(this._input,7,this._ctx);
	                switch(la_) {
	                case 1:
	                    localctx = new ReferenceContext(this, _parentctx, _parentState);
	                    this.pushNewRecursionContext(localctx, _startState, CAQLParser.RULE_reference);
	                    this.state = 106;
	                    if (!( this.precpred(this._ctx, 2))) {
	                        throw new antlr4.error.FailedPredicateException(this, "this.precpred(this._ctx, 2)");
	                    }
	                    this.state = 107;
	                    this.match(CAQLParser.DOT);
	                    this.state = 108;
	                    this.match(CAQLParser.T_STRING);
	                    break;

	                case 2:
	                    localctx = new ReferenceContext(this, _parentctx, _parentState);
	                    this.pushNewRecursionContext(localctx, _startState, CAQLParser.RULE_reference);
	                    this.state = 109;
	                    if (!( this.precpred(this._ctx, 1))) {
	                        throw new antlr4.error.FailedPredicateException(this, "this.precpred(this._ctx, 1)");
	                    }
	                    this.state = 110;
	                    this.match(CAQLParser.T_ARRAY_OPEN);
	                    this.state = 111;
	                    this.expression(0);
	                    this.state = 112;
	                    this.match(CAQLParser.T_ARRAY_CLOSE);
	                    break;

	                } 
	            }
	            this.state = 118;
	            this._errHandler.sync(this);
	            _alt = this._interp.adaptivePredict(this._input,8,this._ctx);
	        }

	    } catch( error) {
	        if(error instanceof antlr4.error.RecognitionException) {
		        localctx.exception = error;
		        this._errHandler.reportError(this, error);
		        this._errHandler.recover(this, error);
		    } else {
		    	throw error;
		    }
	    } finally {
	        this.unrollRecursionContexts(_parentctx)
	    }
	    return localctx;
	}



	compound_value() {
	    let localctx = new Compound_valueContext(this, this._ctx, this.state);
	    this.enterRule(localctx, 8, CAQLParser.RULE_compound_value);
	    try {
	        this.enterOuterAlt(localctx, 1);
	        this.state = 121;
	        this._errHandler.sync(this);
	        switch(this._input.LA(1)) {
	        case CAQLParser.T_ARRAY_OPEN:
	            this.state = 119;
	            this.array();
	            break;
	        case CAQLParser.T_OBJECT_OPEN:
	            this.state = 120;
	            this.object();
	            break;
	        default:
	            throw new antlr4.error.NoViableAltException(this);
	        }
	    } catch (re) {
	    	if(re instanceof antlr4.error.RecognitionException) {
		        localctx.exception = re;
		        this._errHandler.reportError(this, re);
		        this._errHandler.recover(this, re);
		    } else {
		    	throw re;
		    }
	    } finally {
	        this.exitRule();
	    }
	    return localctx;
	}



	function_call() {
	    let localctx = new Function_callContext(this, this._ctx, this.state);
	    this.enterRule(localctx, 10, CAQLParser.RULE_function_call);
	    var _la = 0; // Token type
	    try {
	        this.enterOuterAlt(localctx, 1);
	        this.state = 123;
	        this.match(CAQLParser.T_STRING);
	        this.state = 124;
	        this.match(CAQLParser.T_OPEN);
	        this.state = 126;
	        this._errHandler.sync(this);
	        _la = this._input.LA(1);
	        if(((((_la - 10)) & ~0x1f) == 0 && ((1 << (_la - 10)) & ((1 << (CAQLParser.T_PLUS - 10)) | (1 << (CAQLParser.T_MINUS - 10)) | (1 << (CAQLParser.T_OPEN - 10)) | (1 << (CAQLParser.T_OBJECT_OPEN - 10)) | (1 << (CAQLParser.T_ARRAY_OPEN - 10)) | (1 << (CAQLParser.T_FALSE - 10)))) !== 0) || ((((_la - 47)) & ~0x1f) == 0 && ((1 << (_la - 47)) & ((1 << (CAQLParser.T_NOT - 47)) | (1 << (CAQLParser.T_NULL - 47)) | (1 << (CAQLParser.T_TRUE - 47)) | (1 << (CAQLParser.T_STRING - 47)) | (1 << (CAQLParser.T_INT - 47)) | (1 << (CAQLParser.T_FLOAT - 47)) | (1 << (CAQLParser.T_QUOTED_STRING - 47)))) !== 0)) {
	            this.state = 125;
	            this.expression(0);
	        }

	        this.state = 132;
	        this._errHandler.sync(this);
	        var _alt = this._interp.adaptivePredict(this._input,11,this._ctx)
	        while(_alt!=1 && _alt!=antlr4.atn.ATN.INVALID_ALT_NUMBER) {
	            if(_alt===1+1) {
	                this.state = 128;
	                this.match(CAQLParser.T_COMMA);
	                this.state = 129;
	                this.expression(0); 
	            }
	            this.state = 134;
	            this._errHandler.sync(this);
	            _alt = this._interp.adaptivePredict(this._input,11,this._ctx);
	        }

	        this.state = 136;
	        this._errHandler.sync(this);
	        _la = this._input.LA(1);
	        if(_la===CAQLParser.T_COMMA) {
	            this.state = 135;
	            this.match(CAQLParser.T_COMMA);
	        }

	        this.state = 138;
	        this.match(CAQLParser.T_CLOSE);
	    } catch (re) {
	    	if(re instanceof antlr4.error.RecognitionException) {
		        localctx.exception = re;
		        this._errHandler.reportError(this, re);
		        this._errHandler.recover(this, re);
		    } else {
		    	throw re;
		    }
	    } finally {
	        this.exitRule();
	    }
	    return localctx;
	}



	value_literal() {
	    let localctx = new Value_literalContext(this, this._ctx, this.state);
	    this.enterRule(localctx, 12, CAQLParser.RULE_value_literal);
	    var _la = 0; // Token type
	    try {
	        this.enterOuterAlt(localctx, 1);
	        this.state = 140;
	        _la = this._input.LA(1);
	        if(!(((((_la - 34)) & ~0x1f) == 0 && ((1 << (_la - 34)) & ((1 << (CAQLParser.T_FALSE - 34)) | (1 << (CAQLParser.T_NULL - 34)) | (1 << (CAQLParser.T_TRUE - 34)))) !== 0) || ((((_la - 70)) & ~0x1f) == 0 && ((1 << (_la - 70)) & ((1 << (CAQLParser.T_INT - 70)) | (1 << (CAQLParser.T_FLOAT - 70)) | (1 << (CAQLParser.T_QUOTED_STRING - 70)))) !== 0))) {
	        this._errHandler.recoverInline(this);
	        }
	        else {
	        	this._errHandler.reportMatch(this);
	            this.consume();
	        }
	    } catch (re) {
	    	if(re instanceof antlr4.error.RecognitionException) {
		        localctx.exception = re;
		        this._errHandler.reportError(this, re);
		        this._errHandler.recover(this, re);
		    } else {
		    	throw re;
		    }
	    } finally {
	        this.exitRule();
	    }
	    return localctx;
	}



	array() {
	    let localctx = new ArrayContext(this, this._ctx, this.state);
	    this.enterRule(localctx, 14, CAQLParser.RULE_array);
	    var _la = 0; // Token type
	    try {
	        this.enterOuterAlt(localctx, 1);
	        this.state = 142;
	        this.match(CAQLParser.T_ARRAY_OPEN);
	        this.state = 144;
	        this._errHandler.sync(this);
	        _la = this._input.LA(1);
	        if(((((_la - 10)) & ~0x1f) == 0 && ((1 << (_la - 10)) & ((1 << (CAQLParser.T_PLUS - 10)) | (1 << (CAQLParser.T_MINUS - 10)) | (1 << (CAQLParser.T_OPEN - 10)) | (1 << (CAQLParser.T_OBJECT_OPEN - 10)) | (1 << (CAQLParser.T_ARRAY_OPEN - 10)) | (1 << (CAQLParser.T_FALSE - 10)))) !== 0) || ((((_la - 47)) & ~0x1f) == 0 && ((1 << (_la - 47)) & ((1 << (CAQLParser.T_NOT - 47)) | (1 << (CAQLParser.T_NULL - 47)) | (1 << (CAQLParser.T_TRUE - 47)) | (1 << (CAQLParser.T_STRING - 47)) | (1 << (CAQLParser.T_INT - 47)) | (1 << (CAQLParser.T_FLOAT - 47)) | (1 << (CAQLParser.T_QUOTED_STRING - 47)))) !== 0)) {
	            this.state = 143;
	            this.expression(0);
	        }

	        this.state = 150;
	        this._errHandler.sync(this);
	        var _alt = this._interp.adaptivePredict(this._input,14,this._ctx)
	        while(_alt!=1 && _alt!=antlr4.atn.ATN.INVALID_ALT_NUMBER) {
	            if(_alt===1+1) {
	                this.state = 146;
	                this.match(CAQLParser.T_COMMA);
	                this.state = 147;
	                this.expression(0); 
	            }
	            this.state = 152;
	            this._errHandler.sync(this);
	            _alt = this._interp.adaptivePredict(this._input,14,this._ctx);
	        }

	        this.state = 154;
	        this._errHandler.sync(this);
	        _la = this._input.LA(1);
	        if(_la===CAQLParser.T_COMMA) {
	            this.state = 153;
	            this.match(CAQLParser.T_COMMA);
	        }

	        this.state = 156;
	        this.match(CAQLParser.T_ARRAY_CLOSE);
	    } catch (re) {
	    	if(re instanceof antlr4.error.RecognitionException) {
		        localctx.exception = re;
		        this._errHandler.reportError(this, re);
		        this._errHandler.recover(this, re);
		    } else {
		    	throw re;
		    }
	    } finally {
	        this.exitRule();
	    }
	    return localctx;
	}



	object() {
	    let localctx = new ObjectContext(this, this._ctx, this.state);
	    this.enterRule(localctx, 16, CAQLParser.RULE_object);
	    var _la = 0; // Token type
	    try {
	        this.enterOuterAlt(localctx, 1);
	        this.state = 158;
	        this.match(CAQLParser.T_OBJECT_OPEN);
	        this.state = 160;
	        this._errHandler.sync(this);
	        _la = this._input.LA(1);
	        if(_la===CAQLParser.T_ARRAY_OPEN || _la===CAQLParser.T_STRING || _la===CAQLParser.T_QUOTED_STRING) {
	            this.state = 159;
	            this.object_element();
	        }

	        this.state = 166;
	        this._errHandler.sync(this);
	        var _alt = this._interp.adaptivePredict(this._input,17,this._ctx)
	        while(_alt!=2 && _alt!=antlr4.atn.ATN.INVALID_ALT_NUMBER) {
	            if(_alt===1) {
	                this.state = 162;
	                this.match(CAQLParser.T_COMMA);
	                this.state = 163;
	                this.object_element(); 
	            }
	            this.state = 168;
	            this._errHandler.sync(this);
	            _alt = this._interp.adaptivePredict(this._input,17,this._ctx);
	        }

	        this.state = 170;
	        this._errHandler.sync(this);
	        _la = this._input.LA(1);
	        if(_la===CAQLParser.T_COMMA) {
	            this.state = 169;
	            this.match(CAQLParser.T_COMMA);
	        }

	        this.state = 172;
	        this.match(CAQLParser.T_OBJECT_CLOSE);
	    } catch (re) {
	    	if(re instanceof antlr4.error.RecognitionException) {
		        localctx.exception = re;
		        this._errHandler.reportError(this, re);
		        this._errHandler.recover(this, re);
		    } else {
		    	throw re;
		    }
	    } finally {
	        this.exitRule();
	    }
	    return localctx;
	}



	object_element() {
	    let localctx = new Object_elementContext(this, this._ctx, this.state);
	    this.enterRule(localctx, 18, CAQLParser.RULE_object_element);
	    try {
	        this.enterOuterAlt(localctx, 1);
	        this.state = 185;
	        this._errHandler.sync(this);
	        var la_ = this._interp.adaptivePredict(this._input,19,this._ctx);
	        switch(la_) {
	        case 1:
	            this.state = 174;
	            this.match(CAQLParser.T_STRING);
	            break;

	        case 2:
	            this.state = 175;
	            this.object_element_name();
	            this.state = 176;
	            this.match(CAQLParser.T_COLON);
	            this.state = 177;
	            this.expression(0);
	            break;

	        case 3:
	            this.state = 179;
	            this.match(CAQLParser.T_ARRAY_OPEN);
	            this.state = 180;
	            this.expression(0);
	            this.state = 181;
	            this.match(CAQLParser.T_ARRAY_CLOSE);
	            this.state = 182;
	            this.match(CAQLParser.T_COLON);
	            this.state = 183;
	            this.expression(0);
	            break;

	        }
	    } catch (re) {
	    	if(re instanceof antlr4.error.RecognitionException) {
		        localctx.exception = re;
		        this._errHandler.reportError(this, re);
		        this._errHandler.recover(this, re);
		    } else {
		    	throw re;
		    }
	    } finally {
	        this.exitRule();
	    }
	    return localctx;
	}



	object_element_name() {
	    let localctx = new Object_element_nameContext(this, this._ctx, this.state);
	    this.enterRule(localctx, 20, CAQLParser.RULE_object_element_name);
	    var _la = 0; // Token type
	    try {
	        this.enterOuterAlt(localctx, 1);
	        this.state = 187;
	        _la = this._input.LA(1);
	        if(!(_la===CAQLParser.T_STRING || _la===CAQLParser.T_QUOTED_STRING)) {
	        this._errHandler.recoverInline(this);
	        }
	        else {
	        	this._errHandler.reportMatch(this);
	            this.consume();
	        }
	    } catch (re) {
	    	if(re instanceof antlr4.error.RecognitionException) {
		        localctx.exception = re;
		        this._errHandler.reportError(this, re);
		        this._errHandler.recover(this, re);
		    } else {
		    	throw re;
		    }
	    } finally {
	        this.exitRule();
	    }
	    return localctx;
	}


}

CAQLParser.EOF = antlr4.Token.EOF;
CAQLParser.DOT = 1;
CAQLParser.T_REGEX_MATCH = 2;
CAQLParser.T_REGEX_NON_MATCH = 3;
CAQLParser.T_EQ = 4;
CAQLParser.T_NE = 5;
CAQLParser.T_LT = 6;
CAQLParser.T_GT = 7;
CAQLParser.T_LE = 8;
CAQLParser.T_GE = 9;
CAQLParser.T_PLUS = 10;
CAQLParser.T_MINUS = 11;
CAQLParser.T_TIMES = 12;
CAQLParser.T_DIV = 13;
CAQLParser.T_MOD = 14;
CAQLParser.T_QUESTION = 15;
CAQLParser.T_COLON = 16;
CAQLParser.T_SCOPE = 17;
CAQLParser.T_RANGE = 18;
CAQLParser.T_COMMA = 19;
CAQLParser.T_OPEN = 20;
CAQLParser.T_CLOSE = 21;
CAQLParser.T_OBJECT_OPEN = 22;
CAQLParser.T_OBJECT_CLOSE = 23;
CAQLParser.T_ARRAY_OPEN = 24;
CAQLParser.T_ARRAY_CLOSE = 25;
CAQLParser.T_AGGREGATE = 26;
CAQLParser.T_ALL = 27;
CAQLParser.T_AND = 28;
CAQLParser.T_ANY = 29;
CAQLParser.T_ASC = 30;
CAQLParser.T_COLLECT = 31;
CAQLParser.T_DESC = 32;
CAQLParser.T_DISTINCT = 33;
CAQLParser.T_FALSE = 34;
CAQLParser.T_FILTER = 35;
CAQLParser.T_FOR = 36;
CAQLParser.T_GRAPH = 37;
CAQLParser.T_IN = 38;
CAQLParser.T_INBOUND = 39;
CAQLParser.T_INSERT = 40;
CAQLParser.T_INTO = 41;
CAQLParser.T_K_SHORTEST_PATHS = 42;
CAQLParser.T_LET = 43;
CAQLParser.T_LIKE = 44;
CAQLParser.T_LIMIT = 45;
CAQLParser.T_NONE = 46;
CAQLParser.T_NOT = 47;
CAQLParser.T_NULL = 48;
CAQLParser.T_OR = 49;
CAQLParser.T_OUTBOUND = 50;
CAQLParser.T_REMOVE = 51;
CAQLParser.T_REPLACE = 52;
CAQLParser.T_RETURN = 53;
CAQLParser.T_SHORTEST_PATH = 54;
CAQLParser.T_SORT = 55;
CAQLParser.T_TRUE = 56;
CAQLParser.T_UPDATE = 57;
CAQLParser.T_UPSERT = 58;
CAQLParser.T_WITH = 59;
CAQLParser.T_KEEP = 60;
CAQLParser.T_COUNT = 61;
CAQLParser.T_OPTIONS = 62;
CAQLParser.T_PRUNE = 63;
CAQLParser.T_SEARCH = 64;
CAQLParser.T_TO = 65;
CAQLParser.T_CURRENT = 66;
CAQLParser.T_NEW = 67;
CAQLParser.T_OLD = 68;
CAQLParser.T_STRING = 69;
CAQLParser.T_INT = 70;
CAQLParser.T_FLOAT = 71;
CAQLParser.T_PARAMETER = 72;
CAQLParser.T_QUOTED_STRING = 73;
CAQLParser.SINGLE_LINE_COMMENT = 74;
CAQLParser.MULTILINE_COMMENT = 75;
CAQLParser.SPACES = 76;
CAQLParser.UNEXPECTED_CHAR = 77;
CAQLParser.ERROR_RECONGNIGION = 78;

CAQLParser.RULE_parse = 0;
CAQLParser.RULE_expression = 1;
CAQLParser.RULE_operator_unary = 2;
CAQLParser.RULE_reference = 3;
CAQLParser.RULE_compound_value = 4;
CAQLParser.RULE_function_call = 5;
CAQLParser.RULE_value_literal = 6;
CAQLParser.RULE_array = 7;
CAQLParser.RULE_object = 8;
CAQLParser.RULE_object_element = 9;
CAQLParser.RULE_object_element_name = 10;

class ParseContext extends antlr4.ParserRuleContext {

    constructor(parser, parent, invokingState) {
        if(parent===undefined) {
            parent = null;
        }
        if(invokingState===undefined || invokingState===null) {
            invokingState = -1;
        }
        super(parent, invokingState);
        this.parser = parser;
        this.ruleIndex = CAQLParser.RULE_parse;
    }

	expression() {
	    return this.getTypedRuleContext(ExpressionContext,0);
	};

	EOF() {
	    return this.getToken(CAQLParser.EOF, 0);
	};

	enterRule(listener) {
	    if(listener instanceof CAQLParserListener ) {
	        listener.enterParse(this);
		}
	}

	exitRule(listener) {
	    if(listener instanceof CAQLParserListener ) {
	        listener.exitParse(this);
		}
	}


}



class ExpressionContext extends antlr4.ParserRuleContext {

    constructor(parser, parent, invokingState) {
        if(parent===undefined) {
            parent = null;
        }
        if(invokingState===undefined || invokingState===null) {
            invokingState = -1;
        }
        super(parent, invokingState);
        this.parser = parser;
        this.ruleIndex = CAQLParser.RULE_expression;
        this.eq_op = null; // Token
    }

	value_literal() {
	    return this.getTypedRuleContext(Value_literalContext,0);
	};

	reference() {
	    return this.getTypedRuleContext(ReferenceContext,0);
	};

	operator_unary() {
	    return this.getTypedRuleContext(Operator_unaryContext,0);
	};

	expression = function(i) {
	    if(i===undefined) {
	        i = null;
	    }
	    if(i===null) {
	        return this.getTypedRuleContexts(ExpressionContext);
	    } else {
	        return this.getTypedRuleContext(ExpressionContext,i);
	    }
	};

	T_PLUS() {
	    return this.getToken(CAQLParser.T_PLUS, 0);
	};

	T_MINUS() {
	    return this.getToken(CAQLParser.T_MINUS, 0);
	};

	T_TIMES() {
	    return this.getToken(CAQLParser.T_TIMES, 0);
	};

	T_DIV() {
	    return this.getToken(CAQLParser.T_DIV, 0);
	};

	T_MOD() {
	    return this.getToken(CAQLParser.T_MOD, 0);
	};

	T_RANGE() {
	    return this.getToken(CAQLParser.T_RANGE, 0);
	};

	T_LT() {
	    return this.getToken(CAQLParser.T_LT, 0);
	};

	T_GT() {
	    return this.getToken(CAQLParser.T_GT, 0);
	};

	T_LE() {
	    return this.getToken(CAQLParser.T_LE, 0);
	};

	T_GE() {
	    return this.getToken(CAQLParser.T_GE, 0);
	};

	T_IN() {
	    return this.getToken(CAQLParser.T_IN, 0);
	};

	T_NOT() {
	    return this.getToken(CAQLParser.T_NOT, 0);
	};

	T_EQ() {
	    return this.getToken(CAQLParser.T_EQ, 0);
	};

	T_NE() {
	    return this.getToken(CAQLParser.T_NE, 0);
	};

	T_ALL() {
	    return this.getToken(CAQLParser.T_ALL, 0);
	};

	T_ANY() {
	    return this.getToken(CAQLParser.T_ANY, 0);
	};

	T_NONE() {
	    return this.getToken(CAQLParser.T_NONE, 0);
	};

	T_LIKE() {
	    return this.getToken(CAQLParser.T_LIKE, 0);
	};

	T_REGEX_MATCH() {
	    return this.getToken(CAQLParser.T_REGEX_MATCH, 0);
	};

	T_REGEX_NON_MATCH() {
	    return this.getToken(CAQLParser.T_REGEX_NON_MATCH, 0);
	};

	T_AND() {
	    return this.getToken(CAQLParser.T_AND, 0);
	};

	T_OR() {
	    return this.getToken(CAQLParser.T_OR, 0);
	};

	T_QUESTION() {
	    return this.getToken(CAQLParser.T_QUESTION, 0);
	};

	T_COLON() {
	    return this.getToken(CAQLParser.T_COLON, 0);
	};

	enterRule(listener) {
	    if(listener instanceof CAQLParserListener ) {
	        listener.enterExpression(this);
		}
	}

	exitRule(listener) {
	    if(listener instanceof CAQLParserListener ) {
	        listener.exitExpression(this);
		}
	}


}



class Operator_unaryContext extends antlr4.ParserRuleContext {

    constructor(parser, parent, invokingState) {
        if(parent===undefined) {
            parent = null;
        }
        if(invokingState===undefined || invokingState===null) {
            invokingState = -1;
        }
        super(parent, invokingState);
        this.parser = parser;
        this.ruleIndex = CAQLParser.RULE_operator_unary;
    }

	T_PLUS() {
	    return this.getToken(CAQLParser.T_PLUS, 0);
	};

	expression() {
	    return this.getTypedRuleContext(ExpressionContext,0);
	};

	T_MINUS() {
	    return this.getToken(CAQLParser.T_MINUS, 0);
	};

	T_NOT() {
	    return this.getToken(CAQLParser.T_NOT, 0);
	};

	enterRule(listener) {
	    if(listener instanceof CAQLParserListener ) {
	        listener.enterOperator_unary(this);
		}
	}

	exitRule(listener) {
	    if(listener instanceof CAQLParserListener ) {
	        listener.exitOperator_unary(this);
		}
	}


}



class ReferenceContext extends antlr4.ParserRuleContext {

    constructor(parser, parent, invokingState) {
        if(parent===undefined) {
            parent = null;
        }
        if(invokingState===undefined || invokingState===null) {
            invokingState = -1;
        }
        super(parent, invokingState);
        this.parser = parser;
        this.ruleIndex = CAQLParser.RULE_reference;
    }

	T_STRING() {
	    return this.getToken(CAQLParser.T_STRING, 0);
	};

	compound_value() {
	    return this.getTypedRuleContext(Compound_valueContext,0);
	};

	function_call() {
	    return this.getTypedRuleContext(Function_callContext,0);
	};

	T_OPEN() {
	    return this.getToken(CAQLParser.T_OPEN, 0);
	};

	expression() {
	    return this.getTypedRuleContext(ExpressionContext,0);
	};

	T_CLOSE() {
	    return this.getToken(CAQLParser.T_CLOSE, 0);
	};

	reference() {
	    return this.getTypedRuleContext(ReferenceContext,0);
	};

	DOT() {
	    return this.getToken(CAQLParser.DOT, 0);
	};

	T_ARRAY_OPEN() {
	    return this.getToken(CAQLParser.T_ARRAY_OPEN, 0);
	};

	T_ARRAY_CLOSE() {
	    return this.getToken(CAQLParser.T_ARRAY_CLOSE, 0);
	};

	enterRule(listener) {
	    if(listener instanceof CAQLParserListener ) {
	        listener.enterReference(this);
		}
	}

	exitRule(listener) {
	    if(listener instanceof CAQLParserListener ) {
	        listener.exitReference(this);
		}
	}


}



class Compound_valueContext extends antlr4.ParserRuleContext {

    constructor(parser, parent, invokingState) {
        if(parent===undefined) {
            parent = null;
        }
        if(invokingState===undefined || invokingState===null) {
            invokingState = -1;
        }
        super(parent, invokingState);
        this.parser = parser;
        this.ruleIndex = CAQLParser.RULE_compound_value;
    }

	array() {
	    return this.getTypedRuleContext(ArrayContext,0);
	};

	object() {
	    return this.getTypedRuleContext(ObjectContext,0);
	};

	enterRule(listener) {
	    if(listener instanceof CAQLParserListener ) {
	        listener.enterCompound_value(this);
		}
	}

	exitRule(listener) {
	    if(listener instanceof CAQLParserListener ) {
	        listener.exitCompound_value(this);
		}
	}


}



class Function_callContext extends antlr4.ParserRuleContext {

    constructor(parser, parent, invokingState) {
        if(parent===undefined) {
            parent = null;
        }
        if(invokingState===undefined || invokingState===null) {
            invokingState = -1;
        }
        super(parent, invokingState);
        this.parser = parser;
        this.ruleIndex = CAQLParser.RULE_function_call;
    }

	T_STRING() {
	    return this.getToken(CAQLParser.T_STRING, 0);
	};

	T_OPEN() {
	    return this.getToken(CAQLParser.T_OPEN, 0);
	};

	T_CLOSE() {
	    return this.getToken(CAQLParser.T_CLOSE, 0);
	};

	expression = function(i) {
	    if(i===undefined) {
	        i = null;
	    }
	    if(i===null) {
	        return this.getTypedRuleContexts(ExpressionContext);
	    } else {
	        return this.getTypedRuleContext(ExpressionContext,i);
	    }
	};

	T_COMMA = function(i) {
		if(i===undefined) {
			i = null;
		}
	    if(i===null) {
	        return this.getTokens(CAQLParser.T_COMMA);
	    } else {
	        return this.getToken(CAQLParser.T_COMMA, i);
	    }
	};


	enterRule(listener) {
	    if(listener instanceof CAQLParserListener ) {
	        listener.enterFunction_call(this);
		}
	}

	exitRule(listener) {
	    if(listener instanceof CAQLParserListener ) {
	        listener.exitFunction_call(this);
		}
	}


}



class Value_literalContext extends antlr4.ParserRuleContext {

    constructor(parser, parent, invokingState) {
        if(parent===undefined) {
            parent = null;
        }
        if(invokingState===undefined || invokingState===null) {
            invokingState = -1;
        }
        super(parent, invokingState);
        this.parser = parser;
        this.ruleIndex = CAQLParser.RULE_value_literal;
    }

	T_QUOTED_STRING() {
	    return this.getToken(CAQLParser.T_QUOTED_STRING, 0);
	};

	T_INT() {
	    return this.getToken(CAQLParser.T_INT, 0);
	};

	T_FLOAT() {
	    return this.getToken(CAQLParser.T_FLOAT, 0);
	};

	T_NULL() {
	    return this.getToken(CAQLParser.T_NULL, 0);
	};

	T_TRUE() {
	    return this.getToken(CAQLParser.T_TRUE, 0);
	};

	T_FALSE() {
	    return this.getToken(CAQLParser.T_FALSE, 0);
	};

	enterRule(listener) {
	    if(listener instanceof CAQLParserListener ) {
	        listener.enterValue_literal(this);
		}
	}

	exitRule(listener) {
	    if(listener instanceof CAQLParserListener ) {
	        listener.exitValue_literal(this);
		}
	}


}



class ArrayContext extends antlr4.ParserRuleContext {

    constructor(parser, parent, invokingState) {
        if(parent===undefined) {
            parent = null;
        }
        if(invokingState===undefined || invokingState===null) {
            invokingState = -1;
        }
        super(parent, invokingState);
        this.parser = parser;
        this.ruleIndex = CAQLParser.RULE_array;
    }

	T_ARRAY_OPEN() {
	    return this.getToken(CAQLParser.T_ARRAY_OPEN, 0);
	};

	T_ARRAY_CLOSE() {
	    return this.getToken(CAQLParser.T_ARRAY_CLOSE, 0);
	};

	expression = function(i) {
	    if(i===undefined) {
	        i = null;
	    }
	    if(i===null) {
	        return this.getTypedRuleContexts(ExpressionContext);
	    } else {
	        return this.getTypedRuleContext(ExpressionContext,i);
	    }
	};

	T_COMMA = function(i) {
		if(i===undefined) {
			i = null;
		}
	    if(i===null) {
	        return this.getTokens(CAQLParser.T_COMMA);
	    } else {
	        return this.getToken(CAQLParser.T_COMMA, i);
	    }
	};


	enterRule(listener) {
	    if(listener instanceof CAQLParserListener ) {
	        listener.enterArray(this);
		}
	}

	exitRule(listener) {
	    if(listener instanceof CAQLParserListener ) {
	        listener.exitArray(this);
		}
	}


}



class ObjectContext extends antlr4.ParserRuleContext {

    constructor(parser, parent, invokingState) {
        if(parent===undefined) {
            parent = null;
        }
        if(invokingState===undefined || invokingState===null) {
            invokingState = -1;
        }
        super(parent, invokingState);
        this.parser = parser;
        this.ruleIndex = CAQLParser.RULE_object;
    }

	T_OBJECT_OPEN() {
	    return this.getToken(CAQLParser.T_OBJECT_OPEN, 0);
	};

	T_OBJECT_CLOSE() {
	    return this.getToken(CAQLParser.T_OBJECT_CLOSE, 0);
	};

	object_element = function(i) {
	    if(i===undefined) {
	        i = null;
	    }
	    if(i===null) {
	        return this.getTypedRuleContexts(Object_elementContext);
	    } else {
	        return this.getTypedRuleContext(Object_elementContext,i);
	    }
	};

	T_COMMA = function(i) {
		if(i===undefined) {
			i = null;
		}
	    if(i===null) {
	        return this.getTokens(CAQLParser.T_COMMA);
	    } else {
	        return this.getToken(CAQLParser.T_COMMA, i);
	    }
	};


	enterRule(listener) {
	    if(listener instanceof CAQLParserListener ) {
	        listener.enterObject(this);
		}
	}

	exitRule(listener) {
	    if(listener instanceof CAQLParserListener ) {
	        listener.exitObject(this);
		}
	}


}



class Object_elementContext extends antlr4.ParserRuleContext {

    constructor(parser, parent, invokingState) {
        if(parent===undefined) {
            parent = null;
        }
        if(invokingState===undefined || invokingState===null) {
            invokingState = -1;
        }
        super(parent, invokingState);
        this.parser = parser;
        this.ruleIndex = CAQLParser.RULE_object_element;
    }

	T_STRING() {
	    return this.getToken(CAQLParser.T_STRING, 0);
	};

	object_element_name() {
	    return this.getTypedRuleContext(Object_element_nameContext,0);
	};

	T_COLON() {
	    return this.getToken(CAQLParser.T_COLON, 0);
	};

	expression = function(i) {
	    if(i===undefined) {
	        i = null;
	    }
	    if(i===null) {
	        return this.getTypedRuleContexts(ExpressionContext);
	    } else {
	        return this.getTypedRuleContext(ExpressionContext,i);
	    }
	};

	T_ARRAY_OPEN() {
	    return this.getToken(CAQLParser.T_ARRAY_OPEN, 0);
	};

	T_ARRAY_CLOSE() {
	    return this.getToken(CAQLParser.T_ARRAY_CLOSE, 0);
	};

	enterRule(listener) {
	    if(listener instanceof CAQLParserListener ) {
	        listener.enterObject_element(this);
		}
	}

	exitRule(listener) {
	    if(listener instanceof CAQLParserListener ) {
	        listener.exitObject_element(this);
		}
	}


}



class Object_element_nameContext extends antlr4.ParserRuleContext {

    constructor(parser, parent, invokingState) {
        if(parent===undefined) {
            parent = null;
        }
        if(invokingState===undefined || invokingState===null) {
            invokingState = -1;
        }
        super(parent, invokingState);
        this.parser = parser;
        this.ruleIndex = CAQLParser.RULE_object_element_name;
    }

	T_STRING() {
	    return this.getToken(CAQLParser.T_STRING, 0);
	};

	T_QUOTED_STRING() {
	    return this.getToken(CAQLParser.T_QUOTED_STRING, 0);
	};

	enterRule(listener) {
	    if(listener instanceof CAQLParserListener ) {
	        listener.enterObject_element_name(this);
		}
	}

	exitRule(listener) {
	    if(listener instanceof CAQLParserListener ) {
	        listener.exitObject_element_name(this);
		}
	}


}




CAQLParser.ParseContext = ParseContext; 
CAQLParser.ExpressionContext = ExpressionContext; 
CAQLParser.Operator_unaryContext = Operator_unaryContext; 
CAQLParser.ReferenceContext = ReferenceContext; 
CAQLParser.Compound_valueContext = Compound_valueContext; 
CAQLParser.Function_callContext = Function_callContext; 
CAQLParser.Value_literalContext = Value_literalContext; 
CAQLParser.ArrayContext = ArrayContext; 
CAQLParser.ObjectContext = ObjectContext; 
CAQLParser.Object_elementContext = Object_elementContext; 
CAQLParser.Object_element_nameContext = Object_element_nameContext; 
