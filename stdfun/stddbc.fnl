
ns stddbc

bt-to-str = func(callchain)
	loopy = func(rest-of result)
		if(empty(rest-of)
			result
			call(func()
				next-str = call(print-citem head(rest-of))
				call(loopy rest(rest-of) plus(result next-str))
			end)
		)
	end

	call(loopy callchain '')
end

print-citem = func(citem)
	sprintf(
		'  file: %s line: %d args: %v\n'
		get(citem 'file')
		get(citem 'line')
		get(citem 'args')
	)
end

assert = call(proc()
	import stdos
	import stdrun

	_ env-value = call(stdos.getenv 'FUNL_DBC_CC'):

	func(condition errtext)
		condtype = type(condition)
		condval = case( condtype
			'bool' 	   condition
			'function' call(condition)
			true
		)
		bt = call(stdrun.backtrace)
		prev = if( lt(len(bt) 2)
			map('file' '-' 'line' 0 'args' list())
			head(rest(bt))
		)
		printing = case(env-value
			'nocc'
			errtext

			'prev'
			sprintf('%s:\nfrom: %s' errtext call(print-citem prev))

			sprintf('%s:\ncall chain:\n%s' errtext call(bt-to-str bt))
		)
		if(condval true error(printing))
	end
end)

endns
