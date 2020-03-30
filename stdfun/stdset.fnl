
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
