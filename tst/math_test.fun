
ns math_test

import ut_fwk

testAbsCases = func()
	import stdmath

	result = list(
		call(stdmath.abs float(2))
		call(stdmath.abs float(minus(0 2)))
	)
	expected = list(
		2.0
		2.0
	)
	call(ut_fwk.VERIFY eq(result expected) plus('Unexpected result = ' str(result)))
end

testPowLogCases = func()
	import stdmath

	result = list(
		call(stdmath.exp 0.0)
		call(stdmath.exp2 2.0)
		call(stdmath.expm1 0.0)
		call(stdmath.log 1.0)
		call(stdmath.log10 100.0)
		call(stdmath.log1p 0.0)
		call(stdmath.log2 256.0)
		call(stdmath.logb 8.0)
		call(stdmath.pow 2.0 3.0)
		call(stdmath.pow 1.5 3.0)
		call(stdmath.pow 10.0 3.0)
		call(stdmath.sqrt 16.0)
		call(stdmath.cbrt 27.0)
	)
	expected = list(
		1.0
		4.0
		0.0
		0.0
		2.0
		0.0
		8.0
		3.0
		8.0
		3.375
		1000.0
		4.0
		3.0
	)
	call(ut_fwk.VERIFY eq(result expected) plus('Unexpected result = ' str(result)))
end

testLimitCases = func()
	import stdmath

	result = list(
		call(stdmath.ceil 1.49)
		call(stdmath.ceil 0.0)
		call(stdmath.is-nan call(stdmath.ceil call(stdmath.acos 2.0)))

		call(stdmath.floor 1.51)

		call(stdmath.trunc 10.5)

		call(stdmath.frexp 10.0)

		call(stdmath.ldexp 0.625 4)
		call(stdmath.ldexp call(stdmath.frexp 10.0):)
		
		call(stdmath.modf 10.75)

		call(stdmath.remainder 7.5 2.0)
	)
	expected = list(
		2.0
		0.0
		true
		
		1.0

		10.0
		
		list(0.625 4)
		
		10.0
		10.0
		
		list(10.0 0.75)
		
		1.5
	)
	call(ut_fwk.VERIFY eq(result expected) plus('Unexpected result = ' str(result)))
end

testTrigoCases = func()
	import stdmath

	result = list(
		call(stdmath.acos 1.0)
		call(stdmath.is-nan call(stdmath.acos 2.0))

		call(stdmath.acosh 1.0)
		call(stdmath.is-nan call(stdmath.acosh 0.5))

		call(stdmath.asin 0.0)
		call(stdmath.is-nan call(stdmath.asin 2.0))

		call(stdmath.asin 0.0)

		call(stdmath.atan 0.0)

		call(stdmath.atanh 0.0)
		call(stdmath.is-inf call(stdmath.atanh 1.0) '+')

		call(stdmath.cos stdmath.pi)
		sprintf('%.2f' call(stdmath.cos div(stdmath.pi 2)))

		call(stdmath.cosh 0.0)

		call(stdmath.sin 0.0)
		sprintf('%.2f' call(stdmath.sin stdmath.pi))

		call(stdmath.sinh 0.0)

		call(stdmath.tan 0.0)

		call(stdmath.tanh 0.0)
	)
	expected = list(
		0.0
		true
		
		0.0
		true
		
		0.0
		true

		0.0

		0.0

		0.0
		true
		
		float(minus(0 1))
		'0.00'
		
		1.0
		
		0.0
		'0.00'
		
		0.0
		
		0.0
		
		0.0
	)
	call(ut_fwk.VERIFY eq(result expected) plus('Unexpected result = ' str(result)))
end

endns

