
ns list_test

import ut_fwk

ASSURE = ut_fwk.VERIFY

import common_test_util

testListBasics = proc()
	otherlist = list(1, 2, 3)
	testlist = list('abc', 123, true, list(), otherlist)

	result = list(
		eq(head(testlist), 'abc'),
		eq(head(rest(testlist)), 123),
		eq(head(rest(rest(rest(testlist)))), list()),
		eq(head(rest(rest(rest(rest(testlist))))), otherlist),
		eq(head(rest(rest(rest(rest(testlist))))), last(testlist)),
		eq(last(testlist), otherlist),
		eq(len(otherlist), 3),
		eq(len(testlist), 5),
		eq(len(extend(otherlist, testlist)), 8),
		eq(extend(otherlist, testlist, list(10, 20)), list(1, 2, 3, 'abc', 123, true, list(), otherlist, 10, 20)),
		eq(append(otherlist, 'new'), list(1, 2, 3, 'new')),
		eq(append(otherlist, 'new', 'new2'), list(1, 2, 3, 'new', 'new2')),
		eq(add(otherlist, 'new'), list('new', 1, 2, 3)),
		eq(add(otherlist, 'new', 'new2'), list('new', 'new2', 1, 2, 3)),
		eq(empty(list()), true),
		eq(empty(list(list())), false),
		eq(empty(testlist), false),
		eq(rest(otherlist), list(2, 3)),
		eq(rrest(otherlist), list(1, 2)),
		eq(reverse(otherlist), list(3, 2, 1)),
		eq(reverse(list()), list()),
		eq(in(otherlist, 2), true),
		eq(in(otherlist, 4), false),
		eq(ind(testlist, 2), true),
		eq(ind(testlist, 0), 'abc'),
		try(ind(testlist, 100), true),
		eq(find(otherlist, 3), list(2)),
		eq(find(extend(otherlist, list(1, 3, 3)), 3), list(2, 4, 5)),
		eq(find(list(1,2,3), 4), list()),
		eq(slice(testlist, 1, 3), list(123, true, list())),
		eq(slice(testlist, 1, 1), list(123)),
		eq(slice(testlist, 2), list(true, list(), otherlist)),
		eq(slice(otherlist, 100), list()),
		eq(slice(otherlist, 1, 100), list(2, 3))

		eq( extend(list(1 2 3) list(4 5 6) list(7 8 9)) list(1 2 3 4 5 6 7 8 9))
		eq( extend(list() list() list()) list())
		eq( extend(list(1 2)) list(1 2))
		eq( extend(list()) list())
		eq( extend(list(1 2 3) list()) list(1 2 3))

		eq( extend(list() list() list()) list())
		eq( extend(list() list(4 5 6)) list(4 5 6))

		eq( extend(list(1 2) append(list(3) 4) ) list(1 2 3 4))
		eq( extend(list() append(list(3) 4) ) list(3 4))
		eq( extend(list() append(list() 4) ) list(4))
		eq( extend(list(1 2) append(list(3) 4) append(list(5) 6) ) list(1 2 3 4 5 6))
		eq( extend(list(1 2) append(list(3) 4) list(5 6)) list(1 2 3 4 5 6))
	)
	allRight = call(common_test_util.isAllTrueInList, result)
	call(ASSURE, allRight, plus('Unexpected result = ', str(result)))
end

testListTypeErrors = proc()
	notlist = 'this is not list'
	result = list(
		try(head(notlist), true),
		try(last(notlist), true),
		try(extend(notlist, list(1,2)), true),
		try(add(notlist, 1), true),
		try(append(notlist, 1), true),
		try(empty(notlist), true),
		try(rest(notlist), true),
		try(rrest(notlist), true),
		try(reverse(notlist), true),
		true
	)
	allRight = call(common_test_util.isAllTrueInList, result)
	call(ASSURE, allRight, plus('Unexpected result = ', str(result)))
end

endns

