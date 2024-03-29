
ns stdfu

select-keys = func(m keylist)
	keylooper = func(remaining-keys result)
		if( empty(remaining-keys)
			result
			call(func()
				nextkey = head(remaining-keys)
				found val = getl(m nextkey):
				next-result = if(found put(result nextkey val) result)
				call(keylooper rest(remaining-keys) next-result)
			end)
		)
	end

	call(keylooper keylist map())
end

pre-decorate = func(pre-handler handler)
	func()
		new-args = call(pre-handler argslist():)
		call(handler new-args:)
	end
end

post-decorate = func(handler post-handler)
	func()
		retval = call(handler argslist():)
		new-retval = call(post-handler retval argslist():)
		new-retval
	end
end

p-pre-decorate = func(pre-handler handler)
	proc()
		new-args = call(pre-handler argslist():)
		call(handler new-args:)
	end
end

p-post-decorate = func(handler post-handler)
	proc()
		retval = call(handler argslist():)
		new-retval = call(post-handler retval argslist():)
		new-retval
	end
end

loop = func(handler inlist result)
	while( not(empty(inlist))
		handler
		rest(inlist)
		call(handler head(inlist) result)
		result
	)
end

ploop = proc(handler inlist result)
	while( not(empty(inlist))
		handler
		rest(inlist)
		call(handler head(inlist) result)
		result
	)
end

chain = func(input, handler-list)
	applyHandler = func(handler-chain, output)
		while( not(empty(handler-chain)),
			rest(handler-chain), 
			call(head(handler-chain), output),
			output
		)
	end
	
	call(applyHandler, handler-list, input)
end

foreach = func(lst, handler, initial)
	looper = func(l, cum)
		while( not(empty(l)),
			rest(l),
			call(handler, head(l), cum),
			cum
		)
	end

	if( eq(type(lst), 'list'),
		call(looper, lst, initial),
		error('foreach needs list as argument, not ', type(lst))
	)	
end

max = func(lst greater-than-func)
	comparator = func(item-1 item-2)
		call(greater-than-func item-1 item-2)
	end

	if( eq(type(lst) 'list')
		case( len(lst)
			0 error('empty list')
			1 head(lst)
			call(foreach rest(lst) comparator head(lst))
		)
		error(sprintf('argument not a list (%s)' type(lst)))
	)
end

zip = func(keylist, valuelist)
	looper = func(keyl, values, m)
		while( not(empty(keyl)),
			rest(keyl),
			rest(values),
			if( in(m, head(keyl)),
				m,
				put(m, head(keyl), head(values))
			),
			m
		)
	end
	
	if( eq(len(keylist), len(valuelist)),
		call(looper, keylist, valuelist, map()),
		error('zip needs list to be same length')
	)	
end

generate = func(startv, stopv, generFunc)
	looper = func(i, targetlist)
		while(le(i, stopv),
			plus(i, 1),
			call(func(idx)
				newv = call(generFunc, idx)
				append(targetlist, newv)
			end, i),
			targetlist
		)
	end

	ok1 = eq(type(startv), 'int')
	ok2 = eq(type(stopv), 'int')
	case(list(ok1, ok2),
		list(true, false), error('2nd argument not int'),
		list(false, true), error('1st argument not int'),
		list(false, false), error('1st and 2nd argument are not int'),
		list(true, true), call(looper, startv, list())
	)
end

apply = func(srclist, converter)
	applyConv = func(l, newl, cnt)
		if( empty(l),
			newl,
			call(
				func()
					remainingList = rest(l)
					item = head(l)
					convertedItem = call(converter, item)
					call(applyConv, remainingList, append(newl, convertedItem), plus(cnt,1))
				end
			)
		)
	end

	call(applyConv, srclist, list(), 1)
end

proc-apply = proc(srclist, converter)
	applyConv = proc(l, newl, cnt)
		if( empty(l),
			newl,
			call(
				proc()
					remainingList = rest(l)
					item = head(l)
					convertedItem = call(converter, item)
					call(applyConv, remainingList, append(newl, convertedItem), plus(cnt,1))
				end
			)
		)
	end

	call(applyConv, srclist, list(), 1)
end

filter = func(srcdata, condition)
	map-filter = func(src-map)
		looper = func(kvs, resultm)
			next-kv = if(not(empty(kvs)), head(kvs), 'whatever')
	
			while( not(empty(kvs)),
				rest(kvs),
				if( call(condition, head(next-kv), last(next-kv)), 
					put(resultm, head(next-kv), last(next-kv)),
					resultm
				),
				resultm
			)
		end
		
		call(looper, keyvals(src-map), map())
	end
	
	list-filter = func(srclist)
		handleNext = func(l, newl)
			applyCond = func()
				item = head(l)
				remainingList = rest(l)
	
				if( call(condition, item),
					call(handleNext, remainingList, append(newl, item)),
					call(handleNext, remainingList, newl)
				)
			end
	
			if(	empty(l),
				newl,
				call(applyCond)
			)
		end
	
		call(handleNext, srclist, list())
	end

	case( type(srcdata),
		'list', call(list-filter, srcdata),
		'map',  call(map-filter, srcdata),
		error('non-supported type: ', type(srcdata))
	)
