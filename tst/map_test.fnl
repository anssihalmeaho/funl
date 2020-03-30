
ns map_test

import ut_fwk

ASSURE = ut_fwk.VERIFY

import common_test_util

testDeletingALotFromMap = proc()
	import stdfu
	
	limit = 2000
	nums = call(stdfu.generate, 0, minus(limit, 1), func(i) i end)
	strs = call(stdfu.apply, nums, func(x) str(x) end)	
	big-m1 = call(stdfu.zip, nums, strs)
	#_ = print('big-m1 : ', big-m1)

	dkeys = call(stdfu.generate, 0, conv(div(limit, 1.002), 'int'), func(i) i end)
	#_ = print('dkeys : ', last(dkeys))

	remove-keys = func(l, prevm)
		while( not(empty(l)),
			rest(l),
			del(prevm, head(l)),
			prevm
		)
	end
		
	m2 = try(call(remove-keys, dkeys, big-m1))
	expected-result = map(
		1997, '1997',
		1998, '1998',
		1999, '1999'
	)
	call(ASSURE, eq(m2, expected-result), plus('unexpected result', str(m2)))
end

testRepetitionOfPutAndDelOfSameKey = proc()
	repeat-count = 5000 
	
	repeater = func(cnt, prevm)
		while( lt(cnt, repeat-count),
			plus(cnt, 1),
			put(del(prevm, 1), 1, cnt),
			prevm
		)
	end
	
	m = map(1, 2, 3, 4)
	m2 = call(repeater, 1, m)
	call(ASSURE, eq(m2, map(1, minus(repeat-count, 1), 3, 4)), plus('unexpected result', str(m2)))
end

testDell = proc()
	m = map(0.0, 1, 0, 2, 10, 20, 30, 40)
	empty-map = map()
	non-existing-key = 1000

	result = list(
		call(ASSURE, eq(dell(m, non-existing-key), list(false, m)), 'unexpected result'),
		call(ASSURE, eq(dell(empty-map, non-existing-key), list(false, empty-map)), 'unexpected result'),

		call(ASSURE, eq(dell(map(1, 2, 3, 4), 3), list(true, map(1, 2))), 'unexpected result'),
		call(ASSURE, eq(dell(map(1, 2, 3, 4), 1), list(true, map(3, 4))), 'unexpected result'),
		call(ASSURE, eq(dell(map(0.0, 1, 0, 2), 0.0), list(true, map(0, 2))), 'unexpected result'),
		call(ASSURE, eq(dell(map(0.0, 1, 0, 2), 0), list(true, map(0.0, 1))), 'unexpected result'),

		true
	)

	import stdfu
	failed-list = call(stdfu.filter, result, func(x) not(x) end)
	empty(failed-list)
end

testDelUnexistingKey = proc()
	m = map(0.0, 1, 0, 2, 10, 20, 30, 40)
	empty-map = map()
	non-existing-key = 1000

	result = list(
		call(ASSURE, eq(try(del(m, non-existing-key), 'fail'), 'fail'), 'should fail to del non-existing key'),
		call(ASSURE, eq(try(del(empty-map, non-existing-key), 'fail'), 'fail'), 'should fail to del non-existing key'),
		true
	)

	import stdfu
	failed-list = call(stdfu.filter, result, func(x) not(x) end)
	empty(failed-list)
end

