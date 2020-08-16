
ns stdser

import stdjson
import stdbytes
import stdfu

tags = list('int' 'float' 'bool' 'string' 'list' 'map')

encode = func(val)
	handle-item = func(inval)
		vtype = type(inval)
		case( vtype
			'int'    list('int' inval)
			'float'  list('float' inval)
			'bool'   list('bool' inval)
			'string' list('string' inval)
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

	handle-pair = func(pairval)
		tag value = pairval:
		case( tag
			'int'    value
			'float'  value
			'bool'   value
			'string' value
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

