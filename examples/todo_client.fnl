
ns main

import stdio
import stdhttp
import stdjson
import stdbytes

import todo_common

# HTTP status codes
Status-OK                    = 200 # OK HTTP status code
Status-Bad-Request           = 400 # Bad Request HTTP status code
Status-Internal-Server-Error = 500 # Internal Server Error HTTP status code

# none is value which content does not matter
none = 'whatever'

# Server URL
server-endpoint = sprintf('http://localhost%s/items' call(todo_common.get-port))

# Printing put/del result
print-result = proc(response)
	code = get(response 'status-code')
	text = get(response 'status-text')
	case( code
		Status-OK call(stdio.printline 'ok')
		call(stdio.printline sprintf('error: %d, %s: %s' code text call(stdbytes.string get(response 'body'))))
	)
end

# Inquires items from server
get-items = proc()
	print-inq-result = proc(data)
		import stdfu

		item-form = func(item)
			item-formatter = func(item-as-map out)
				key val = item-as-map:
				sprintf('%s\n  - %s : %v' out key val)
			end
			call(stdfu.foreach keyvals(item) item-formatter '')
		end

		out = call(stdfu.foreach data func(item output) plus(output '\n item: ' call(item-form item) '\n') end '')
		call(stdio.printline out)
	end

	result = call(stdhttp.do 'GET' server-endpoint map())
	statuscode = get(result 'status-code')
	case( statuscode
		Status-OK call(print-inq-result last(call(stdjson.decode get(result 'body'))) )
		sprintf('error from server: %d : %s' statuscode get(result 'status-text'))
	)
end

# Sends new item to server
put-item = proc(cmd input)
	send-item-to-server = proc(body-content)
		header = map('Content-Type' 'application/json')
		response = call(stdhttp.do 'POST' server-endpoint header body-content)
		response
	end

	if( eq(cmd input)
		'no item given'
		call(proc()
			json-item = slice(input len(cmd)) # remove command part from front of string
			response = call(send-item-to-server call(stdbytes.str-to-bytes json-item))
			call(print-result response)
		end)
	)
end

# Removes item from server
del-item = proc(cmd input)
	import stdstr

	remove-item-from-server = proc(id)
		delete-endpoint = plus(server-endpoint '/id/' id)
		response = call(stdhttp.do 'DELETE' delete-endpoint map())
		response
	end

	if( eq(cmd input)
		'no id given'
		call(proc()
			param = slice(input len(cmd)) # remove command part from front of string
			id = call(stdstr.strip param)
			if( and( not(eq(id '')) call(stdstr.is-digit id) )
				call(proc()
					response = call(remove-item-from-server id)
					call(print-result response)
				end)
				call(stdio.printline 'invalid id: ' id)
			)
		end)
	)
end

# Executes commands
exec-cmd = proc(input)
	cmd = head(split(input))
	_ = case( cmd
		'get' call(get-items)
		'put' call(put-item cmd input)
		'del' call(del-item cmd input)
		none
	)
	none
end

# Printing help for user
print-help = proc()
	_ = call(stdio.printline 'Input can be:')
	_ = call(stdio.printline '  help         -> prints this help')
	_ = call(stdio.printline '  ?            -> prints this help')
	_ = call(stdio.printline '  quit         -> exits repl')
	_ = call(stdio.printline '  exit         -> exits repl')
	_ = call(stdio.printline '')
	_ = call(stdio.printline '  put <JSON value for item> -> adds item')
	_ = call(stdio.printline '  get                       -> prints all items')
	_ = call(stdio.printline '  del <id of item>          -> removes item')
	_ = call(stdio.printline '')
	true
end

# Main "loop" of program
exec-commands = proc()
	_ = call(stdio.printout 'todo> ')
	input = call(stdio.readinput)
	continue = not(in(list('quit' 'exit') input))

	result = case( input
		''     true
		'quit' false
		'exit' false
		'help' call(print-help)
		'?'    call(print-help)
		call(exec-cmd input)
	)

	while(continue 'bye')
end

# main (entrypoint)
main = proc()
	_ = call(stdio.printline 'Welcome to Todo application client')
	call(exec-commands)
end

endns
