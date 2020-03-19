package funl

import (
	"fmt"
	"log"
	"unicode"
)

const printingOn = false

type lexState int

type tokenType int

const (
	tokenSymbol tokenType = iota
	tokenAs
	tokenStartNS
	tokenEndNS
	tokenNumber
	tokenString
	tokenDot
	tokenOpenBracket
	tokenClosingBracket
	tokenComma
	tokenOperator
	tokenTrue
	tokenFalse
	tokenFuncBegin
	tokenProcBegin
	tokenFuncEnd
	tokenImport
	tokenEqualsSign
	tokenExpander

	receivingWhitespace lexState = iota
	receivingAlphabetsAndNumbers
	receivingCharacters
	receivingCommentLine
	receivingCommentBlock
	unknown
)

func getStateLinking() map[lexState]processor {
	stateLinks := make(map[lexState]processor)
	stateLinks[receivingWhitespace] = &wsState{}
	stateLinks[receivingAlphabetsAndNumbers] = &alnumState{}
	stateLinks[receivingCharacters] = &strState{}
	stateLinks[receivingCommentLine] = &lineCommentState{}
	stateLinks[receivingCommentBlock] = &multilineCommentState{}
	stateLinks[unknown] = &noState{}
	return stateLinks
}

func (tt tokenType) GoString() string {
	return tokenAsStr(tt)
}

func (tt tokenType) String() string {
	return tokenAsStr(tt)
}

type collector interface {
	addChr(rune)
	addChrStr(string)
	getTokenFromBuf(tokenType)
	reportError(error)
}

type processor interface {
	setTokenizer(collector)

	tryProcess(rune) (lexState, bool)

	processWS() lexState
	processAlphanumeric(rune) lexState
	processStringLimiter() lexState
	processSpecCharacter(rune) lexState
	processNewline() lexState
	processLineCommentStart() lexState
	processMultilineStart() lexState
	processMultilineEnd() lexState
	endingOfProcess() error
}

type noState struct {
	tokenizer collector
}

func (s *noState) setTokenizer(tokenizerAPI collector) {
	s.tokenizer = tokenizerAPI
}

func (s *noState) tryProcess(ch rune) (lexState, bool) {
	return unknown, false
}

func (s *noState) endingOfProcess() error {
	return nil
}

func (s *noState) processLineCommentStart() lexState {
	return receivingCommentLine
}

func (s *noState) processNewline() lexState {
	return unknown
}

func (s *noState) processMultilineStart() lexState {
	return receivingCommentBlock
}

func (s *noState) processMultilineEnd() lexState {
	return unknown
}

func (s *noState) processWS() lexState {
	return receivingWhitespace
}

func (s *noState) processAlphanumeric(ch rune) lexState {
	s.tokenizer.addChr(ch)
	return receivingAlphabetsAndNumbers
}

func (s *noState) processStringLimiter() lexState {
	return receivingCharacters
}

func (s *noState) processSpecCharacter(ch rune) lexState {
	s.tokenizer.addChr(ch)
	s.tokenizer.getTokenFromBuf(tokenTypeForCharacter(ch))
	return unknown
}

type wsState struct {
	tokenizer collector
}

func (s *wsState) setTokenizer(tokenizerAPI collector) {
	s.tokenizer = tokenizerAPI
}

func (s *wsState) endingOfProcess() error {
	return nil
}

func (s *wsState) tryProcess(ch rune) (lexState, bool) {
	return unknown, false
}

func (s *wsState) processLineCommentStart() lexState {
	return receivingCommentLine
}

func (s *wsState) processMultilineStart() lexState {
	return receivingCommentBlock
}

func (s *wsState) processMultilineEnd() lexState {
	s.tokenizer.reportError(fmt.Errorf("comment ending not found"))
	return unknown
}

func (s *wsState) processNewline() lexState {
	return unknown
}

func (s *wsState) processWS() lexState {
	return receivingWhitespace
}

func (s *wsState) processAlphanumeric(ch rune) lexState {
	s.tokenizer.addChr(ch)
	return receivingAlphabetsAndNumbers
}

func (s *wsState) processStringLimiter() lexState {
	return receivingCharacters
}

func (s *wsState) processSpecCharacter(ch rune) lexState {
	s.tokenizer.addChr(ch)
	s.tokenizer.getTokenFromBuf(tokenTypeForCharacter(ch))
	return unknown
}

type alnumState struct {
	tokenizer collector
}

func (s *alnumState) setTokenizer(tokenizerAPI collector) {
	s.tokenizer = tokenizerAPI
}

func (s *alnumState) tryProcess(ch rune) (lexState, bool) {
	return unknown, false
}