testDelAndPutSeveralTimes = proc()
	m1 = map(0.0, 1, 0, 2, 10, 20, 30, 40)
	m2 = del(m1, 0.0)
	m3 = put(m2, 0.0, 'new-1')
	m4 = del(m3, 0.0)
	m5 = put(m4, 0.0, 'new-2')

	result = list(
		call(ASSURE, eq(len(m3), 4), plus('m2 len not ok: ', str(len(m3)))),
		call(ASSURE, eq(len(m4), 3), plus('m4 len not ok: ', str(len(m4)))),
		call(ASSURE, eq(len(m5), 4), plus('m5 len not ok: ', str(len(m5)))),

		call(ASSURE, not(in(m2, 0.0)), plus('m2 should not have 0.0 key: ', str(m2)) ),
		call(ASSURE, not(in(m4, 0.0)), plus('m4 should not have 0.0 key: ', str(m4)) ),
		call(ASSURE, in(m2, 0), plus('m2 should have 0 key: ', str(m2)) ),
		call(ASSURE, in(m4, 0), plus('m4 should have 0 key: ', str(m4)) ),

		call(ASSURE, in(m3, 0.0), plus('m3 should have 0.0 key: ', str(m3)) ),
		call(ASSURE, in(m5, 0.0), plus('m5 should have 0.0 key: ', str(m5)) ),

		call(ASSURE, eq(get(m3, 0.0), 'new-1'), plus('m3 len not ok: ', str(m3)) ),
		call(ASSURE, eq(get(m3, 0), 2), plus('value with key 0 not ok: ', str(m3)) ),

		call(ASSURE, eq(get(m5, 0.0), 'new-2'), plus('m5 len not ok: ', str(m5)) ),
		call(ASSURE, eq(get(m5, 0), 2), plus('value with key 0 not ok: ', str(m5)) ),

		call(ASSURE, eq(try(get(m2, 0.0), 'fail'), 'fail'), 'should fail to access key'),
		call(ASSURE, eq(try(get(m4, 0.0), 'fail'), 'fail'), 'should fail to access key'),
		
		true
	)

	import stdfu
	failed-list = call(stdfu.filter, result, func(x) not(x) end)
	empty(failed-list)
end

testDelFromMapWithSameHashAndOthers = proc()
	m1  = map(0.0, 1, 0, 2, 10, 20, 30, 40)
	m2  = del(m1, 0.0)
	m3  = del(m1, 0)

	result = list(
		call(ASSURE, eq(len(m2), 3), plus('m2 len not ok: ', str(len(m2)))),
		call(ASSURE, eq(len(m3), 3), plus('m3 len not ok: ', str(len(m3)))),
		call(ASSURE, eq(len(m1), 4), plus('m1 len not ok: ', str(len(m1)))),
		true
	)

	import stdfu
	failed-list = call(stdfu.filter, result, func(x) not(x) end)
	empty(failed-list)
end

testDelFromMapWithSameHash = proc()
	m1  = map(0.0, 1, 0, 2)
	m2  = del(m1, 0.0)
	m3  = del(m1, 0)
	m2e = del(m2, 0)
	m3e = del(m3, 0.0)

	result = list(
		call(ASSURE, eq(len(m2), 1), plus('m2 len not ok: ', str(len(m2)))),
		call(ASSURE, eq(len(m3), 1), plus('m3 len not ok: ', str(len(m3)))),
		call(ASSURE, eq(len(m1), 2), plus('m1 len not ok: ', str(len(m1)))),
		call(ASSURE, eq(len(m2e), 0), plus('m2e len not ok: ', str(len(m2e)))),
		call(ASSURE, eq(len(m3e), 0), plus('m3e len not ok: ', str(len(m3e)))),

		call(ASSURE, eq(get(m2, 0), 2), plus('m2 invalid: ', str(m2)) ),
		call(ASSURE, eq(get(m3, 0.0), 1), plus('m3 invalid: ', str(m3)) ),
		call(ASSURE, eq(get(m1, 0), 2), plus('m1 invalid: ', str(m1)) ),
		call(ASSURE, eq(get(m1, 0.0), 1), plus('m1 invalid: ', str(m1)) ),

		call(ASSURE, eq(try(get(m2, 0.0), 'fail'), 'fail'), 'm2 key shouldnt be found' ),
		call(ASSURE, eq(try(get(m3, 0), 'fail'), 'fail'), 'm3 key shouldnt be found' ),
		call(ASSURE, in(m2, 0), 'm2 key should be found' ),
		call(ASSURE, in(m3, 0.0), 'm3 key should be found' ),
		call(ASSURE, not(in(m2, 0.0)), 'm2 key shouldnt be found (in)' ),
		call(ASSURE, not(in(m3, 0)), 'm3 key shouldnt be found (in)' ),

		call(ASSURE, empty(m3e), 'm3e not empty'),
		call(ASSURE, empty(m2e), 'm2e not empty'),

		true
	)

	import stdfu
	failed-list = call(stdfu.filter, result, func(x) not(x) end)
	empty(failed-list)
