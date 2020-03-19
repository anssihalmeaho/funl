package funl

import (
	"fmt"
	"testing"
)

func lextester(t *testing.T, text string, ttokens []token) {
	lexer := newTokenizer(NewDefaultOperators())
	tokens, err := lexer.scan(text)
	if err != nil {
		t.Fatalf("error : %v", err)
	}
	if l1, l2 := len(tokens), len(ttokens); l1 != l2 {
		t.Log(tokens)
		t.Fatalf("wrong (%v)(%v)", l1, l2)
		return
	}
	for i, token := range ttokens {
		if tokens[i].Type != token.Type {
			t.Fatalf("Unexpected type (%d): %v, %v (%s)(%s)", i, tokens[i].Type, token.Type, tokens[i].Value, token.Value)
		}
		if tokens[i].Value != token.Value {
			t.Fatalf("Unexpected value (%d): %s, %s", i, tokens[i].Value, token.Value)
		}
	}
}

func TestJustSymbols(t *testing.T) {
	text := "2010"
	tokens := []token{
		token{
			Type:  tokenNumber,
			Value: "2010",
		},
	}
	lextester(t, text, tokens)
	text = "dum_name"
	tokens = []token{
		token{
			Type:  tokenSymbol,
			Value: "dum_name",
		},
	}
	lextester(t, text, tokens)
}

func TestEmptyInput(t *testing.T) {
	lexer := newTokenizer(NewDefaultOperators())
	tokens, err := lexer.scan("")
	if err != nil {
		t.Fatalf("error : %v", err)
	}
	if l := len(tokens); l != 0 {
		t.Fatalf("should be empty (%d)", l)
	}
}

func TestNonAllowedChar(t *testing.T) {
	lexer := newTokenizer(NewDefaultOperators())
	_, err := lexer.scan("eq(dum.xyz, ? 'any text')")
	if err == nil {
		t.Fatalf("should fail")
	}
}

func TestLexBasic(t *testing.T) {
	text := "not(dum.sub, 'aabbcc text'any_more'what text')"
	tokens := []token{
		token{
			Type:  tokenOperator,
			Value: "not",
		},
		token{
			Type:  tokenOpenBracket,
			Value: "(",
		},
		token{
			Type:  tokenSymbol,
			Value: "dum",
		},
		token{
			Type:  tokenDot,
			Value: ".",
		},
		token{
			Type:  tokenSymbol,
			Value: "sub",
		},
		token{
			Type:  tokenComma,
			Value: ",",
		},
		token{
			Type:  tokenString,
			Value: "aabbcc text",
		},
		token{
			Type:  tokenSymbol,
			Value: "any_more",
		},
		token{
			Type:  tokenString,
			Value: "what text",
		},
		token{
			Type:  tokenClosingBracket,
			Value: ")",
		},
	}
	lextester(t, text, tokens)
}

func TestLexNumbers(t *testing.T) {
	text := "not(2007, 'dum, 2010 text', value1, value2)"
	tokens := []token{
		token{
			Type:  tokenOperator,
			Value: "not",
		},
		token{
			Type:  tokenOpenBracket,
			Value: "(",
		},
		token{
			Type:  tokenNumber,
			Value: "2007",
		},
		token{
			Type:  tokenComma,
			Value: ",",
		},
		token{
			Type:  tokenString,
			Value: "dum, 2010 text",
		},
		token{
			Type:  tokenComma,
			Value: ",",
		},
		token{
			Type:  tokenSymbol,
			Value: "value1",
		},
		token{
			Type:  tokenComma,
			Value: ",",
		},
		token{
			Type:  tokenSymbol,
			Value: "value2",
		},
		token{
			Type:  tokenClosingBracket,
			Value: ")",
		},
	}
	lextester(t, text, tokens)
}

