
ns expansion_test

import ut_fwk

ASSURE = ut_fwk.VERIFY

import common_test_util

returnArguments = func(p1 _ p3)
	list(
		list(p1 p3)
		argslist()
	)
end

returnArguments2 = func()
	argslist()
end

variadicF = func()
	list(argslist():)
end

returnArguments3 = func()
	call(variadicF, argslist():)
end

testArgslist = proc()
	arg-val-1 = 'some test value 1'
	dum-val   = 'dum'
	arg-val-2 = 'some test value 2'

	l1 l2 = call(returnArguments arg-val-1 dum-val arg-val-2):
	p1 p2 = l1:

	and(
		eq(p1 arg-val-1)
		eq(p2 arg-val-2)
	)
end 

testArgslistWithSubFunction = proc()
	arg-val-1 = 'some test value 1'
	dum-val   = 'dum'
	arg-val-2 = 'some test value 2'

	subReturnArguments = func(p1 _ p3)
		list(
			list(p1 p3)
			argslist()
		)
	end

	l1 l2 = call(subReturnArguments arg-val-1 dum-val arg-val-2):
	p1 p2 = l1:

	and(
		eq(p1 arg-val-1)
		eq(p2 arg-val-2)
	)
end 

testArgslistAsVariadicFunction = proc()
	argl = list(1 2 3)
	l = call(returnArguments2 argl:)
	eq(l argl)
end 

testArgslistAsVariadicFunction2 = proc()
	argl = list(1 2 3)
	l = call(returnArguments3 argl:)
	eq(l argl)
end 

testArgslistWithNoArguments = proc()
	l = call(returnArguments2)
	eq(l list())
end 

testExpandOperCall = proc()
	l1 = list(2 4 1 3)
	l2 = list('a' 'b' 'c' 'd')
	l3 = list('e' 'f' 'g' 'h')
	l4 = list()
	l5 = list(10)
	and(
		eq(plus(l1:) 10)
		eq(plus(l2:) 'abcd')
		eq(plus(l2: l3:) 'abcdefgh')
		eq(plus('<'  l2: '_' l3: '>') '<abcd_efgh>')
		eq(plus('<'  l2: '_' l3: l4: '>') '<abcd_efgh>')
		eq( list(list():) list() )
		eq( str(l5:) '10' )
	)
end

testExpandLetDefInFunction = proc()
	l = list(1 2 3)
	a b c = l:
	x _ z = l:
	v = list(10):

	and(
		eq(a 1)
		eq(b 2)
		eq(c 3)
		eq(x 1)
		eq(z 3)
		eq(v 10)
	)
end

G_l = list(1 2 3)
G_a G_b G_c = G_l:
G_x _ G_z = G_l:
G_v = list(10):

testExpandLetDefInNamespace = proc()
	and(
		eq(G_a 1)
		eq(G_b 2)
		eq(G_c 3)
		eq(G_x 1)
		eq(G_z 3)
		eq(G_v 10)
	)
end

endns
