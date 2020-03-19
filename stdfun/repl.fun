
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
