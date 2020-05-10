
ns main

import stdhttp
import stdlog

# HTTP status codes
Status-OK                    = 200 # OK HTTP status code
Status-Bad-Request           = 400 # Bad Request HTTP status code
Status-Internal-Server-Error = 500 # Internal Server Error HTTP status code

# Logger to use in this server
log = call(stdlog.get-default-logger map('prefix' 'todo-server: ' 'date' true 'time' true))

# none is value which content does not matter
none = 'whatever'

# Starts server  fiber, returns channel with which item server is accessed
run-server = proc()
	import stdjson

	# helper for writing OK HTTP response
	write-success-response = proc(w data)
		_ = call(stdhttp.set-response-header w map('Content-Type' 'application/json'))
		status = call(stdhttp.write-response w Status-OK data)

		if( not(eq(status ''))
			call(log 'error: ' status)
			none
		)
	end

	# helper for writing Error HTTP response
	write-error-response = proc(w error-text error-code)
		_ = call(stdhttp.set-response-header w map('Content-Type' 'application/json'))
		_ _ error-body = call(stdjson.encode error-text):
		call(stdhttp.write-response w error-code error-body)
	end

	# handles GET /items
	get-items = proc(w r items)
		ok err data = call(stdjson.encode items):
		_ = if( ok
			call(write-success-response w data)
			call(proc()
				_ = call(log 'error in encoding items: ' err)
				call(write-error-response w 'error in encoding items' Status-Internal-Server-Error)
			end)
		)
		items
	end

	# handles POST /items
	put-item = proc(w r items id)
		ok err item = call(stdjson.decode get(r 'body')):
		_ = if( not(ok)
			call(proc()
				_ = call(log 'error in decoding: ' err)
				call(write-error-response w 'error in decoding request' Status-Bad-Request)
			end)
			none
		)
		if(ok append(items put(item 'id' id)) items)
	end

	# handles DELETE /items/id/<id>
	del-item = proc(w r items)
		import stdfu

		item-id = conv(last(split(get(r 'URI') '/')) 'int')
		call(stdfu.filter items func(item) not( eq(get(item 'id') item-id) ) end)
	end

	# server fiber implementation
	item-server = proc(server-ch)
		serving = proc(items id-counter)
			while( true
				call(proc()
					w r replych = recv(server-ch):
					new-items = case( get(r 'method')
						'GET'    call(get-items w r items)
						'POST'   call(put-item w r items id-counter)
						'DELETE' call(del-item w r items)
						items
					)
					_ = send(replych none)
					new-items
				end)
				plus(id-counter 1)
				none
			)
		end

		call(serving list() 100)
	end

	# create channel, start server fiber and return channel
	ch = chan()
	_ = spawn(call(item-server ch))
	ch
end

# returns handler for /items (POST, GET, DELETE)
get-item-handler = func(ch supported-methods)
	send-and-wait = proc(w r)
		replych = chan()
		_ = send(ch list(w r replych))
		recv(replych) # need to wait reply, otherwise connection will be lost
	end

	proc(w r)
		if( in(supported-methods get(r 'method'))
			call(send-and-wait w r)
			none
		)
	end
end

# HTTP handlers and server
server-main = proc(mux item-server-ch)
	import todo_common

	_ = call(stdhttp.reg-handler mux '/items' call(get-item-handler item-server-ch list('GET' 'POST')))
	_ = call(stdhttp.reg-handler mux '/items/id/' call(get-item-handler item-server-ch list('DELETE')))
	_ = call(log '...listening...')
	retv = call(stdhttp.listen-and-serve mux call(todo_common.get-port))
	plus('Quit serving: ' retv)
end

# main sets up item-server fiber and HTTP server and signal handler
main = proc()
	import stdos

	mux = call(stdhttp.mux)
	sig-handler = proc(signum sigtext)
		_ = call(log 'signal received' signum sigtext)
		call(stdhttp.shutdown mux)
	end
	_ = call(stdos.reg-signal-handler sig-handler 2)

	call(server-main mux call(run-server))
end

endns

