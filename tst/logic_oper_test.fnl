
ns logic_oper_test

testAndOperBasics = func()
	result = and(true, true)
	eq(result, true)
end

testOrOperBasics = func()
	result = or(true, true)
	eq(result, true)
end

testAndOperWithIllegalValue = proc()
	expectText = 'error as expected'
	result = try(and(true, 'this is illegal'), expectText)
	eq(result, expectText)
end

testOrOperWithIllegalValue = proc()
	expectText = 'error as expected'
	result = try(or('this is illegal', true), expectText)
	eq(result, expectText)
end

testNotOperWithIllegalValue = proc()
	expectText = 'error as expected'
	result = try(not('this is illegal'), expectText)
	eq(result, expectText)
end

testNotOperBasics = func()
	result = list(
		not(true),
		not(not(true)),
		not(not(not(true)))
	)
	expected = list(false, true, false)

	eq(result, expected)
end

testCombinations = func()
	ttrue = true
	tfalse = false

	result = list(
		and( or(ttrue, tfalse), not(tfalse) ),
		and( not(and(tfalse, tfalse)), or(tfalse, ttrue), and(not(false), true) , false),
		or( and(not(false), ttrue), false, not(true), and(ttrue, ttrue, false)),

		not( and( or(ttrue, tfalse), not(tfalse) ) ),
		not( and( not(and(tfalse, tfalse)), or(tfalse, ttrue), and(not(false), true) , false) ),
		not( or( and(not(false), ttrue), false, not(true), and(ttrue, ttrue, false)) )

		and(true)
		and(false)

		or(true)
		or(false)
	)

	expected = list(true, false, true, false, true, false true false true false)

	eq(result, expected)
end

endns
