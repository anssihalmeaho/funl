
ns stdvar_test

import ut_fwk

ASSURE = ut_fwk.VERIFY

import common_test_util
import stdvar

test-change-v2-ok = proc()
	var = call(stdvar.new 50)
	retv = call(stdvar.change-v2 var func(prev inp) list(plus(prev inp) list(prev inp)) end 10)
	_ = call(ASSURE eq(retv list(true '' 60 list(50 10))) plus('Unexpected result = ' str(retv)))
	call(ASSURE eq(call(stdvar.value var) 60) 'unexpected value')
end

test-change-v2-ok-no-extra-parameter = proc()
	var = call(stdvar.new 50)
	retv = call(stdvar.change-v2 var func(prev) list(plus(prev 1) list(prev 1)) end)
	_ = call(ASSURE eq(retv list(true '' 51 list(50 1))) plus('Unexpected result = ' str(retv)))
	call(ASSURE eq(call(stdvar.value var) 51) 'unexpected value')
end

test-change-v2-RTE-in-func = proc()
	var = call(stdvar.new 50)
	retv = call(stdvar.change-v2 var func(prev inp) list(plus(prev inp) list(prev inp)) end true)
	_ = call(ASSURE eq(retv list(false 'mismatching types as arguments' '' '')) plus('Unexpected result = ' str(retv)))
	call(ASSURE eq(call(stdvar.value var) 50) 'unexpected value')
end

test-change-v2-ret-value-not-list = proc()
	var = call(stdvar.new 50)
	retv = call(stdvar.change-v2 var func(prev inp) 'crappy return value' end true)
	_ = call(ASSURE eq(retv list(false 'List value expected' '' '')) plus('Unexpected result = ' str(retv)))
	call(ASSURE eq(call(stdvar.value var) 50) 'unexpected value')
end

test-change-v2-ret-value-is-empty-list = proc()
	var = call(stdvar.new 50)
	retv = call(stdvar.change-v2 var func(prev inp) list() end true)
	_ = call(ASSURE eq(retv list(false 'Too short list received (empty)' '' '')) plus('Unexpected result = ' str(retv)))
	call(ASSURE eq(call(stdvar.value var) 50) 'unexpected value')
end

test-change-v2-ret-value-list-too-short = proc()
	var = call(stdvar.new 50)
	retv = call(stdvar.change-v2 var func(prev inp) list(1) end true)
	_ = call(ASSURE eq(retv list(false 'Too short list received (one item)' '' '')) plus('Unexpected result = ' str(retv)))
	call(ASSURE eq(call(stdvar.value var) 50) 'unexpected value')
end

endns
