
ns comparison_test

import ut_fwk

testIFcases = func()
	result = list(
		if(true, 'a', 'b'),
		if(not(true), 'a', 'b')
	)
	expected = list(
		'a',
		'b'
	)
	call(ut_fwk.VERIFY, eq(result, expected), plus('Unexpected result = ', str(result)))
end

testEQcases = func()
	l1 = list('aa', 11, 22)

	result = list(
		eq('abc', plus('a', 'bc')),
		eq(100, plus(50,50)),
		eq(true, not(false)),
		eq( list('aa', 11, mul(11, 2)), l1),
		eq(1, '1'),
		eq(1, 1, 2),
		eq(1, 1, 1),
		eq(float(1), float(1), float(1)),
		eq(float(1), float(1), float(2))
	)
	expected = list(
		true,
		true,
		true,
		true,
		false,
		false,
		true,
		true,
		false
	)

	call(ut_fwk.VERIFY, eq(result, expected), plus('Unexpected result = ', str(result)))
end

testGTcases = proc()
	result = list(
		gt(2, 1),
		gt(1, 2),
		gt(1, 1),
		try(gt(1, '1'), true),
		try(gt(1, 1, 2), true),
		try(gt(1), true),

		gt(float(2), float(1)),
		gt(float(2), 1),
		gt(2, float(1)),
		gt(float(1), float(1)),
		gt(float(1), 1),
		gt(1, float(1)),
		gt(float(1), float(2)),
		gt(float(1), 2),
		gt(1, float(2)),

		true
	)
	expected = list(
		true,
		false,
		false,
		true,
		true,
		true,

		true,
		true,
		true,
		false,
		false,
		false,
		false,
		false,
		false,

		true
	)

	call(ut_fwk.VERIFY, eq(result, expected), plus('Unexpected result = ', str(result)))
end

testGEcases = proc()
	result = list(
		ge(2, 1),
		ge(1, 2),
		ge(1, 1),
		try(ge(1, '1'), true),
		try(ge(1, 1, 2), true),
		try(ge(1), true),

		ge(float(2), float(1)),
		ge(float(2), 1),
		ge(2, float(1)),
		ge(float(1), float(1)),
		ge(float(1), 1),
		ge(1, float(1)),
		ge(float(1), float(2)),
		ge(float(1), 2),
		ge(1, float(2)),

		true
	)
	expected = list(
		true,
		false,
		true,
		true,
		true,
		true,

		true,
		true,
		true,
		true,
		true,
		true,
		false,
		false,
		false,

		true
	)

	call(ut_fwk.VERIFY, eq(result, expected), plus('Unexpected result = ', str(result)))
end

testLTcases = proc()
	result = list(
		lt(1, 2),
		lt(2, 1),
		lt(1, 1),
		try(lt(1, '1'), true),
		try(lt(1, 1, 2), true),
		try(lt(1), true),

		lt(float(1), float(2)),
		lt(float(1), 2),
		lt(1, float(2)),
		lt(float(1), float(1)),
		lt(float(1), 1),
		lt(1, float(1)),
		lt(float(2), float(1)),
		lt(float(2), 1),
		lt(2, float(1)),

		true
	)
	expected = list(
		true,
		false,
		false,
		true,
		true,
		true,

		true,
		true,
		true,
		false,
		false,
		false,
		false,
		false,
		false,

		true
	)

	call(ut_fwk.VERIFY, eq(result, expected), plus('Unexpected result = ', str(result)))
end

testLEcases = proc()
	result = list(
		le(1, 2),
		le(2, 1),
		le(1, 1),
		try(le(1, '1'), true),
		try(le(1, 1, 2), true),
		try(le(1), true),

		le(float(1), float(2)),
		le(float(1), 2),
		le(1, float(2)),
		le(float(1), float(1)),
		le(float(1), 1),
		le(1, float(1)),
		le(float(2), float(1)),
		le(float(2), 1),
		le(2, float(1)),

		true
	)
	expected = list(
		true,
		false,
		true,
		true,
		true,
		true,

		true,
		true,
		true,
		true,
		true,
		true,
		false,
		false,
		false,

		true
	)

	call(ut_fwk.VERIFY, eq(result, expected), plus('Unexpected result = ', str(result)))
end

testCondCase = proc()
	yes = true
	result = list(
		try(cond(
			eq(10 11) '1'
			not(true) '2'
			yes       '3'
			'default'
		))
		try(cond(
			eq(10 11) '1'
			true 	  '2'
			yes       '3'
			'default'
		))
		try(cond(
			eq(10 11) '1'
			not(true) '2'
			false     '3'
			'default'
		))
		try(cond(
			eq(10 11) '1'
			'default'
		))
		try(cond(
			eq(10 10) '1'
			'default'
		))
		try(cond(
			eq(10 11)  '1'
			true       '2'
			not(false) '3'
			'default'
		))
		try(cond(
			eq(11 11)  '1'
			true       '2'
			not(false) '3'
			'default'
		))
	)
	
	expected = list(
		'3' '2' 'default' 'default' '1' '2' '1'
	)

	call(ut_fwk.VERIFY, eq(result, expected), plus('Unexpected result = ', str(result)))
end

testCaseCase = func()
	numberAsString = func(d) str(d) end 
	result = list(
		case(
			call(numberAsString, 1),
			'1', 1,
			'2', 2,
			'default'
		),
		case(
			call(numberAsString, 2),
			'1', 1,
			'2', 2,
			'default'
		),
		case(
			call(numberAsString, 3),
			'1', 1,
			'2', 2,
			'default'
		),
		case(
			call(numberAsString, 2),
			'1', 1,
			'2', 2
		)
	)
	expected = list(
		1,
		2,
		'default',
		2
	)

	call(ut_fwk.VERIFY, eq(result, expected), plus('Unexpected result = ', str(result)))
end

testEQMismatchedTypes = proc()
	text = 'got error'
	result = list(
		try(eq(), text)
	)
	expected = list(
		text
	)

	call(ut_fwk.VERIFY, eq(result, expected), plus('Unexpected result = ', str(result)))
end

endns
