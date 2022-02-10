
ns main

import stdhttp

main = proc(port)
	hello-handler = proc(w r)
		import stdbytes
		call(stdhttp.write-response w 200 call(stdbytes.str-to-bytes 'Hi There !'))
	end

	mux = call(stdhttp.mux)
	_ = call(stdhttp.reg-handler mux '/hello' hello-handler)
	address = plus(':' str(port))
	call(stdhttp.listen-and-serve mux address)
end

endns

