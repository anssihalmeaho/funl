
ns stdlex_test

import ut_fwk

ASSURE = ut_fwk.VERIFY

import common_test_util
import stdlex

test-tokenize-ok = proc()
	ok err tokens = call(stdlex.tokenize 'symbol 123 #line comment\ntrue'):
	assumed-tokens = list(
		map(
			'line'  1
			'value' 'symbol'
			'pos'   7
			'type'  'tokenSymbol'
		)
		map(
			'line'  1
			'value' '123'
			'pos'   11
			'type'  'tokenNumber'
		)
		map(
			'line'  1
			'value' 'line comment '
			'pos'   26
			'type'  'tokenLineComment'
		)
		map(
			'line'  2
			'value' 'true'
			'pos'   4
			'type'  'tokenTrue'
		)
	)
	and(
		call(ASSURE ok err)
		call(ASSURE eq(tokens assumed-tokens ) sprintf('unexpected tokens: %v' tokens))
	)
end

test-tokenize-nok = proc()
	ok err tokens = call(stdlex.tokenize 'symbol 123 ? #line comment\ntrue'):
	and(
		call(ASSURE not(ok) 'Unexpected success')
	)
end

endns

