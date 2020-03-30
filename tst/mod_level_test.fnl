
ns mod_level_test

import ut_fwk

ASSURE = ut_fwk.VERIFY

import common_test_util

someTestInt = 1
someTestStr = 'this is test'
someTestFunc = func() 123 end
someTestChan = chan()
someTestFloat = div(float(5), float(10))

testModLevelSymbol = func()
	result = list(
		someTestInt,
		type(someTestStr),
		call(someTestFunc),
		type(someTestChan),
		type(someTestFloat),
		true
	)
	expected = list(
		1,
		'string',
		123,
		'channel',
		'float',
		true
	)
	call(ASSURE, eq(result, expected), plus('Unexpected result = ', str(result)))
end

endns
