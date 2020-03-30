
ns func_proc_call_test

import ut_fwk

defaultText = 'ok'
p2 = 'xxx'

functionToTest = func(p1, p2x, p3)
	retval = defaultText
	plus(retval, ':', str(p1), ':', p2x, ':', str(p3))
end

testFunctionCall = func()
	result = call(functionToTest, 1, '2', true)
	call(ut_fwk.VERIFY, eq(result, 'ok:1:2:true'), plus('Unexpected result = ', str(result)))
end

testAnonymousFuncCall = func()
	p3 = 'yep'
	result = call(func() v = p3 v end)
	call(ut_fwk.VERIFY, eq(result, 'yep'), plus('Unexpected result = ', str(result)))
end

testAnonymousFuncCall2 = func()
	p3 = 'yep'
	retf = func(p1)
		func() v = plus(p1, p3) v end
	end
	f = call(retf, '...')
	result = call(f)
	call(ut_fwk.VERIFY, eq(result, '...yep'), plus('Unexpected result = ', str(result)))
end

testClosureAndSubfunction = func()
	v11 = 'should not be used'
	subf = func(p11)
		v111 = p11
		subsubf = func(p111)
			call(functionToTest, v111, plus(p111, p2), not(true))
		end

		subsubf
	end

	clos = call(subf, 100)
	result = call(clos, 'any')
	call(ut_fwk.VERIFY, eq(result, 'ok:100:anyxxx:false'), plus('Unexpected result = ', str(result)))
end

testNotEnoughArguments = proc()
	assumed = 'expected error'
	result = try(call(functionToTest, 1, '2'), assumed)
	call(ut_fwk.VERIFY, eq(result, assumed), plus('Unexpected result = ', str(result)))
end

testNonExistingFunctionCalled = proc()
	assumed = 'expected error'
	result = try(call(nonExistFunction), assumed)
	call(ut_fwk.VERIFY, eq(result, assumed), plus('Unexpected result = ', str(result)))
end

testNonFunctionCalled = proc()
	assumed = 'expected error'
	result = try(call('this is not function'), assumed)
	call(ut_fwk.VERIFY, eq(result, assumed), plus('Unexpected result = ', str(result)))
end

testFunctionCallingProcedure = proc()
	someProcedure = proc() 'hello' 	end
	middleFunction = func() call(someProcedure) end

	assumed = 'expected error'
	result = try(call(middleFunction), assumed)
	call(ut_fwk.VERIFY, eq(result, assumed), plus('Unexpected result = ', str(result)))
end

/*
testImportFailing = proc()
	testf = func()
		import non_existent_module
		'you shouldnt see this'
	end

	assumed = 'expected error'
	result = try(call(testf), assumed)
	call(ut_fwk.VERIFY, eq(result, assumed), plus('Unexpected result = ', str(result)))
end
*/

testWhile = func()
	subf = func(n, l)
		nexti = minus(n, 1)

		while(
			not(eq(n, 0)),
			nexti,
			append(l, n),
			l
		)
	end

	result = call(subf, 5, list())
	assumed = list(5, 4, 3, 2, 1)
	call(ut_fwk.VERIFY, eq(result, assumed), plus('Unexpected result = ', str(result)))
end

testWhile2 = func()
	getInnerF = func()
		func(p) minus(p, 1) end
	end

	subf = func(n, l)
		innerF = call(getInnerF)
		nexti = call(innerF, n)

		while(
			not(eq(n, 0)),
			nexti,
			append(l, n),
			l
		)
	end

	result = call(subf, 5, list())
	assumed = list(5, 4, 3, 2, 1)
	call(ut_fwk.VERIFY, eq(result, assumed), plus('Unexpected result = ', str(result)))
end

testWhile3 = func()
	repeat = func(x, n, l)
		counter = minus(n, 1)
		while(
			not(eq(counter, 0)),
			x,
			counter,
			append(l, mul(x, 10)),
			l
		)
	end

	subf = func(n, l)
		nexti = minus(n, 1)

		while(
			not(eq(n, 0)),
			nexti,
			call(repeat, n, 4, l),
			l
		)
	end

	result = call(subf, 5, list())
	assumed = list(50, 50, 50, 40, 40, 40, 30, 30, 30, 20, 20, 20, 10, 10, 10)
	call(ut_fwk.VERIFY, eq(result, assumed), plus('Unexpected result = ', str(result)))
end

endns
