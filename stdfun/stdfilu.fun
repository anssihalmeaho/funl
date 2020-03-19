
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
