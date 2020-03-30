
ns main

# reads help for all operators to list
get-all-helps = func()
	looper = func(opelist result)
		while( not(empty(opelist))
			rest(opelist)
			append(result help(head(opelist)))
			result
		)
	end
	
	call(looper help('operators') list())
end

# writes list items to file
write-to-file = proc(help-list)
	looper = proc(lst)
		while( not(empty(lst))
			call(proc()
				_ = call(stdfiles.writeln fh head(lst))
				rest(lst)
			end)
			true
		)		
	end
	
	import stdfiles
	fh = call(stdfiles.create 'helpit.txt')
	_ = call(looper help-list)
	call(stdfiles.close fh)
end

# main
main = proc()
	just-oper-names = if( not(empty(argslist()))
		eq(head(argslist()) 'opnames')
		false
	)
	list-to-write = if( just-oper-names
		help('operators')
		call(get-all-helps)
	)
	call(write-to-file list-to-write)	
	
	#call(write-to-file call(get-all-helps))	
end

endns