func (s *alnumState) endingOfProcess() error {
	s.tokenizer.getTokenFromBuf(tokenSymbol)
	return nil
}

func (s *alnumState) processLineCommentStart() lexState {
	s.tokenizer.getTokenFromBuf(tokenSymbol)
	return receivingCommentLine
}

func (s *alnumState) processMultilineStart() lexState {
	s.tokenizer.getTokenFromBuf(tokenSymbol)
	return receivingCommentBlock
}

func (s *alnumState) processMultilineEnd() lexState {
	return unknown
}

func (s *alnumState) processNewline() lexState {
	return unknown
}

func (s *alnumState) processWS() lexState {
	s.tokenizer.getTokenFromBuf(tokenSymbol)
	return receivingWhitespace
}

func (s *alnumState) processAlphanumeric(ch rune) lexState {
	s.tokenizer.addChr(ch)
	return receivingAlphabetsAndNumbers
}

func (s *alnumState) processStringLimiter() lexState {
	s.tokenizer.getTokenFromBuf(tokenSymbol)
	return receivingCharacters
}

func (s *alnumState) processSpecCharacter(ch rune) lexState {
	s.tokenizer.getTokenFromBuf(tokenSymbol)
	s.tokenizer.addChr(ch)
	s.tokenizer.getTokenFromBuf(tokenTypeForCharacter(ch))
	return unknown
}

type strState struct {
	tokenizer  collector
	escapeNext bool
}

func (s *strState) setTokenizer(tokenizerAPI collector) {
	s.tokenizer = tokenizerAPI
}