end

testDelWithBigMap = proc()
	import stdfu
	
	limit = 1000
	nums = call(stdfu.generate, 0, minus(limit, 1), func(i) i end)
	strs = call(stdfu.apply, nums, func(x) str(x) end)	
	big-m1 = call(stdfu.zip, nums, strs)
	big-m2 = del(del(big-m1, 500), 700)

	keys-m1 = keys(big-m1)
	keys-m2 = keys(big-m2)

	vals-m1 = vals(big-m1)
	vals-m2 = vals(big-m2)

	kvs-m1 = keyvals(big-m1)
	kvs-m2 = keyvals(big-m2)
	
	result = list(
		eq(len(big-m2), minus(limit, 2)),
		not(in(big-m2, 500)),
		not(in(big-m2, 700)),
		in(big-m2, 701),
		in(big-m2, 501),
		in(big-m1, 700),
		in(big-m1, 500),

		eq(len(keys-m1), limit),
		eq(len(keys-m2), minus(limit, 2)),
		not(in(keys-m2, 500)),
		not(in(keys-m2, 700)),
		in(keys-m2, 701),
		in(keys-m2, 501),
		in(keys-m1, 700),
		in(keys-m1, 500),

		eq(len(vals-m1), limit),
		eq(len(vals-m2), minus(limit, 2)),
		not(in(vals-m2, '500')),
		not(in(vals-m2, '700')),
		in(vals-m2, '701'),
		in(vals-m2, '501'),
		in(vals-m1, '700'),
		in(vals-m1, '500'),

		eq(len(kvs-m1), limit),
		eq(len(kvs-m2), minus(limit, 2)),
		
		true
	)
	allRight = call(common_test_util.isAllTrueInList, result)
	call(ASSURE, allRight, plus('Unexpected result = ', str(result)))
end

