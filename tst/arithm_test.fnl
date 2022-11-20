
ns arithm_test

testPlusWithSingleArgument = func()
	and(
		eq(plus(10) 10)
		eq(plus(10.25) 10.25)
		eq(plus('ok') 'ok')
	)
end

testMulWithSingleArgument = func()
	and(
		eq(mul(10) 10)
		eq(mul(0.5) 0.5)
	)
end

testSumOfIntegers = func()
	result = plus(1 ,2, 3)
	eq(result, 6)
end

testDivisionByZeroForFloat = proc()
	text = 'assumed'
	result = try(div(10, float(0)), text)
	result2 = try(div(10, 0.0), text)
	and(eq(result, text), eq(result2, text))
end

testSumOfFloats = func()
	fval1 = float(10)
	fval2 = div(float(3), float(10))
	fval3 = div(float(2), float(10))

	expected = div(float(105), float(10)) # 10.5

	result = plus(fval1, fval2, fval3, 0.0)
	eq(result, expected)
end

testDiffOfFloats = func()
	x = div(float(5), float(10))
	y = div(float(2), float(10))
	result = minus(x, y)
	eq(result, div(float(3), float(10)))
end

testDiffOfIntegers = func()
	result = minus(6, 2)
	eq(result, 4)
end

testFloatAndIntDivisions = func()
	result = list(
		div(10, 2),
		type(div(10, 2)),
		div(float(10), 2),
		type(div(float(10), 2)),
		div(10, float(2)),
		type(div(10, float(2))),
		div(float(10), float(2)),
		type(div(float(10), float(2)))
	)
	expected = list(
		5,
		'int',
		float(5),
		'float',
		float(5),
		'float',
		float(5),
		'float'
	)
	eq(result, expected)
end

testMultiplyOfFloats = func()
	x = div(float(5), float(100)) # 0.05
	y = float(10)
	expected = div(float(5), float(10)) # 0.5
	result = mul(x, y)
	and(eq(result, expected), eq(type(result), 'float'))
end

testMultiplyOfFloats2 = func()
	x = 0.05
	y = 10.0
	expected = 0.5
	result = mul(x, y)
	and(eq(result, expected), eq(type(result), 'float'))
end

testMultiplyOfIntsAndFloats = func()
	x = div(float(5), float(100)) # 0.05
	y = 10
	expected = float(2)
	result = mul(x, y, 4)
	and(eq(result, expected), eq(type(result), 'float'))
end

testMultiplyOfIntegers = func()
	result = mul(6, 2, 10)
	eq(result, 120)
end

testDivisionOfIntegers = func()
	result = list(
		div(20, 3),
		div(20, 10),
		div(2, 6),
		mod(20, 3),
		mod(20, 10),
		mod(2, 6)
	)
	expected = list(6, 2, 0, 2, 0, 2)
	eq(result, expected)
end

testCombinations = func()
	result = list(
		plus( mul(mul(minus(9,1), 2), div(4, 2)), 8),
		plus( minus(1, 1), minus(1, 5), div(4, 1))
	)
	expected = list(40, 0)
	eq(result, expected)
end

testDivisionByZero = proc()
	text = 'assumed'
	result = try(div(10, 0), text)
	eq(result, text)
end

testInvalidTypes = proc()
	text = 'assumed'
	result = list(
		try(div('abc', 1), text),
		try(mod(true, 1), text),
		try(plus(10, list()), text),
		try(minus('abc', 0), text),
		try(mul('abc', 1), text)
	)

	eq(result, list(text,text,text,text,text))
end

endns
