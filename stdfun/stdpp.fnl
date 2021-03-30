
ns stdpp

pprint = proc(val)
	import stdio

	value-str = call(pform val)
	_ = call(stdio.printline value-str)
	val
end

form = func(val)
	import stdfu

	indent-mark = if( gt(len(argslist()) 1)
		call(func()
			_ marker = argslist():
			if( eq(type(marker) 'string') marker '\t')
		end argslist():)
		'\t' # tab is default
	)

	print-list = func(v indent)
		call(stdfu.foreach v func(x cum) call(print-value x plus(indent indent-mark) cum) end '')
	end

	print-map = func(v indent)
		call(stdfu.foreach keyvals(v) func(kv cum)
			mkey mvalue = kv:
			cum2 = call(print-value mkey plus(indent indent-mark) cum)
			call(print-value mvalue plus(indent indent-mark) cum2)
		end '')
	end

	print-value = func(v indent res)
		item-str = case( type(v)
			'string' plus('\n' indent '\'' v '\'')
			'list'   plus('\n' indent plus('list(' call(print-list v indent) '\n' indent ')') )
			'map'    plus('\n' indent plus('map(' call(print-map v indent) '\n' indent ')') )
			plus('\n' indent str(v))
		)
		plus(res item-str)
	end

	call(print-value val '' '')
end

pform = proc(val)
	import stdfu

	indent-mark = if( gt(len(argslist()) 1)
		call(proc()
			_ marker = argslist():
			if( eq(type(marker) 'string') marker '\t')
		end argslist():)
		'\t' # tab is default
	)

	print-list = proc(v indent)
		call(stdfu.ploop proc(x cum) call(print-value x plus(indent indent-mark) cum) end v '')
	end

	print-map = proc(v indent)
		call(stdfu.ploop proc(kv cum)
			mkey mvalue = kv:
			cum2 = call(print-value mkey plus(indent indent-mark) cum)
			call(print-value mvalue plus(indent indent-mark) cum2)
		end keyvals(v) '')
	end

	print-value = proc(v indent res)
		item-str = case( type(v)
			'string' plus('\n' indent '\'' v '\'')
			'list'   plus('\n' indent plus('list(' call(print-list v indent) '\n' indent ')') )
			'map'    plus('\n' indent plus('map(' call(print-map v indent) '\n' indent ')') )
			plus('\n' indent str(v))
		)
		plus(res item-str)
	end

	call(print-value val '' '')
end

endns

