
ns string_test

import ut_fwk

testSprintf = func()
	result = list(
		sprintf('%s:%d:%s' 'blaa' 100 plus('aa' 'bb'))
		sprintf('%.2f:%v:%x' 0.25 false 'abc')
		sprintf('')
	)
	expected = list(
		'blaa:100:aabb'
		'0.25:false:616263'
		''
	)
	call(ut_fwk.VERIFY eq(result expected) plus('Unexpected result = ' str(result)))
end

testStringCatenation = func()
	str1 = 'aaa'
	str2 = 'bbb'
	getstr = func() 'ddd' end

	result = plus(str1, str2, 'ccc', call(getstr))
	expected = 'aaabbbcccddd'
	cond1 = call(ut_fwk.VERIFY, eq(result, expected), plus('Unexpected result = ', str(result)))
	cond2 = and( eq(len(str1), 3), eq(len(result), 12) )
	and(cond1, cond2)
end

testStringLen = func()
	str1 = ''
	str2 = 'abc'
	result = list(
		len(str1),
		len(str2)
	)
	expected = list(
		0,
		3
	)
	call(ut_fwk.VERIFY, eq(result, expected), plus('Unexpected result = ', str(result)))
end

testSubstringQuery = func()
	mainstr = 'abc def ghi'
	substr1 = ' def'
	substr2 = 'chi'
	substr3 = 'ghi'
	substr4 = 'abc'

	result = list(
		in(mainstr, ' def'),
		in(mainstr, substr2),
		in(mainstr, 'ghi'),
		in(mainstr, substr4)
	)
	expected = list(
		true,
		false,
		true,
		true
	)
	call(ut_fwk.VERIFY, eq(result, expected), plus('Unexpected result = ', str(result)))
end

testConvToStr = func()
	x = 100
	result = list(
		str(10),
		str(true),
		str(x),
		str('text'),
		str(list(1, 2, 3))
	)
	expected = list(
		'10',
		'true',
		'100',
		'text',
		'list(1, 2, 3)'
	)
	call(ut_fwk.VERIFY, eq(result, expected), plus('Unexpected result = ', str(result)))
end

testStartsWith = func()
	import stdstr
	startswith = stdstr.startswith
	
	s1 = 'begins blaa'
	s2 = 'begins'
	s3 = 'not this'
	result = list(
		call(startswith, 'abc def ghi', 'a'),
		call(startswith, 'abc def ghi', 'abc d'),
		call(startswith, 'abc def ghi', 'b'),
		call(startswith, 'abc def ghi', 'bc'),
		call(startswith, 'abc def ghi', 'ghi'),
		call(startswith, 'abc def ghi', 'abc def ghi'),
		call(startswith, s1, s2),
		call(startswith, s1, s3),
		call(startswith, 'zzz', s2)
	)
	expected = list(
		true,
		true,
		false,
		false,
		false,
		true,
		true,
		false,
		false
	)
	call(ut_fwk.VERIFY, eq(result, expected), plus('Unexpected result = ', str(result)))
end

testEndsWith = func()
	import stdstr
	endswith = stdstr.endswith

	s1 = 'begins blaa'
	s2 = 'blaa'
	s3 = 'not this'
	result = list(
		call(endswith, 'abc def ghi', 'a'),
		call(endswith, 'abc def ghi', 'abc d'),
		call(endswith, 'abc def ghi', 'i'),
		call(endswith, 'abc def ghi', 'hi'),
		call(endswith, 'abc def ghi', 'f ghi'),
		call(endswith, 'abc def ghi', 'abc def ghi'),
		call(endswith, s1, s2),
		call(endswith, s1, s3),
		call(endswith, 'zzz', s2)
	)
	expected = list(
		false,
		false,
		true,
		true,
		true,
		true,
		true,
		false,
		false
	)
	call(ut_fwk.VERIFY, eq(result, expected), plus('Unexpected result = ', str(result)))
end

testSplitString = func()
	result = list(
		split('abcdabcdabcd', 'bc'),
		split('abcdabcdabcd', 'b'),
		split('abcd  abcd    abcd'),
		split('abcd 1234 abcd', ' '),

		split('abcdefgh', '')
	)

	expected = list(
		list('a', 'da', 'da', 'd'),
		list('a', 'cda', 'cda', 'cd'),
		list('abcd', 'abcd', 'abcd'),
		list('abcd', '1234', 'abcd'),

		list('a', 'b', 'c', 'd', 'e', 'f', 'g', 'h')
	)
	call(ut_fwk.VERIFY, eq(result, expected), plus('Unexpected result = ', str(result)))
end

testSliceOfString = func()
	result = list(
		slice('some nonsense sentence', 5),
		slice('some nonsense sentence', 5, 12)
	)
	
	expected = list(
		'nonsense sentence',
		'nonsense'
	)
	call(ut_fwk.VERIFY, eq(result, expected), plus('Unexpected result = ', str(result)))
end

testFindSubstring = func()
	result = list(
		find('there is substring here', 'substring'),
		find('there is substring here and here substring again', 'substring'),
		find('there is substring here', 'xxx')
	)
	expected = list(
		list(9),
		list(9, 33),
		list()
	)
	call(ut_fwk.VERIFY, eq(result, expected), plus('Unexpected result = ', str(result)))
end

testIndForString = proc()
	result = list(
		ind('0123456789', 6),
		ind('0123456789', 4),
		ind('0123456789', 0),
		ind('0123456789', 9),
		try(ind('0123456789', 10), 'index over')
	)
	expected = list(
		'6',
		'4',
		'0',
		'9',
		'index over'
	)
	call(ut_fwk.VERIFY, eq(result, expected), plus('Unexpected result = ', str(result)))
end

testLenOfString = func()
	result = list(
		len('123456789'),
		len('')
	)
	expected = list(
		9,
		0
	)
	call(ut_fwk.VERIFY, eq(result, expected), plus('Unexpected result = ', str(result)))
end

testStringEscapeCaharacters = func()
	result = list(
		split('123\'456\'789', '\''),
		split('123\'\'456\'789', '\''),
		split('123\\456\\789', '\\'),
		split('123\\\\456\\\\789', '\\\\'),
		split('123\\\\456\\\\789', '\\'),
		split('123\\\'\\789', '\''),
		split('123\'\\\'789', '\'')
	)
	expected = list(
		list('123', '456', '789'),
		list('123', '', '456', '789'),
		list('123', '456', '789'),
		list('123', '456', '789'),
		list('123', '', '456', '', '789'),
		list('123\\', '\\789'),
		list('123', '\\', '789')
	)
	call(ut_fwk.VERIFY, eq(result, expected), plus('Unexpected result = ', str(result)))
end

testStrForTypes = func()
	import stdstr

	result = list(
		str('123456789'),
		str(123),
		str(0.5),
		str(list(1,2,3)),
		str(map(1, 10, 2, 20)),
		str(true),
		str(not(true)),
		call(stdstr.startswith str(func() 1 end) 'func-value:'),
		call(stdstr.startswith str(proc() 1 end) 'func-value:'),
		str(chan())
		# missing: opaque type
		# missing: external proc
	)
	expected = list(
		'123456789',
		'123',
		'0.5',
		'list(1, 2, 3)',
		'map(1 : 10, 2 : 20)',
		'true',
		'false',
		true
		true
		'chan-value'
	)
	call(ut_fwk.VERIFY, eq(result, expected), plus('Unexpected result = ', str(result)))
end

endns

