
ns todo_common

# returns TCP port which is listened
get-port = proc()
	import stdos

	found port = call(stdos.getenv 'TODO_SRV_PORT'):
	if(not(found) ':8003' plus(':' port))
end

endns