func (s *strState) tryProcess(ch rune) (lexState, bool) {
	if string(ch) == `\` {
		if !s.escapeNext {
			s.escapeNext = true
			return receivingCharacters, true
		}
	}
	s.clearEscaped()
	s.tokenizer.addChr(ch)
	return receivingCharacters, true
}

func (s *strState) clearEscaped() {
	s.escapeNext = false
}

func (s *strState) endingOfProcess() error {
	s.clearEscaped()
	return fmt.Errorf("String not completed properly")
}

func (s *strState) processLineCommentStart() lexState {
	s.clearEscaped()
	s.tokenizer.addChr('#')
	return receivingCharacters
}

func (s *strState) processMultilineStart() lexState {
	s.clearEscaped()
	s.tokenizer.addChr('*')
	return receivingCharacters
}

func (s *strState) processMultilineEnd() lexState {
	s.clearEscaped()
	s.tokenizer.addChr('/')
	return receivingCharacters
}

func (s *strState) processNewline() lexState {
	s.clearEscaped()
	return unknown
}

func (s *strState) processWS() lexState {
	s.clearEscaped()
	s.tokenizer.addChr(' ')
	return receivingCharacters
}

func (s *strState) processAlphanumeric(ch rune) lexState {
	if s.escapeNext {
		switch ch {
		case 'n':
			s.tokenizer.addChrStr("\n")
			s.clearEscaped()
			return receivingCharacters
		case 'a':
			s.tokenizer.addChrStr("\a")
			s.clearEscaped()
			return receivingCharacters
		case 't':
			s.tokenizer.addChrStr("\t")
			s.clearEscaped()
			return receivingCharacters
		case 'f':
			s.tokenizer.addChrStr("\f")
			s.clearEscaped()
			return receivingCharacters
		case 'r':
			s.tokenizer.addChrStr("\r")
			s.clearEscaped()
			return receivingCharacters
		case 'v':
			s.tokenizer.addChrStr("\v")
			s.clearEscaped()
			return receivingCharacters
		}
	}

	s.clearEscaped()
	s.tokenizer.addChr(ch)
	return receivingCharacters
}

func (s *strState) processStringLimiter() lexState {
	if s.escapeNext {
		s.clearEscaped()
		s.tokenizer.addChrStr("'")
		return receivingCharacters
	}
	s.clearEscaped()
	s.tokenizer.getTokenFromBuf(tokenString)
	return unknown
}

func (s *strState) processSpecCharacter(ch rune) lexState {
	s.clearEscaped()
	s.tokenizer.addChr(ch)
	return receivingCharacters
}

type lineCommentState struct {
	tokenizer collector
}

func (s *lineCommentState) setTokenizer(tokenizerAPI collector) {
	s.tokenizer = tokenizerAPI
}

func (s *lineCommentState) tryProcess(ch rune) (lexState, bool) {
	return receivingCommentLine, true
}

func (s *lineCommentState) endingOfProcess() error {
	return nil
}

func (s *lineCommentState) processLineCommentStart() lexState {
	return receivingCommentLine
}

func (s *lineCommentState) processMultilineStart() lexState {
	s.tokenizer.addChr('*')
	return receivingCommentLine
}

func (s *lineCommentState) processMultilineEnd() lexState {
	s.tokenizer.addChr('/')
	return receivingCommentLine
}

func (s *lineCommentState) processNewline() lexState {
	return unknown
}

func (s *lineCommentState) processWS() lexState {
	return receivingCommentLine
}

func (s *lineCommentState) processAlphanumeric(ch rune) lexState {
	return receivingCommentLine
}

func (s *lineCommentState) processStringLimiter() lexState {
	return receivingCommentLine
}

func (s *lineCommentState) processSpecCharacter(ch rune) lexState {
	return receivingCommentLine
}

type multilineCommentState struct {
	tokenizer collector
}

func (s *multilineCommentState) setTokenizer(tokenizerAPI collector) {
	s.tokenizer = tokenizerAPI
}

func (s *multilineCommentState) tryProcess(ch rune) (lexState, bool) {
	return receivingCommentBlock, true
}

func (s *multilineCommentState) endingOfProcess() error {
	return fmt.Errorf("Comment not completed properly")
}

func (s *multilineCommentState) processLineCommentStart() lexState {
	return receivingCommentBlock
}

func (s *multilineCommentState) processMultilineStart() lexState {
	return receivingCommentBlock
}

func (s *multilineCommentState) processMultilineEnd() lexState {
	return unknown
}

func (s *multilineCommentState) processNewline() lexState {
	return receivingCommentBlock
}

func (s *multilineCommentState) processWS() lexState {
	return receivingCommentBlock
}

func (s *multilineCommentState) processAlphanumeric(ch rune) lexState {
	return receivingCommentBlock
}

func (s *multilineCommentState) processStringLimiter() lexState {
	return receivingCommentBlock
}

func (s *multilineCommentState) processSpecCharacter(ch rune) lexState {
	return receivingCommentBlock
}

func (s *tokenizer) scan(sourceText string) ([]token, error) {
	var currentState = s.stateLinks[unknown]
	var newState lexState
	var prevC rune
	var makeSureItsStarNext bool

	s.lineno = 1
	s.pos = 0
	for _, c := range sourceText {
		s.pos++

		if makeSureItsStarNext && c != '*' {
			return []token{}, fmt.Errorf("Line:%d Pos:%d: Illegal comment start (%s)", s.lineno, s.pos, string(c))
		}
		makeSureItsStarNext = false

		if c == '/' {
			if prevC != 0 && prevC == '*' {
				newState = currentState.processMultilineEnd()
				goto nextPlease
			}
		} else if c == '*' {
			if prevC != 0 && prevC == '/' {
				newState = currentState.processMultilineStart()
				goto nextPlease
			}
		}

		if isAlNum(c) {
			newState = currentState.processAlphanumeric(c)
		} else if isNewline(c) {
			newState = currentState.processNewline()
			s.lineno++
			s.pos = 0
		} else if isLineCommentStart(c) {
			newState = currentState.processLineCommentStart()
		} else if isWhitespace(c) {
			newState = currentState.processWS()
		} else if isStrLimiter(c) {
			newState = currentState.processStringLimiter()
		} else if isDot(c) {
			newState = currentState.processSpecCharacter(c)
		} else if isOpeningBracket(c) {
			newState = currentState.processSpecCharacter(c)
		} else if isClosingBracket(c) {
			newState = currentState.processSpecCharacter(c)
		} else if isComma(c) {
			newState = currentState.processSpecCharacter(c)
		} else if isEqualSign(c) {
			newState = currentState.processSpecCharacter(c)
		} else if isExpander(c) {
			newState = currentState.processSpecCharacter(c)
		} else {
			if nState, ok := currentState.tryProcess(c); ok {
				newState = nState
			} else if c == '/' {
				prevC = c
				makeSureItsStarNext = true
				continue
			} else {
				if printingOn {
					log.Printf("unexpected: %s", string(c))
				}
				return []token{}, fmt.Errorf("Line:%d Pos:%d: Illegal character (%s)", s.lineno, s.pos, string(c))
			}
		}
		if s.err != nil {
			return []token{}, fmt.Errorf("Line:%d Pos:%d: %v", s.lineno, s.pos, s.err)
		}
	nextPlease:
		currentState = s.stateLinks[newState]
		prevC = c
	}
	err := currentState.endingOfProcess()
	return s.tokens, err
}

type token struct {
	Type   tokenType
	Value  string
	Lineno int
	Pos    int
}

func isAlNum(ch rune) bool {
	if ch == '-' || ch == '_' {
		return true
	}
	if unicode.IsDigit(ch) {
		return true
	}
	if 'a' <= ch && ch <= 'z' {
		return true
	}
	if 'A' <= ch && ch <= 'Z' {
		return true
	}
	return false
}

func isNewline(ch rune) bool {
	if ch == '\n' {
		return true
	}
	return false
}

func isLineCommentStart(ch rune) bool {
	if ch == '#' {
		return true
	}
	return false
}

func isWhitespace(ch rune) bool {
	return unicode.IsSpace(ch)
}

func isStrLimiter(ch rune) bool {
	if string(ch) == "'" {
		return true
	}
	return false
}

func isEqualSign(ch rune) bool {
	return ch == '='
}

func isDot(ch rune) bool {
	return ch == '.'
}

func isOpeningBracket(ch rune) bool {
	return ch == '('
}

func isClosingBracket(ch rune) bool {
	return ch == ')'
}

func isExpander(ch rune) bool {
	return ch == ':'
}

func isComma(ch rune) bool {
	return ch == ','
}

type lexerAPI interface {
	scan(string) ([]token, error)
}

type tokenizer struct {
	buffer     []byte
	tokens     []token
	stateLinks map[lexState]processor
	operators  Operators
	lineno     int
	pos        int
	err        error
}

func (s *tokenizer) reportError(err error) {
	s.err = err
}

func (s *tokenizer) addChrStr(ch string) {
	s.buffer = append([]byte(s.buffer), ch...)
}

func (s *tokenizer) addChr(ch rune) {
	s.buffer = append([]byte(s.buffer), string(ch)...)
}

func (s *tokenizer) getTokenFromBuf(tokenKind tokenType) {
	getOperatorToken := func() tokenType {
		if s.operators.isOperator(string(s.buffer)) {
			return tokenOperator
		}
		return tokenSymbol
	}

	getTokenForNumeral := func() tokenType {
		for _, r := range string(s.buffer) {
			if !unicode.IsDigit(r) {
				return tokenSymbol
			}
		}
		return tokenNumber
	}

	if tokenKind == tokenSymbol {
		tokenKind = getOperatorToken()
	}
	if tokenKind == tokenSymbol {
		tokenKind = getTokenForNumeral()
	}
	if tokenKind == tokenSymbol {
		switch string(s.buffer) {
		case "ns":
			tokenKind = tokenStartNS
		case "endns":
			tokenKind = tokenEndNS
		case "true":
			tokenKind = tokenTrue
		case "false":
			tokenKind = tokenFalse
		case "func":
			tokenKind = tokenFuncBegin
		case "proc":
			tokenKind = tokenProcBegin
		case "end":
			tokenKind = tokenFuncEnd
		case "as":
			tokenKind = tokenAs
		case "import":
			tokenKind = tokenImport
		}
	}
	if printingOn {
		log.Printf("token: <%s> (%d)", s.buffer, tokenKind)
	}
	token := token{
		Type:   tokenKind,
		Value:  string(s.buffer),
		Lineno: s.lineno,
		Pos:    s.pos,
	}
	s.tokens = append(s.tokens, token)
	s.buffer = []byte{}
}

func newTokenizer(operators Operators) *tokenizer {
	newtokenizer := tokenizer{stateLinks: getStateLinking(), operators: operators}
	for _, v := range newtokenizer.stateLinks {
		v.setTokenizer(&newtokenizer)
	}
	return &newtokenizer
}

func tokenAsStr(tt tokenType) string {
	str, ok := map[tokenType]string{
		tokenEqualsSign:     "tokenEqualsSign",
		tokenAs:             "tokenAs",
		tokenFuncBegin:      "tokenFuncBegin",
		tokenProcBegin:      "tokenProcBegin",
		tokenFuncEnd:        "tokenFuncEnd",
		tokenImport:         "tokenImport",
		tokenSymbol:         "tokenSymbol",
		tokenStartNS:        "tokenStartNS",
		tokenEndNS:          "tokenEndNS",
		tokenNumber:         "tokenNumber",
		tokenString:         "tokenString",
		tokenDot:            "tokenDot",
		tokenOpenBracket:    "tokenOpenBracket",
		tokenClosingBracket: "tokenClosingBracket",
		tokenComma:          "tokenComma",
		tokenOperator:       "tokenOperator",
		tokenTrue:           "tokenTrue",
		tokenFalse:          "tokenFalse",
		tokenExpander:       "tokenExpander",
	}[tt]
	if !ok {
		return fmt.Sprintf("Unknown (%d)", int(tt))
	}
	return str
}

func tokenTypeForCharacter(ch rune) tokenType {
	if isComma(ch) {
		return tokenComma
	} else if isDot(ch) {
		return tokenDot
	} else if isOpeningBracket(ch) {
		return tokenOpenBracket
	} else if isEqualSign(ch) {
		return tokenEqualsSign
	} else if isExpander(ch) {
		return tokenExpander
	}
	return tokenClosingBracket
}
