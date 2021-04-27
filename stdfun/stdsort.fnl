
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