testDelAddSeveralTimes = proc()
	m1 = map(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	m2 = del(m1, 5)
	m3 = put(m2, 5, 'new')
	m4 = del(m3, 5)
	m5 = put(m4, 5, 'even-newer')

	result = list(
		in(m1, 5),
		not(in(m2, 5)),
		in(m3, 5),
		not(in(m4, 5)),
		in(m1, 5),

		eq(get(m1, 5), 6),
		eq(get(m3, 5), 'new'),
		eq(get(m5, 5), 'even-newer'),
		eq(try(get(m2, 5), 'ok'), 'ok'),
		eq(try(get(m4, 5), 'ok'), 'ok'),
		
		true
	)
	allRight = call(common_test_util.isAllTrueInList, result)
	call(ASSURE, allRight, plus('Unexpected result = ', str(result)))
end

testDelOperatorForMap = proc()
	m1      = map(1, 2, 3, 4, 5, 6)
	m2a     = del(m1, 3)
	m2b     = del(m1, 1)
	m2a2    = del(m2a, 5)
	m3      = del(m1, 5)
	m-empty = del(m2a2, 1)

	result = list(
		# m1
		eq(len(m1), 3),
		eq(m1, map(3, 4, 1, 2, 5, 6)),
		not(empty(m1)),

		# m-empty
		empty(m-empty),
		eq(len(m-empty), 0),
		
		# m3
		eq(m3, map(1, 2, 3, 4)),
		eq(len(m3), 2),
		
		# m2a2
		in(m2a2, 1),
		not(in(m2a2, 3)),
		not(in(m2a2, 5)),
		eq(get(m2a2, 1), 2),
		eq(try(eq(get(m2a2, 3), 4), 'ok'), 'ok'),
		eq(try(eq(get(m2a2, 5), 6), 'ok'), 'ok'),
		eq(m2a2, map(1, 2)),
		eq(len(m2a2), 1),
		
		# m2a
		in(m2a, 1),
		in(m2a, 5),
		not(in(m2a, 3)),
		eq(len(m2a), 2),
		not(empty(m2a)),
		eq(get(m2a, 1), 2),
		eq(get(m2a, 5), 6),
		eq(try(eq(get(m2a, 3), 4), 'ok'), 'ok'),
		eq(m2a, map(5, 6, 1, 2)),

		# m2b
		in(m2b, 3),
		in(m2b, 5),
		not(in(m2b, 1)),
		eq(len(m2b), 2),
		not(empty(m2b)),
		eq(get(m2b, 3), 4),
		eq(get(m2b, 5), 6),
		eq(try(eq(get(m2b, 1), 2), 'ok'), 'ok'),
		eq(m2b, map(5, 6, 3, 4)),

		true
	)
	allRight = call(common_test_util.isAllTrueInList, result)
	call(ASSURE, allRight, plus('Unexpected result = ', str(result)))
end

testConcurrentMapAccess = proc()
	fillMap = func(cnt, maxc, m)
		while( lt(cnt, maxc),
			plus(cnt, 1),
			maxc,
			put(m, cnt, str(plus(cnt, 100))),
			m
		)
	end
	
	m = call(fillMap, 0, 10, map())
	ackch = chan()

	tstfib = proc(offset)
		addToMap = func(cnt, maxc, mp)
			while( lt(cnt, maxc),
				plus(cnt, 1),
				maxc,
				put(mp, cnt, str(plus(cnt, 100))),
				mp
			)
		end

		_ = call(addToMap, offset, plus(offset, 20), m)
		_ = send(ackch, 'ok')
		true
	end
	
	_ = spawn(call(tstfib, 200))
	_ = spawn(call(tstfib, 300))
	_ = spawn(call(tstfib, 400))
	_ = spawn(call(tstfib, 500))

	call(proc(n)
		while( lt(n, 4),
			call(
				proc(nn)
					_ = recv(ackch)
					plus(nn, 1)
				end,
				n
			),
			true
		)
	end, 0)
end

testEqForMap = proc()
	m1 = map(1, 2, 3, 4, 5, 6)
	m2 = map(1, 2, 3, 4)
	m3 = m1
	m4 = map(5, 6, 3, 4, 1, 2)
	m5 = map(1, 2, 3, 4, 50, 6)
	m6 = map(1, 2, 3, 4, 5, 60)
	m7 = map(1, 2, 3, 4, 5, 100)

	mx1 = map(1, map(10, 100, 20, 200), 2, 'two')
	mx2 = map(2, 'two', 1, map(10, 100, 20, 200))

	result = list(
		eq(map(), map()),
		eq(m1, m1),
		eq(m1, m3),
		eq(m1, m4),
		eq(m4, m1),
		eq(m1, m4, m3),
		not(eq(m1, map(), m4)),
		not(eq(m1, m2)),
		not(eq(m2, m1)),
		not(eq(m1, m7)),
		not(eq(m7, m1)),
		not(eq(m1, m2, map())),
		not(eq(m1, m5)),
		not(eq(m1, m6)),
		not(eq(list(), map())),
		eq(map(0, 20, 0.0, 10), map(0.0, 10, 0, 20)),
		eq(mx1, mx2),
		true
	)
	allRight = call(common_test_util.isAllTrueInList, result)
	call(ASSURE, allRight, plus('Unexpected result = ', str(result)))
end

testBaseOperatorsForMap = proc()
	m = map(1, 2, 3, 4, 5, 6)
	m2 = map(0.0, 1, 0, 2)
	result = list(
		eq(type(m), 'map'),
		eq(str(map(1,2)), 'map(1 : 2)'),
		eq(conv(map(1,2), 'string'), 'map(1 : 2)'),
		in(m, 1),
		in(m, 3),
		in(m, 5),
		not(in(m, 30)),
		not(in(m, 2)),
		not(in(m, 4)),
		in(m2, 0.0),
		in(m2, 0),
		not(in(m2, 0.1)),
		eq(len(m), 3),
		eq(len(m2), 2),
		eq(len( put(put(m, 10, 11), 100, 110)), 5),
		empty(map()),
		not(empty(m)),
		true
	)
	allRight = call(common_test_util.isAllTrueInList, result)
	call(ASSURE, allRight, plus('Unexpected result = ', str(result)))
end

testKeyvals = proc()
	m = map(1, 2, 3, 4, 5, 6)
	mvals = keyvals(m)
	result = list(
		in(mvals, list(1, 2)),
		in(mvals, list(3, 4)),
		in(mvals, list(5, 6)),
		eq(len(mvals), len(m)),
		true
	)
	allRight = call(common_test_util.isAllTrueInList, result)
	call(ASSURE, allRight, plus('Unexpected result = ', str(result)))
end

testVals = proc()
	m = map(
		10,         'hundred',
		'123',      123,
		list(),     list(0,0),
		'xxx',      123,
		0.5,        'half',
		false,      'not true',
		'some map', map(10, 20, 30, 40),
		'double',   123
	)
	mvals = vals(m)
	result = list(
		in(mvals, 'hundred'),
		in(mvals, list(0,0)),
		in(mvals, 'half'),
		in(mvals, 'not true'),
		eq(len(mvals), len(m)),
		eq(len(find(mvals, 123)), 3),
		in(mvals, map(10, 20, 30, 40)),
		true
	)
	allRight = call(common_test_util.isAllTrueInList, result)
	call(ASSURE, allRight, plus('Unexpected result = ', str(result)))
end

testMapInMap = proc()
	m = map(
		1, 10,
		2, map(1, 'A', 2, 'B'),
		3, 30
	)
	v = get(get(m, 2), 1)
	call(ASSURE, eq(v, 'A'), plus('Wrong value: ', v)) 
end

testKeys = proc()
	m = map(
		10,         'hundred',
		'123',      123,
		list(),     list(0,0),
		list(1, 2), map('a', 'b', 'c', 'd'),
		0.5,        'half',
		false,      'not true'
	)
	mkeys = keys(m)

	result = list(
		in(mkeys, 10),
		in(mkeys, '123'),
		in(mkeys, list()),
		in(mkeys, list(1,2)),
		in(mkeys, 0.5),
		in(mkeys, false),
		eq(len(mkeys), len(m)),
		true
	)
	allRight = call(common_test_util.isAllTrueInList, result)
	call(ASSURE, allRight, plus('Unexpected result = ', str(result)))
end

testGetl = proc()
	someVal = 'some value'
	someKey = 'some key'
	m = map(
		someKey, someVal,
		100,     200,
		true,    false,
		list(),  list(1, 2, 3)
	)
	result = list(
		eq( getl(m, someKey), list(true, someVal) ),
		eq( getl(m, 'not to be found'), list(false false) ),
		true
	)
	allRight = call(common_test_util.isAllTrueInList, result)
	call(ASSURE, allRight, plus('Unexpected result = ', str(result)))
end

testGetWithDefaultWhenKeyNotFound = proc()
	someVal = 'some value'
	someKey = 'some key'
	keyThatIsNotFound = 'this key is not found'
	m = map(
		someKey, someVal,
		100,     200,
		true,    false,
		list(),  list(1, 2, 3)
	)
	defaultVal = 'this is default value when key is not found'
	result = list(
		eq( get(m, keyThatIsNotFound, defaultVal), defaultVal ),
		eq( get(m, keyThatIsNotFound, 'this is default value when key is not found'), defaultVal ),
		eq( get(m, keyThatIsNotFound, call(func(p) p end, defaultVal)), defaultVal ),
		eq( get(m, someKey, defaultVal), someVal ),
		true
	)
	allRight = call(common_test_util.isAllTrueInList, result)
	call(ASSURE, allRight, plus('Unexpected result = ', str(result)))
end

testBoolMap = proc()
	m = map(
		1, 'one',
		false, 'false',
		true, 'true',
		0, 'zero',
		0.0, 'float-zero',
		1.0, 'float-one'
	)
	result = list(
		eq( get(m, 1), 'one' ),
		eq( get(m, 0), 'zero' ),
		eq( get(m, true), 'true' ),
		eq( get(m, false), 'false' ),
		eq( get(m, 1.0), 'float-one' ),
		eq( get(m, 0.0), 'float-zero' ),
		eq( getl(m, 0.0), list(true, 'float-zero') ),
		true
	)
	allRight = call(common_test_util.isAllTrueInList, result)
	call(ASSURE, allRight, plus('Unexpected result = ', str(result)))
end

testCreateMapWithSeveralItems = proc()
	m = map(
		10,  '10',
		0.2, '1/5'
	)

	result = list(
		eq(get(m, 10), '10'),
		eq(get(m, 0.2), '1/5'),
		true
	)
	allRight = call(common_test_util.isAllTrueInList, result)
	call(ASSURE, allRight, plus('Unexpected result = ', str(result)))
end

testMapWithSeveralKeyTypes = proc()
	v = 100
	m = map(
		123, v,
		0.05, v,
		'some', v,
		list(1, list(10, 20), 3), v,
		true, v,
	  not(true), v
	)
	valuesok = list(
		eq(get(m, 123), v),
		eq(get(m, 0.05), v),
		eq(get(m, 'some'), v),
		eq(get(m, list(1, list(10, 20), 3)), v),
		true
	)
	allRight1 = call(common_test_util.isAllTrueInList, valuesok)
	ok1 = call(ASSURE, allRight1, plus('Unexpected result = ', str(valuesok)))

	result = list(
		try(put(m, chan()), true),
		try(put(m, func() 'some' end), true),
		try(put(m, map()), true),
		true
	)
	allRight = call(common_test_util.isAllTrueInList, result)
	ok2 = call(ASSURE, allRight, plus('Unexpected result = ', str(result)))

	and(ok1, ok2)
end

testMapWithSeveralValueTypes = proc()
	txt = 'this is func'
	m1 = map()
	m2 = put(m1, 'func-value', func() txt end)
	m3 = put(m2, 100, list(1, 2, 3))
	m4 = put(m3, 'map-as-value', m2)

	res = call(get(m4, 'func-value'))
	ok1 = call(ASSURE, eq(res, txt), plus('Unexpected result: ', str(res)))

	l = get(m4, 100)
	ok2 = call(ASSURE, eq(l, list(1,2,3)), plus('Unexpected result: ', str(l)))

	mapval = get(m4, 'map-as-value')
	res2 = call(get(mapval, 'func-value'))
	ok3 = call(ASSURE, eq(res2, txt), plus('Unexpected result: ', str(res2)))
	
	and(ok1, ok2, ok3)
end

testMapBasic1 = proc()
	m1 = map()
	m2 = put(put(m1, 10, 'ten'), 20, 'twenty')
	m3 = put(put(m1, 30, 'thirty'), '40', 'forty')
	m4 = put(put(m3, '50', 'fifty'), 60, 'sixty')

	m5 = put(m4, 100, '100')
	m6 = put(m4, 100, 'hundred')
		
	result = list(
		eq( get(m4, plus('4', '0')), 'forty' ),
		eq( get(m2, 10), 'ten' ),
		eq( get(m2, 20), 'twenty' ),
		eq( get(m3, 30), 'thirty' ),
		eq( get(m3, '40'), 'forty' ),
		eq( get(m4, '50'), 'fifty' ),
		eq( get(m4, 60), 'sixty' ),
		eq( get(m4, 30), 'thirty' ),
		eq( get(m4, '40'), 'forty' ),

		eq( get(m5, 100), '100' ),
		eq( get(m6, 100), 'hundred' ),
		true
	)
	allRight = call(common_test_util.isAllTrueInList, result)
	allFoundThatShould = call(ASSURE, allRight, plus('Unexpected result = ', str(result)))

	result2 = list(
		try(get(m1, 10), true),
		try(get(m1, 30), true),
		try(get(m1, '50'), true),
		try(get(m1, 60), true),
		try(get(m2, 30), true),
		try(get(m2, '40'), true),
		try(get(m2, 60), true),
		try(get(m2, '50'), true),
		try(get(m3, '50'), true),
		try(get(m3, 60), true),
		try(get(m3, 10000), true),
		
		try(put(m2, 10, 'anything'), true), # should make RTE as key already exists
		true
	)

	allRight2 = call(common_test_util.isAllTrueInList, result2)
	allNotFoundThatShouldnt = call(ASSURE, allRight2, plus('Unexpected result = ', str(result2)))

	and(allFoundThatShould, allNotFoundThatShouldnt)
end

endns

