
ns stdsort_test

import ut_fwk

ASSURE = ut_fwk.VERIFY

import common_test_util

should-be = func(input expected-output)
	import stdsort

	out = call(stdsort.sort input)
	call(ASSURE eq(out expected-output) sprintf('Result: %v, \nExpected: %v' out expected-output))
end

get-should-be-with-func = func(comparator)
	func(input expected-output)
		import stdsort

		out = call(stdsort.sort input comparator)
		call(ASSURE eq(out expected-output) sprintf('Result: %v, \nExpected: %v' out expected-output))
	end
end

test-list-sortings = func()
	and(list(
		call(should-be list(1) list(1))
		call(should-be list() list())
		call(should-be list(10 5 20) list(5 10 20))
		call(should-be list(1 1) list(1 1))
		call(should-be list(1 1 1) list(1 1 1))
		call(should-be list(10 minus(0 5) 20 0) list(minus(0 5) 0 10 20))
		call(should-be list(minus(0 5) 1 minus(0.0 2.5)) list(minus(0 5) minus(0.0 2.5) 1))
		call(should-be list(10 100 5 50 20) list(5 10 20 50 100))
		call(should-be list(10.5 100 5.5 50.5 20) list(5.5 10.5 20 50.5 100))
	):)
end

test-list-sortings-with-comparison-func = func()
	len-comparator = func(a b)
		lt(len(a) len(b))
	end

	checker = call(get-should-be-with-func len-comparator)

	and(list(
		call(checker list() list())
		call(checker list('a') list('a'))
		call(checker list('ab' 'a') list('a' 'ab'))
		call(checker list('ab' 'abcd' 'a' 'ab' 'abc') list('a' 'ab' 'ab' 'abc' 'abcd'))
	):)
end

endns

