package funl 

func init() {

	stdfunMap["repl"] = `
ns main

import stdio
import stdstr

# Note. some symbols are named with prefix ___
#       so that clash between user let-definitions would be less probable

___handleHelp = proc()
	_ = call(stdio.printline 'Input can be:')
	_ = call(stdio.printline '  help         -> prints this help')
	_ = call(stdio.printline '  ?            -> prints this help')
	_ = call(stdio.printline '  quit         -> exits repl')
	_ = call(stdio.printline '  exit         -> exits repl')
	_ = call(stdio.printline '  <expression> -> evaluates expression and prints result')
	_ = call(stdio.printline '')
	_ = call(stdio.printline 'Adding < to end of line causes repl to gather more input before evaluation')
	_ = call(stdio.printline '')
	true
end

___strip-last = func(prev)
	case( len(prev)
		0 error('odd input')
		1 ''
		slice(prev 0 minus(len(prev) 2))
	)
end

___getmore = proc(prev-input)
	_ = call(stdio.printout 'funl>... ')
	more-input = plus(prev-input call(stdio.readinput))

	if( call(stdstr.endswith more-input '<')
		call(___getmore call(___strip-last more-input))
		if( call(stdstr.endswith more-input '<')
			call(___strip-last more-input)
			more-input
		)
	)
end

___repl = proc()
	_ = call(stdio.printout 'funl> ')
	___input = call(stdio.readinput)
	___continue = not(in(list('quit' 'exit') ___input))

	___real-input = if( call(stdstr.endswith ___input '<')
		call(___getmore call(___strip-last ___input))
		___input
	)
	#_ = print(':' real-input ':')

	result = case( ___input
		''     true
		'quit' false
		'exit' false
		'help' call(___handleHelp)
		'?'    call(___handleHelp)
		call(stdio.printline try(eval(___real-input)))
	)

	while(___continue 'done')
end

main = proc()
	_ = call(stdio.printline 'Welcome to FunL REPL (interactive command shell)')
	call(___repl)
end

endns
`

	stdfunMap["stddbc"] = `
ns stddbc

assert = func(condition, errtext)
	condtype = type(condition)
	condval = case( condtype,
		'bool', 	condition,
		'function', call(condition),
		true
	)

	if(condval, true, error(errtext))
end

endns
`

	stdfunMap["stdfilu"] = `
ns stdfilu

import stdfiles
import stdfu

get-files-by-ext = proc(path, extension)
	matcher = func(filename)
		in(filename, plus('.', extension))
	end

	result = call(stdfiles.read-dir, path)
	if( eq(type(result), 'string'),
		result,
		call(stdfu.filter, keys(result), matcher)
	)
end 

get-subdirs = proc(path)
	matcher = proc(v)
		finfo = call(stdfiles.finfo-map v)
		and(
			in(finfo 'is-dir')
			get(finfo 'is-dir')
		)
	end

	looper = proc(kvs results)
		while( not(empty(kvs))
			rest(kvs)
			call(proc()
				finfo = head(kvs)
				fname fm = finfo:
				if( call(matcher fm)
					append(results fname)
					results
				)
			end)
			results
		)
	end
		
	fmap = call(stdfiles.read-dir path)
	call(looper keyvals(fmap) list())
end 

get-nondirs = proc(path)
	matcher = proc(v)
		finfo = call(stdfiles.finfo-map v)
		and(
			in(finfo 'is-dir')
			not(get(finfo 'is-dir'))
		)
	end

	looper = proc(kvs results)
		while( not(empty(kvs))
			rest(kvs)
			call(proc()
				finfo = head(kvs)
				fname fm = finfo:
				if( call(matcher fm)
					append(results fname)
					results
				)
			end)
			results
		)
	end
		
	fmap = call(stdfiles.read-dir path)
	call(looper keyvals(fmap) list())
end 

endns
`

	stdfunMap["stdfu"] = `
ns stdfu

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

endns
`

	stdfunMap["stdmeta"] = `
ns stdmeta

/*
todo

- list('list' list-schema)
    -> validates with list schema

- list('exact-keys' 'f1' 'f2')
    -> should have exactly following keys (not less, not more)

- list('len' 5)
     -> length of list or map (key-values) checked
	 -> ranges too ?

- list('all-in-list' list('type' 'int'))
     -> requirement(s) for all items in list
	 -> could it be also for all key-values in map ?

- validation with customized function(s)...?

*/

# gathers documentation from schema
get-doc = func(schema)
	import stdfu

	visit-map = func(map-schema indent result-str)
		map-fields-visit = func(kvitem output)
			get-field-visitor = func(mkey)
				func(check out-str)
					fout = case( head(check)
						'required'
							'required, '
						'type'
							plus('type: ' head(rest(check)) ', ')
						'map'
							call(visit-map head(rest(check)) plus(indent '    ') 'map: ')
						'doc'
							call(func()
								doc-lines = call(stdfu.apply rest(check) func(item) plus('\n' indent '  -> ' item) end)
								plus(doc-lines: ', ')
							end)
						'in'
							plus('allowed: [ ' call(stdfu.apply rest(check) func(x) plus(str(x) ' ') end): '], ')
						error('illegal tag: ' str(head(check)))
					)
					plus(out-str fout)
				end
			end

			keyv checklist = kvitem:
			field-output = call(stdfu.loop call(get-field-visitor keyv) checklist '')
			plus(output '\n' indent str(keyv) ' : ' field-output)
		end

		if( eq(type(map-schema) 'map')
			call(func()
				map-output = call(stdfu.loop map-fields-visit keyvals(map-schema) '')
				plus('\n' indent 'map: ' map-output)
			end)
			'requires map'
		)
	end

	if( and( eq(type(schema) 'list') not(empty(schema)))
		case( head(schema)
			'map' call(visit-map head(rest(schema)) '' '')
			'none'
		)
		'requires non-empty list'
	)
end

# validates data against schema
validate = func(schema srcdata)
	import stdfu

	validate-map = func(map-schema map-data field-path)
		map-fields-checker = func(kvitem results)
			get-field-checker = func(mkey)
				func(check result-list)
					_ msg-list = result-list:
					passed message = case( head(check)
						'required'
							if( in(map-data mkey)
								result-list
								list(false append(msg-list sprintf('required field %v not found (%s)' mkey field-path)))
							)

						'in'
							call(func()
								key-found val = getl(map-data mkey):
								allowed = rest(check)
								cond(
									not(key-found)
										result-list
									not(in(allowed val))
										list(false append(msg-list sprintf('field %v is not in allowed set (%v not in: %v)(%s)' mkey val allowed field-path)))
									result-list
								)
							end)

						'type'
							call(func()
								key-found val = getl(map-data mkey):
								required-type = head(rest(check))
								cond(
									not(key-found)
										result-list
									not(eq(type(val) required-type))
										list(false append(msg-list sprintf('field %v is not required type (got: %v, expected: %v)(%s)' mkey type(val) required-type field-path)))
									result-list
								)
							end)

						'map'
							call(func()
								submap-found submap-data = getl(map-data mkey):
								cond(
									not(submap-found)
										result-list
									not(eq(type(submap-data) 'map'))
										list(false append(msg-list sprintf('field %v is not map (%s)' mkey field-path)))
									call(validate-map head(rest(check)) submap-data plus(field-path ' -> ' str(mkey)))
								)
							end)

						'doc'
							result-list

						list(false append(msg-list sprintf('unknown validator: %s' str(head(check)))))
					):
					list(passed message)
				end
			end

			keyv checklist = kvitem:
			call(stdfu.loop call(get-field-checker keyv) checklist results)
		end

		if( eq(type(map-schema) 'map')
			call(stdfu.loop map-fields-checker keyvals(map-schema) list(true list()))
			list(false list('requires map'))
		)
	end

	if( and( eq(type(schema) 'list') not(empty(schema)))
		case( head(schema)
			'map' call(validate-map head(rest(schema)) srcdata '')
			list(false list(sprintf('unknown validator: %s' str(head(schema)) )))
		)
		list(false list('requires non-empty list'))
	)
end

endns

`

	stdfunMap["stdpp"] = `
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

`

	stdfunMap["stdpr"] = `
ns stdpr

get-pr = func(do-print)
	if( do-print
		func(text value)
			_ = print(text value)
			value
		end

		func(_ value) value end
	)
end

get-pp-pr = func(do-print)
	if( do-print
		func(text value)
			import stdpp

			_ = print(text call(stdpp.form value))
			value
		end

		func(_ value) value end
	)
end

endns
`

	stdfunMap["stdser"] = `
ns stdser

import stdjson
import stdbytes
import stdfu
import stdbase64

tags = list('int' 'float' 'bool' 'string' 'list' 'map' 'bytearray')

encode = func(val)
	enc-bytearray = func(inval)
		enc-ok enc-err retv = call(stdbase64.encode inval):
		_ = if(enc-ok '' error(enc-err))
		retv
	end

	handle-item = func(inval)
		vtype = type(inval)
		case( vtype
			'int'    list('int' inval)
			'float'  list('float' inval)
			'bool'   list('bool' inval)
			'string' list('string' inval)
			'opaque:bytearray' list('bytearray' call(enc-bytearray inval))
			'list'   list('list' call(stdfu.apply inval func(item) call(handle-item item) end))
			'map'    list('map' call(stdfu.apply keyvals(inval) func(item) k v = item: list(call(handle-item k) call(handle-item v)) end))
			error('unsupported type: ' vtype)
		)
	end

	call(stdjson.encode call(handle-item val))
end

decode = func(val)
	handle-map = func(ml)
		mapper = func(pair resultm)
			kpair vpair = pair:
			put(resultm call(handle-pair kpair) call(handle-pair vpair))
		end

		call(stdfu.loop mapper ml map())
	end

	dec-bytearray = func(inval)
		dec-ok dec-err retv = call(stdbase64.decode inval):
		_ = if(dec-ok '' error(dec-err))
		retv
	end

	handle-pair = func(pairval)
		tag value = pairval:
		case( tag
			'int'    value
			'float'  value
			'bool'   value
			'string' value
			'bytearray' call(dec-bytearray value)
			'list'   call(stdfu.apply value func(pair) call(handle-pair pair) end)
			'map'    call(handle-map value)
			error('unsupported tag: ' tag)
		)
	end

	ok err fuval = call(stdjson.decode val):
	if( ok
		list(true '' call(handle-pair fuval))
		list(false err fuval)
	)
end

endns

`

	stdfunMap["stdset"] = `
ns stdset

# newset creates new set
newset = func()
	map()
end

# is-empty returns true if there is no items in set, otherwise true
is-empty = func(set)
	eq(len(set), 0)
end

# setlen returns number of items in set
setlen = func(set)
	len(set)
end

# as-list returns set items in list
as-list = func(set)
	keys(set)
end

# has-item returns true if item is in set, false otherwise
has-item = func(set, item)
	in(set, item)
end

# removes item from set (if it is in it)
remove-from-set = func(set item)
	if( in(set item)
		del(set item)
		set
	)
end

# add-to-set adds one item to set
add-to-set = func(set, item)
	if( in(set, item),
		set,
		put(set, item, true)
	)
end

# list-to-set adds one item to set
list-to-set = func(set, itemlist)
	looper = func(iteml, setv)
		while( not(empty(iteml)),
			rest(iteml),
			call(add-to-set, setv, head(iteml)),
			setv
		)
	end
	
	call(looper, itemlist, set)
end

# union creates union of two sets given as arguments
union = func(set1, set2)
	set-to-add = if( gt(len(set1), len(set2)),
		set2,
		set1
	)
	target-set = if( gt(len(set1), len(set2)),
		set1,
		set2
	)
	call(list-to-set, target-set, keys(set-to-add))
end

# helper function
get-matching-subset = func(source-list, condition)
	looper = func(iteml, resultl)
		while( not(empty(iteml)),
			rest(iteml),
			call(func()
				item = head(iteml)
				if( call(condition, item),
					append(resultl, item),
					resultl
				)
			end),
			resultl
		)
	end

	itemlist = call(looper, source-list, list())
	call(list-to-set, call(newset), itemlist)	
end

# intersection creates intersection set of two sets given as arguments
intersection = func(set1, set2)
	keys1 = keys(set1)
	keys2 = keys(set2)

	condition = func(item) and( in(keys1, item), in(keys2, item) ) end
	call(get-matching-subset, extend(keys1, keys2), condition)
end

# difference returns set with elements in set1 but not in set2
difference = func(set1, set2)
	keys1 = keys(set1)
	keys2 = keys(set2)

	condition = func(item) and( in(keys1, item), not(in(keys2, item)) ) end
	call(get-matching-subset, extend(keys1, keys2), condition)
end

# is-subset return true if subset -argument is subset of set -argument 
is-subset = func(set, subset)
	looper = func(iteml, result)
		while( not(empty(iteml)),
			rest(iteml),
			and(result, in(set, head(iteml))),
			result
		)
	end
	
	call(looper, keys(subset), true)
end

# equal returns true if two sets given as arguments are having same items, false otherwise
equal = func(set1, set2)
	len1 = call(setlen, set1)
	len2 = call(setlen, set2)

	if( eq(len1, len2),
		call(is-subset, set1, set2),
		false
	)
end

endns
`

	stdfunMap["stdsort"] = `
ns stdsort

# merges two sorted lists to one
merge = func(lst1 lst2)
	do-merge = func(rl1 rl2 res)
		ready nrl1 nrl2 nres = cond(
			empty(rl1) cond(
				empty(rl2) list(true rl1 rl2 res)
				list(false rl1 rest(rl2) append(res head(rl2)))
			)

			empty(rl2) cond(
				empty(rl1) list(true rl1 rl2 res)
				list(false rest(rl1) rl2 append(res head(rl1)))
			)

			if( lt(head(rl1) head(rl2))
				list(false rest(rl1) rl2 append(res head(rl1)))
				list(false rl1 rest(rl2) append(res head(rl2)))
			)
		):

		while( not(ready)
			nrl1
			nrl2
			nres
			nres
		)
	end

	call(do-merge lst1 lst2 list())
end

# merges two sorted lists to one
merge-with-func = func(lst1 lst2 compa-func)
	do-merge = func(rl1 rl2 res)
		ready nrl1 nrl2 nres = cond(
			empty(rl1) cond(
				empty(rl2) list(true rl1 rl2 res)
				list(false rl1 rest(rl2) append(res head(rl2)))
			)

			empty(rl2) cond(
				empty(rl1) list(true rl1 rl2 res)
				list(false rest(rl1) rl2 append(res head(rl1)))
			)

			if( call(compa-func head(rl1) head(rl2))
				list(false rest(rl1) rl2 append(res head(rl1)))
				list(false rl1 rest(rl2) append(res head(rl2)))
			)
		):

		while( not(ready)
			nrl1
			nrl2
			nres
			nres
		)
	end

	call(do-merge lst1 lst2 list())
end

# implements sort by using mergesort
sort = func(src-lst)
	use-func compa-func = if( eq(len(argslist()) 2)
		list(true last(argslist()))
		list(false func() 'not used' end)
	):

	real-sort = func(lst)
		l = len(lst)

		get-slices = func()
			middle = div(l 2)
			left = slice(lst 0 middle)
			right = slice(lst plus(middle 1) l)
			list(left right)
		end

		left-lst right-lst = call(get-slices):
		case( l
			0 error('unexpected empty list')
			1 lst
			2 if( use-func
				call(merge-with-func list(head(lst)) list(last(lst)) compa-func)
				call(merge list(head(lst)) list(last(lst)) )
			)
			if( use-func
				call(merge-with-func call(real-sort left-lst) call(real-sort right-lst) compa-func)
				call(merge call(real-sort left-lst) call(real-sort right-lst))
			)
		)
	end

	if( empty(src-lst)
		list()
		call(real-sort src-lst)
	)
end

endns

`
}