func TestVariousTokens(t *testing.T) {
	text := "  gt  ( or(.dum.xyz.sub, or( '  any  (( ) and not, textTEXT ' , 2010, xy_z-201)) , not('xxx aa__-c--  zz ')  "
	tokens := []token{
		token{
			Type:  tokenOperator,
			Value: "gt",
		},
		token{
			Type:  tokenOpenBracket,
			Value: "(",
		},
		token{
			Type:  tokenOperator,
			Value: "or",
		},
		token{
			Type:  tokenOpenBracket,
			Value: "(",
		},
		token{
			Type:  tokenDot,
			Value: ".",
		},
		token{
			Type:  tokenSymbol,
			Value: "dum",
		},
		token{
			Type:  tokenDot,
			Value: ".",
		},
		token{
			Type:  tokenSymbol,
			Value: "xyz",
		},
		token{
			Type:  tokenDot,
			Value: ".",
		},
		token{
			Type:  tokenSymbol,
			Value: "sub",
		},
		token{
			Type:  tokenComma,
			Value: ",",
		},
		token{
			Type:  tokenOperator,
			Value: "or",
		},
		token{
			Type:  tokenOpenBracket,
			Value: "(",
		},
		token{
			Type:  tokenString,
			Value: "  any  (( ) and not, textTEXT ",
		},
		token{
			Type:  tokenComma,
			Value: ",",
		},
		token{
			Type:  tokenNumber,
			Value: "2010",
		},
		token{
			Type:  tokenComma,
			Value: ",",
		},
		token{
			Type:  tokenSymbol,
			Value: "xy_z-201",
		},
		token{
			Type:  tokenClosingBracket,
			Value: ")",
		},
		token{
			Type:  tokenClosingBracket,
			Value: ")",
		},
		token{
			Type:  tokenComma,
			Value: ",",
		},
		token{
			Type:  tokenOperator,
			Value: "not",
		},
		token{
			Type:  tokenOpenBracket,
			Value: "(",
		},
		token{
			Type:  tokenString,
			Value: "xxx aa__-c--  zz ",
		},
		token{
			Type:  tokenClosingBracket,
			Value: ")",
		},
	}
	lextester(t, text, tokens)
}

func TestLexBool(t *testing.T) {
	text := " not( len(val1.sub, true),  eq(false, val2.f2) ) "
	tokens := []token{
		token{
			Type:  tokenOperator,
			Value: "not",
		},
		token{
			Type:  tokenOpenBracket,
			Value: "(",
		},
		token{
			Type:  tokenOperator,
			Value: "len",
		},
		token{
			Type:  tokenOpenBracket,
			Value: "(",
		},
		token{
			Type:  tokenSymbol,
			Value: "val1",
		},
		token{
			Type:  tokenDot,
			Value: ".",
		},
		token{
			Type:  tokenSymbol,
			Value: "sub",
		},
		token{
			Type:  tokenComma,
			Value: ",",
		},
		token{
			Type:  tokenTrue,
			Value: "true",
		},
		token{
			Type:  tokenClosingBracket,
			Value: ")",
		},
		token{
			Type:  tokenComma,
			Value: ",",
		},
		token{
			Type:  tokenOperator,
			Value: "eq",
		},
		token{
			Type:  tokenOpenBracket,
			Value: "(",
		},
		token{
			Type:  tokenFalse,
			Value: "false",
		},
		token{
			Type:  tokenComma,
			Value: ",",
		},
		token{
			Type:  tokenSymbol,
			Value: "val2",
		},
		token{
			Type:  tokenDot,
			Value: ".",
		},
		token{
			Type:  tokenSymbol,
			Value: "f2",
		},
		token{
			Type:  tokenClosingBracket,
			Value: ")",
		},
		token{
			Type:  tokenClosingBracket,
			Value: ")",
		},
	}
	lextester(t, text, tokens)
}

func TestTokenConvToStr(t *testing.T) {
	if tokStr := tokenAsStr(tokenSymbol); tokStr != "tokenSymbol" {
		t.Fatalf("Unexpected (%s)", tokStr)
	}
	if tokStr := tokenAsStr(tokenType(100)); tokStr != "Unknown (100)" {
		t.Fatalf("Unexpected (%s)", tokStr)
	}
	if s := fmt.Sprintf("%s", tokenSymbol); s != "tokenSymbol" {
		t.Fatalf("Unexpected (%s)", s)
	}
	if s := fmt.Sprintf("%#v", tokenSymbol); s != "tokenSymbol" {
		t.Fatalf("Unexpected (%s)", s)
	}
}

func TestLexStringNotComplete(t *testing.T) {
	lexer := newTokenizer(NewDefaultOperators())
	_, err := lexer.scan("'string not completed ok")
	if err == nil {
		t.Fatalf("Should not be ok")
	}
}

func TestEscapingSubStr(t *testing.T) {
	bq := `\`
	text := "'any " + bq + "'sub dum" + bq + "'text'"
	tokens := []token{
		token{
			Type:  tokenString,
			Value: "any 'sub dum'text",
		},
	}
	lextester(t, text, tokens)
	text = "'any " + bq + bq + "sub dum" + bq + bq + "text'"
	tokens = []token{
		token{
			Type:  tokenString,
			Value: "any " + bq + "sub dum" + bq + "text",
		},
	}
	lextester(t, text, tokens)
}

func TestExpander(t *testing.T) {
	text := "xyz: dum"
	tokens := []token{
		token{
			Type:  tokenSymbol,
			Value: "xyz",
		},
		token{
			Type:  tokenExpander,
			Value: ":",
		},
		token{
			Type:  tokenSymbol,
			Value: "dum",
		},
	}
	lextester(t, text, tokens)
}
