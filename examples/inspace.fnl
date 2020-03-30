
ns main

# https://pragprog.com/book/bhwb/exercises-for-programmers
# http://api.open-notify.org/astros.json

get-astros = proc()
	import stdhttp 
	import stdjson

	response = call(stdhttp.do 'GET' 'http://api.open-notify.org/astros.json' map())
	_ _ val = call(stdjson.decode get(response 'body')):
	val
end

by-craft = func(astro-list)
	import stdfu
	f = func(item)
		craft = get(item 'craft')
		astro-name = get(item 'name')
		list(craft astro-name)
	end
	call(stdfu.group-by astro-list f)
end

main = proc()
	result = call(get-astros)
	astro-list = get(result 'people')
	call(by-craft astro-list)
end

endns
