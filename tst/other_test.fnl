
ns other_test

import ut_fwk

ASSURE = ut_fwk.VERIFY

import common_test_util

testTypeOper = func()
	result = list(
		eq(type('some text'), 'string'),
		eq(type(123), 'int'),
		eq(type(true), 'bool'),
		eq(type(func() 1 end), 'function'),
		eq(type(proc() 1 end), 'function'), # should it be something else..?
		eq(type(chan()), 'channel'),
		eq(type(list(1, 2, 3)), 'list'),
		true
	)
	allRight = call(common_test_util.isAllTrueInList, result)
	call(ASSURE, allRight, plus('Unexpected result = ', str(result)))
end

testStrOper = func()
	import stdstr
	result = list(
		eq(str(123), '123'),
		eq(str(true), 'true'),
		eq(str(not(true)), 'false'),
		eq(str('abc def'), 'abc def'),
		eq(str(list(1, 2, 3)), 'list(1, 2, 3)'),
		call(stdstr.startswith str(func() 1 end), 'func-value'), # maybe should be something else...
		eq(str(chan()), 'chan-value'), # maybe should be something else...
		eq(str( div(float(5), float(10))), '0.5'),
		eq(str(float(50)), '50.0'),
		true
	)
	allRight = call(common_test_util.isAllTrueInList, result)
	call(ASSURE, allRight, plus('Unexpected result = ', str(result)))
end

testTypeConversion = func()
	floatVal1 = div(float(5), float(10))
	floatVal2 = div(float(55), float(10))

	result = list(
		eq(conv(floatVal1, 'float'), floatVal1),
		eq(conv(123, 'float'), float(123)),

		eq(conv(floatVal1, 'string'), '0.5'),
		eq(conv(floatVal2, 'string'), '5.5'),
		eq(conv(float(123), 'int'), 123),
		eq(conv(floatVal1, 'int'), 0),
		eq(conv(floatVal2, 'int'), 5),

		eq(conv('123', 'int'), 123),
		eq(conv(list(1,2,3), 'string'), 'list(1, 2, 3)'),
		eq(conv(123, 'string'), '123'),
		eq(conv(123, 'int'), 123),
		eq(conv(false, 'bool'), false),
		eq(conv('123', 'string'), '123'),
		eq(conv(123, 'int'), 123),
		eq(conv(list(1,2,3), 'list'), list(1, 2, 3)),
		true
	)
	allRight = call(common_test_util.isAllTrueInList, result)
	call(ASSURE, allRight, plus('Unexpected result = ', str(result)))
end

testNameOper = func()
	some_symbol = 123
	result = eq(name(some_symbol), 'some_symbol')
	call(ASSURE, result, plus('Unexpected result = ', str(result)))	
end

testErrorOper = proc()
	result = list(
		eq(try(error('some explanation')), 'RTE:some explanation'),
		eq(try(error()), 'RTE:'),
		eq(try(error('some explanation', '...and more')), 'RTE:some explanation...and more'),
		eq(try(error('some explanation:', 123)), 'RTE:some explanation:123'),
		true
	)
	allRight = call(common_test_util.isAllTrueInList, result)
	call(ASSURE, allRight, plus('Unexpected result = ', str(result)))
end

testSymvalOper = proc()
	some_symbol = 123
	some_other_symbol = 456
	result = list(
		try(symval('nonexisting'), true),
		try(symval('other_test'), true),
		eq(symval('some_symbol'), 123),
		not(eq(symval('some_other_symbol'), 123)),
		eq(symval(name(some_symbol)), 123),
		true
	)
	allRight = call(common_test_util.isAllTrueInList, result)
	call(ASSURE, allRight, plus('Unexpected result = ', str(result)))	
end

testTrylOper = proc()
	result = list(
		eq(tryl(mul(2 3 2)) list(true '' 12))
		eq(tryl(mul(2 3 'X')) list(false 'Invalid type for mul' ''))
	)
	allRight = call(common_test_util.isAllTrueInList, result)
	call(ASSURE, allRight, plus('Unexpected result = ', str(result)))
end

endns
