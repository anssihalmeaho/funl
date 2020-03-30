
ns main

PRINT-OVERALL = 1
PRINT-MOD = 2
PRINT-ALL = 3

import stdio
import stdfilu
import stdfu
import stdstr

get-files = proc(target-dir)
 	result = call(stdfilu.get-files-by-ext target-dir 'fnl')
	is-error = eq(type(result) 'string')
	_ = if(is-error call(stdio.printf 'error in reading: %s\n' result) '')
	if( is-error
		list(false list())
		list(true result)
	)
end

get-list-of-mods = func(file-list)
	fl1 = call(stdfu.apply file-list func(fit) head(split(fit '.fnl')) end)
	fl2 = call(stdfu.filter fl1 func(fit) call(stdstr.endswith fit '_test') end)
	list(true fl2)
end

get-tests-from-mod = proc(testmod)
	modmap = try(eval(sprintf('imp(%s)' testmod)))
	import-nok = if( eq(type(modmap) 'string') print('Import failed: ' modmap) false)
	if( import-nok
		map()
		call(proc()
			test-proc-filter = func(proc-name _) 
				or( 
					call(stdstr.startswith proc-name 'test')
					call(stdstr.startswith proc-name 'Test')
				)
			end
			call(stdfu.filter modmap test-proc-filter)
		end)
	)
end

get-test-list-for-mod = proc(target-dir mod-name)
	ok file-list = call(get-files target-dir):

	_ modlist = call(get-list-of-mods file-list):
	target-mod-list = call(stdfu.filter modlist func(mname) eq(mname mod-name) end)
	mod-found = not(empty(target-mod-list))
	tests-from-mod = if( mod-found 
		call(get-tests-from-mod head(target-mod-list)) 
		map()
	)
	list(mod-found tests-from-mod)
end

get-all-tests = proc(target-dir)
	ok file-list = call(get-files target-dir):
	#_ = print(file-list)

	gather-from-one-mod = proc(onemod result-map)
		looper = func(kvlist resm)
			while( not(empty(kvlist))
				rest(kvlist)
				put(resm plus(onemod ' : ' head(head(kvlist))) last(head(kvlist)))
				resm
			)
		end
		
		testmap = call(get-tests-from-mod onemod)
		call(looper keyvals(testmap) result-map)
	end
		
	gather-from-mods = proc(modlist)
		looper = proc(mlist result-map)
			while( not(empty(mlist))
				rest(mlist)
				call(gather-from-one-mod head(mlist) result-map)
				result-map
			)
		end
		
		list(true call(looper modlist map()))
	end
	
	get-list-of-tests = proc()
		mod-ok modlist = call(get-list-of-mods file-list):
		if( mod-ok
			call(gather-from-mods modlist)
			list(false list())
		)
	end
	
	if( ok
		call(get-list-of-tests)
		list(false list())
	)
end

get-test-list-for-test = proc(target-dir tcase)
	ok test-case-map = call(get-all-tests target-dir):
	matching-map = call(stdfu.filter test-case-map 
		func(tc-name _) 
			tc-part = call(stdstr.strip last(split(tc-name ':')))
			eq(tc-part tcase) 
		end
	)
	list(
		and(ok not(empty(matching-map)))
		matching-map
	)
end

main = proc()
	args = argslist()
	target-defined target-dir = if( and( not(empty(args)) eq(type(head(args)) 'string') )
		list(true head(args))
		list(true '.')
	):
	options-map = if( and( gt(len(args) 1) eq(type(ind(args 1)) 'map'))
		ind(args 1)
		map()
	)
	is-mod-defined = in(options-map 'mod')
	is-test-defined = in(options-map 'test')
	are-all-executed = not(or(is-mod-defined is-test-defined))
	found print-val = getl(options-map 'print'):
	print-opt = case( if(found print-val 'overall')
		'all'     PRINT-ALL
		'mod'     PRINT-MOD
		'overall' PRINT-OVERALL
		PRINT-OVERALL
	)
	
	run-tests = proc()
		ok test-case-map = cond(
			and(is-test-defined is-mod-defined) 
			call(get-test-list-for-test-and-mod target-dir get(options-map 'test') get(options-map 'mod'))
			
			is-test-defined 
			call(proc()
				tcase = get(options-map 'test')
				tc-found tcm = call(get-test-list-for-test target-dir tcase):
				_ = if(not(tc-found) 
					call(stdio.printf 'Test case not found (%s)\n' tcase) 
					'')
				list(tc-found tcm)
			end)
			
			is-mod-defined
			call(proc()
				tmod = get(options-map 'mod')
				mod-found tcm = call(get-test-list-for-mod target-dir tmod):
				_ = if(not(mod-found) 
					call(stdio.printf 'Test module not found (%s)\n' tmod) 
					'')
				list(mod-found tcm)
			end)
			
			call(get-all-tests target-dir)
		):
		
		run-one-test = proc(test-name test-proc)
			retv = try(call(test-proc))
			rte-text = if( eq(type(retv) 'string') sprintf('(%s)' retv) '')
			tc-result = cond(
				eq(type(retv) 'string') false
				retv                    true
				false
			)
			_ = if( in(list(PRINT-ALL PRINT-MOD) print-opt)
				call(stdio.printf '%s %s : %s %s\n' if(tc-result '---' '<<<') if(tc-result 'PASS' 'FAIL') test-name rte-text)
				''
			)
			tc-result
		end

		run-cases = proc(tc-kvs)
			looper = proc(tckvs cnt-list)
				while( not(empty(tckvs))
					rest(tckvs)
					call(proc()
						tcname tcproc = head(tckvs):
						tc-passed = call(run-one-test tcname tcproc)
						pass-count fail-count = cnt-list:
						if( tc-passed
							list(plus(pass-count 1) fail-count)
							list(pass-count plus(fail-count 1))
						)
					end)
					cnt-list
				)
			end
			
			passed failed = call(looper tc-kvs list(0 0)):
			all-passed = eq(failed 0)
			_ = if( all-passed
				call(stdio.printf 'PASSED (%d)\n' passed)
				call(stdio.printf 'FAILED (passed=%d)(failed=%d)\n' passed failed)
			)
			all-passed
		end
		
		if( ok
			call(run-cases keyvals(test-case-map))
			false
		)
	end
	
	if( target-defined
		call(run-tests)
		call(proc() _ = call(stdio.printf 'no target directory defined\n') false end)
	)
end

all = proc()
	call(main list(argslist(): map('print' 'all')):)
end

endns
