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
			t.Logf("Tokens: %#v", tokens)
			t.Fatalf("Unexpected type (%d): %v, %v (%s)(%s)", i, tokens[i].Type, token.Type, tokens[i].Value, token.Value)
		}
		if tokens[i].Value != token.Value {
			t.Logf("Tokens: %#v", tokens)
			t.Fatalf("Unexpected value (%d): %s, %s", i, tokens[i].Value, token.Value)
		}
	}
}

func TestJustSymbols(t *testing.T) {
	text := "2010"
	tokens := []token{
		{
			Type:  tokenNumber,
			Value: "2010",
		},
	}
	lextester(t, text, tokens)
	text = "dum_name"
	tokens = []token{
		{
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
		{
			Type:  tokenSymbol,
			Value: "not",
		},
		{
			Type:  tokenOpenBracket,
			Value: "(",
		},
		{
			Type:  tokenSymbol,
			Value: "dum",
		},
		{
			Type:  tokenDot,
			Value: ".",
		},
		{
			Type:  tokenSymbol,
			Value: "sub",
		},
		{
			Type:  tokenComma,
			Value: ",",
		},
		{
			Type:  tokenString,
			Value: "aabbcc text",
		},
		{
			Type:  tokenSymbol,
			Value: "any_more",
		},
		{
			Type:  tokenString,
			Value: "what text",
		},
		{
			Type:  tokenClosingBracket,
			Value: ")",
		},
	}
	lextester(t, text, tokens)
}

func TestLexNumbers(t *testing.T) {
	text := "not(2007, 'dum, 2010 text', value1, value2)"
	tokens := []token{
		{
			Type:  tokenSymbol,
			Value: "not",
		},
		{
			Type:  tokenOpenBracket,
			Value: "(",
		},
		{
			Type:  tokenNumber,
			Value: "2007",
		},
		{
			Type:  tokenComma,
			Value: ",",
		},
		{
			Type:  tokenString,
			Value: "dum, 2010 text",
		},
		{
			Type:  tokenComma,
			Value: ",",
		},
		{
			Type:  tokenSymbol,
			Value: "value1",
		},
		{
			Type:  tokenComma,
			Value: ",",
		},
		{
			Type:  tokenSymbol,
			Value: "value2",
		},
		{
			Type:  tokenClosingBracket,
			Value: ")",
		},
	}
	lextester(t, text, tokens)
}

func TestVariousTokens(t *testing.T) {
	text := "  gt  ( or(.dum.xyz.sub, or( '  any  (( ) and not, textTEXT ' , 2010, xy_z-201)) , not('xxx aa__-c--  zz ')  "
	tokens := []token{
		{
			Type:  tokenSymbol,
			Value: "gt",
		},
		{
			Type:  tokenOpenBracket,
			Value: "(",
		},
		{
			Type:  tokenSymbol,
			Value: "or",
		},
		{
			Type:  tokenOpenBracket,
			Value: "(",
		},
		{
			Type:  tokenDot,
			Value: ".",
		},
		{
			Type:  tokenSymbol,
			Value: "dum",
		},
		{
			Type:  tokenDot,
			Value: ".",
		},
		{
			Type:  tokenSymbol,
			Value: "xyz",
		},
		{
			Type:  tokenDot,
			Value: ".",
		},
		{
			Type:  tokenSymbol,
			Value: "sub",
		},
		{
			Type:  tokenComma,
			Value: ",",
		},
		{
			Type:  tokenSymbol,
			Value: "or",
		},
		{
			Type:  tokenOpenBracket,
			Value: "(",
		},
		{
			Type:  tokenString,
			Value: "  any  (( ) and not, textTEXT ",
		},
		{
			Type:  tokenComma,
			Value: ",",
		},
		{
			Type:  tokenNumber,
			Value: "2010",
		},
		{
			Type:  tokenComma,
			Value: ",",
		},
		{
			Type:  tokenSymbol,
			Value: "xy_z-201",
		},
		{
			Type:  tokenClosingBracket,
			Value: ")",
		},
		{
			Type:  tokenClosingBracket,
			Value: ")",
		},
		{
			Type:  tokenComma,
			Value: ",",
		},
		{
			Type:  tokenSymbol,
			Value: "not",
		},
		{
			Type:  tokenOpenBracket,
			Value: "(",
		},
		{
			Type:  tokenString,
			Value: "xxx aa__-c--  zz ",
		},
		{
			Type:  tokenClosingBracket,
			Value: ")",
		},
	}
	lextester(t, text, tokens)
}

func TestLexBool(t *testing.T) {
	text := " not( len(val1.sub, true),  eq(false, val2.f2) ) "
	tokens := []token{
		{
			Type:  tokenSymbol,
			Value: "not",
		},
		{
			Type:  tokenOpenBracket,
			Value: "(",
		},
		{
			Type:  tokenSymbol,
			Value: "len",
		},
		{
			Type:  tokenOpenBracket,
			Value: "(",
		},
		{
			Type:  tokenSymbol,
			Value: "val1",
		},
		{
			Type:  tokenDot,
			Value: ".",
		},
		{
			Type:  tokenSymbol,
			Value: "sub",
		},
		{
			Type:  tokenComma,
			Value: ",",
		},
		{
			Type:  tokenTrue,
			Value: "true",
		},
		{
			Type:  tokenClosingBracket,
			Value: ")",
		},
		{
			Type:  tokenComma,
			Value: ",",
		},
		{
			Type:  tokenSymbol,
			Value: "eq",
		},
		{
			Type:  tokenOpenBracket,
			Value: "(",
		},
		{
			Type:  tokenFalse,
			Value: "false",
		},
		{
			Type:  tokenComma,
			Value: ",",
		},
		{
			Type:  tokenSymbol,
			Value: "val2",
		},
		{
			Type:  tokenDot,
			Value: ".",
		},
		{
			Type:  tokenSymbol,
			Value: "f2",
		},
		{
			Type:  tokenClosingBracket,
			Value: ")",
		},
		{
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
		{
			Type:  tokenString,
			Value: "any 'sub dum'text",
		},
	}
	lextester(t, text, tokens)
	text = "'any " + bq + bq + "sub dum" + bq + bq + "text'"
	tokens = []token{
		{
			Type:  tokenString,
			Value: "any " + bq + "sub dum" + bq + "text",
		},
	}
	lextester(t, text, tokens)
}

func TestExpander(t *testing.T) {
	text := "xyz: dum"
	tokens := []token{
		{
			Type:  tokenSymbol,
			Value: "xyz",
		},
		{
			Type:  tokenExpander,
			Value: ":",
		},
		{
			Type:  tokenSymbol,
			Value: "dum",
		},
	}
	lextester(t, text, tokens)
}

func TestLineComments(t *testing.T) {
	text := "abc efg # this is line comment\n hmm"
	tokens := []token{
		{
			Type:  tokenSymbol,
			Value: "abc",
		},
		{
			Type:  tokenSymbol,
			Value: "efg",
		},
		{
			Type:  tokenLineComment,
			Value: " this is line comment ",
		},
		{
			Type:  tokenSymbol,
			Value: "hmm",
		},
	}
	lextester(t, text, tokens)

	text = "abc efg # this is line comment"
	tokens = []token{
		{
			Type:  tokenSymbol,
			Value: "abc",
		},
		{
			Type:  tokenSymbol,
			Value: "efg",
		},
		{
			Type:  tokenLineComment,
			Value: " this is line comment",
		},
	}
	lextester(t, text, tokens)

	text = "# this is line comment"
	tokens = []token{
		{
			Type:  tokenLineComment,
			Value: " this is line comment",
		},
	}
	lextester(t, text, tokens)
}

func TestMultiLineComments(t *testing.T) {
	text := "abc efg /* multi line\ncomment is \nhere */ hmm"

	lexer := newTokenizer(NewDefaultOperators())
	tokens, err := lexer.scan(text)
	if err != nil {
		t.Fatalf("error : %v", err)
	}
	t.Logf("Tokens: %#v", tokens)
}
