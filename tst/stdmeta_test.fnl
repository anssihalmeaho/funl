
ns stdmeta_test

import ut_fwk

ASSURE = ut_fwk.VERIFY

import common_test_util
import stdmeta

test-schema-check-fails = func()
	# define schemas
	subsubschema = map(
		10.5 list(list('type' 'int') list('in' 10 15 25))
	)
	subchema = map(
		'subf1' list(list('required') list('doc' 'this is subfield' '...of sub-map'))
		'subf2' list(list('required') list('type' 'int'))
		'subf3' list(list('map' subsubschema))
	)
	schema = list('map' map(
		'field-1' list(list('required') list('type' 'string'))
		'field-2' list(list('required') list('type' 'list'))
		'field-3' list(list('map' subchema))
	))

	# define data
	subdata = map(
		'subf2' 'not supposed to be string'
		'subf3' map(10.5 10.0 'b' 20)
	)
	data = map(
		'field-1' 'some text'
		'field-2' list(1 2 3)
		'field-3' subdata
	)

	# and validate
	ok msglist = call(stdmeta.validate schema data):

	expected-msg-list = list(
		'required field subf1 not found ( -> field-3)'
		'field subf2 is not required type (got: string, expected: int)( -> field-3)'
		'field 10.5 is not in allowed set (10 not in: list(10, 15, 25))( -> field-3 -> subf3)'
		'field 10.5 is not required type (got: float, expected: int)( -> field-3 -> subf3)'
	)

	import stdfu
	messages-ok = and(
		eq(len(msglist) len(expected-msg-list))
		call(stdfu.foreach msglist func(msg result) and(result in(expected-msg-list msg)) end true)
		call(stdfu.foreach expected-msg-list func(msg result) and(result in(msglist msg)) end true)
	)

	call(ASSURE and(not(ok) messages-ok) plus('Unexpected result = ' str(msglist)))
end

test-schema-ok = func()
	# define schemas
	subsubschema = map(
		10.5 list(list('type' 'int') list('in' 10 15 25))
	)
	subchema = map(
		'subf1' list(list('required') list('doc' 'this is subfield' '...of sub-map'))
		'subf2' list(list('required') list('type' 'int'))
		'subf3' list(list('map' subsubschema))
	)
	schema = list('map' map(
		'field-1' list(list('required') list('type' 'string'))
		'field-2' list(list('required') list('type' 'list'))
		'field-3' list(list('map' subchema))
	))

	# define data
	subdata = map(
		'subf1' 100
		'subf2' 200
		'subf3' map(10.5 10 'b' 20)
	)
	data = map(
		'field-1' 'some text'
		'field-2' list(1 2 3)
		'field-3' subdata
	)

	# and validate
	ok msglist = call(stdmeta.validate schema data):
	call(ASSURE and(ok empty(msglist)) plus('Unexpected result = ' str(msglist)))
end

endns