end

# true if condition(item) is true for all, false otherwise
applies-for-all = func(srclist, condition)
	result-list = call(filter, srclist, condition)
	eq(len(result-list), len(srclist))
end

# true if condition(item) is true for any, false otherwise
applies-for-any = func(srclist, condition)
	result-list = call(filter, srclist, condition)
	not(eq(result-list, list()))
end

group-by = func(srcdata, grouper)
	append-to-list = func(resm akey aval)
		prevl = get(resm akey)
		newl = append(prevl aval)
		newm = del(resm akey)
		put(newm akey newl)
	end
	
	map-group-by = func(src-map)
		looper = func(kvs result)
			next-kv = if(not(empty(kvs)) head(kvs) 'whatever')
	
			while( not(empty(kvs))
				rest(kvs)
				call(func()
					kv vv = head(kvs):
					key value = call(grouper kv vv):
					if( in(result key)
						call(append-to-list result key value)
						put(result key list(value))
					)
				end)
				result
			)
		end
		
		call(looper, keyvals(src-map), map())
	end
	
	list-group-by = func(srclist)
		looper = func(kvl result)
			while( not(empty(kvl))
				rest(kvl)
				call(func()
					key value = call(grouper head(kvl)):
					if( in(result key)
						call(append-to-list result key value)
						put(result key list(value))
					)
				end)
				result
			)
		end
		
		call(looper srclist map())
	end

	case( type(srcdata),
		'list', call(list-group-by, srcdata),
		'map',  call(map-group-by, srcdata),
		error('non-supported type: ', type(srcdata))
	)
end

merge = func()
	args = argslist()
	is-conflict-handler-given = eq(type(head(args)) 'function')
	offset = if(is-conflict-handler-given 1 0)

	conflict-handler = if( is-conflict-handler-given
		head(args)
		func(key val1 val2)
			list(false val2)
		end
	)

	looper = func(kvs result)
		while( not(empty(kvs))
			rest(kvs)
			call(func()
				key value = head(kvs):
				is-key-in = in(result key)
				do-add chosen-val = if( is-key-in
					call(conflict-handler key get(result key) value)
					list(true value)
				):
				if(do-add 
					if(is-key-in
						put(del(result key) key chosen-val)
						put(result key chosen-val)
					)
					result
				)
			end)
			result
		)
	end
	
	map-looper = func(mlist result)
		while( not(empty(mlist))
			rrest(mlist)
			call(looper keyvals(last(mlist)) result)
			result
		)
	end
	
	call(map-looper slice(args offset) map())
end

pipe = func(func-list)
	func(input)
		looper = func(flist inp)
			while( not(empty(flist))
				rest(flist)
				call(head(flist) inp)
				inp
			)
		end

		call(looper func-list input)
	end
end

proc-pipe = func(proc-list)
	proc(input)
		looper = proc(flist inp)
			while( not(empty(flist))
				rest(flist)
				call(head(flist) inp)
				inp
			)
		end

		call(looper proc-list input)
	end
end

pairs-to-map = func(kv-list)
	pair-to-map = func(item result)
		put(result head(item) last(item))
	end

	call(loop pair-to-map kv-list map())
end

# overwrites key-value if found, ignored if not found
write-if-found = func(src delta)
	add-to = func(result kvs)
		if(empty(kvs)
			result
			call(func()
				key val = head(kvs):
				next = if(in(result key)
					put(del(result key) key val)
					result
				)
				call(add-to next rest(kvs))
			end)
		)
	end

	call(add-to src keyvals(delta))
end

# writes key-values regardless is it found or not
write-all = func(src delta)
	add-to = func(result kvs)
		if(empty(kvs)
			result
			call(func()
				key val = head(kvs):
				next = if(in(result key)
					put(del(result key) key val)
					put(result key val)
				)
				call(add-to next rest(kvs))
			end)
		)
	end

	call(add-to src keyvals(delta))
end

# write key-values which are not yet found
write-if-not-found = func(src delta)
	add-to = func(result kvs)
		if(empty(kvs)
			result
			call(func()
				key val = head(kvs):
				next = if(in(result key)
					result
					put(result key val)
				)
				call(add-to next rest(kvs))
			end)
		)
	end

	call(add-to src keyvals(delta))
end

endns
